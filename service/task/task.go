package task

import (
	"container/heap" // Golang提供的heap库
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

const (
	MIN_TIMER_INTERVAL = 500 * time.Millisecond // 循环定时器的最小时间间隔
)

var (
	nextAddSeq uint = 1 // 用于为每个定时器对象生成一个唯一的递增的序号
)

// 定时器对象
type Timer struct {
	fireTime     time.Time              // 触发时间
	interval     time.Duration          // 时间间隔（用于循环定时器）
	callback     CallbackFunc           // 回调函数
	param        map[string]interface{} // 回调函数参数
	repeat       bool                   // 是否循环
	deadlineTime time.Time              // 截止时间
	deadline     bool                   // 是否记次
	cancelled    bool                   // 是否已经取消
	addseq       uint                   // 序号
}

// 取消一个定时器，这个定时器将不会被触发
func (t *Timer) Cancel() {
	t.cancelled = true
}

// 判断定时器是否已经取消
func (t *Timer) IsActive() bool {
	return !t.cancelled
}

// 使用一个heap管理所有的定时器
type _TimerHeap struct {
	timers []*Timer
}

// Golang要求heap必须实现下面这些函数，这些函数的含义都是不言自明的

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

// 使用触发时间和需要对定时器进行比较
func (h *_TimerHeap) Less(i, j int) bool {
	//log.Println(h.timers[i].fireTime, h.timers[j].fireTime)
	t1, t2 := h.timers[i].fireTime, h.timers[j].fireTime
	if t1.Before(t2) {
		return true
	}

	if t1.After(t2) {
		return false
	}
	// t1 == t2, making sure Timer with same deadline is fired according to their add order
	return h.timers[i].addseq < h.timers[j].addseq
}

func (h *_TimerHeap) Swap(i, j int) {
	var tmp *Timer
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *_TimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(*Timer))
}

func (h *_TimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

// 定时器回调函数的类型定义
type CallbackFunc func(param map[string]interface{}) bool

var (
	timerHeap     _TimerHeap // 定时器heap对象
	timerHeapLock sync.Mutex // 一个全局的锁
)

func init() {
	heap.Init(&timerHeap) // 初始化定时器heap
}

// 设置一个一次性的回调，这个回调将在d时间后触发，并调用callback函数
func AddCallback(d time.Duration, callback CallbackFunc) *Timer {
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		// param:        ,
		deadlineTime: time.Now(),
		deadline:     false,
		repeat:       false,
	}
	timerHeapLock.Lock() // 使用锁规避竞争条件
	t.addseq = nextAddSeq
	nextAddSeq += 1

	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

// 设置一个指定截止时间的定时触发回调，这个回调将在d时间后第一次触发，以后每隔d时间重复触发，并调用callback函数，直到达到截止时间
func AddTimerWithDeadLine(d, deadline time.Duration, param map[string]interface{}, callback CallbackFunc) *Timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}

	t := &Timer{
		fireTime:     time.Now().Add(d),
		interval:     d,
		callback:     callback,
		param:        param,
		deadlineTime: time.Now().Add(deadline),
		deadline:     true,
		repeat:       true, // 设置为循环定时器
	}
	timerHeapLock.Lock()
	t.addseq = nextAddSeq // set addseq when locked
	nextAddSeq += 1

	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

// 设置一个指定截止时间的定时触发回调，这个回调将在 first 时间后第一次触发，以后每隔d时间重复触发，并调用callback函数
func AddTimerWithNoDeadLine(d, first time.Duration, param map[string]interface{}, callback CallbackFunc) *Timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}

	t := &Timer{
		fireTime:     time.Now().Add(first),
		interval:     d,
		callback:     callback,
		param:        param,
		deadlineTime: time.Now(),
		deadline:     false,
		repeat:       true, // 设置为循环定时器
	}
	timerHeapLock.Lock()
	t.addseq = nextAddSeq // set addseq when locked
	nextAddSeq += 1

	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

// 设置一个定时触发的回调，这个回调将在d时间后第一次触发，以后每隔d时间重复触发，并调用callback函数
func AddTimer(d time.Duration, callback CallbackFunc) *Timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}

	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		// param:        "",
		deadlineTime: time.Now(),
		deadline:     false,
		repeat:       true, // 设置为循环定时器
	}
	timerHeapLock.Lock()
	t.addseq = nextAddSeq // set addseq when locked
	nextAddSeq += 1

	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

// 对定时器模块进行一次Tick
//
// 一般上层模块需要在一个主线程的goroutine里按一定的时间间隔不停的调用Tick函数，从而确保timer能够按时触发，并且
// 所有Timer的回调函数也在这个goroutine里运行。
func Tick() {
	now := time.Now()
	timerHeapLock.Lock()
	var result bool
	for {
		if timerHeap.Len() <= 0 { // 没有任何定时器，立刻返回
			break
		}
		nextFireTime := timerHeap.timers[0].fireTime
		if nextFireTime.After(now) { // 没有到时间的定时器，返回
			break
		}

		t := heap.Pop(&timerHeap).(*Timer)
		// fmt.Printf("%+v\n", t)
		if t.cancelled { // 忽略已经取消的定时器
			continue
		}
		deadlineTime := t.deadlineTime
		if t.deadline && deadlineTime.Before(now) { //如果达到截止时间则取消任务
			t.Cancel()
			continue
		}
		if !t.repeat {
			t.cancelled = true
		}
		timerHeapLock.Unlock()
		result = runCallback(t.callback, t.param) // 运行回调函数并捕获panic
		if result {                               // 任务执行成功，删除
			t.Cancel()
		}
		timerHeapLock.Lock()

		if t.repeat { // 如果是循环timer就把Timer重新放回heap中
			// add Timer back to heap
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) {
				t.fireTime = now.Add(t.interval)
			}
			t.addseq = nextAddSeq
			nextAddSeq += 1
			heap.Push(&timerHeap, t)
		}
	}
	timerHeapLock.Unlock()
}

// 创建一个 goroutine 对定时器模块进行定时的 Tick
func InitTaskTicks(tickInterval time.Duration) {
	go selfTickRoutine(tickInterval)
}

func selfTickRoutine(tickInterval time.Duration) {
	for {
		time.Sleep(tickInterval)
		Tick()
	}
}

// 运行定时器的回调函数，并捕获panic，将panic转化为错误输出
func runCallback(callback CallbackFunc, param map[string]interface{}) bool {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback %v paniced: %v\n", callback, err)
			debug.PrintStack()
		}
	}()
	return callback(param)
}
