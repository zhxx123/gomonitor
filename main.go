// main.go 文件
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/controller"
	"github.com/zhxx123/gomonitor/cron"
	"github.com/zhxx123/gomonitor/middleware"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/router"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/log"
	"github.com/zhxx123/gomonitor/service/mail"
	"github.com/zhxx123/gomonitor/service/order"
	"github.com/zhxx123/gomonitor/service/pay"
	"github.com/zhxx123/gomonitor/service/set"
	"github.com/zhxx123/gomonitor/service/spider"
	"github.com/zhxx123/gomonitor/service/task"
	"github.com/zhxx123/gomonitor/service/wallet"
)

/**
 * 初始化 iris app
 */

func InitIris() {
	api := iris.New()
	api.Logger().SetLevel("info")
	api.Use(logger.New())
	// api 调用统计
	api.Use(middleware.APIStatsD())
	// 404 错误
	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(controller.ApiResource(model.ErrorCode.NotFound, nil, "404 Not Found"))
	})
	// 500 错误
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.WriteString("Oups something went wrong, try again")
	})
	// iris 中断获取
	iris.RegisterOnInterrupt(func() {
		db.CloseAllDB()
	})
	// 路由注册
	router.RegisterRoute(api)

	// 设置最大提交form 表单限制  (default is 32 MiB)
	maxSize := int64(config.ServerConfig.MaxMultipartMemory)
	apiPostMaxSize := maxSize << 20 // 2 MiB

	addr := config.ServerConfig.Addr
	tlsEnanled := config.ServerConfig.APITLSEnabled
	var err error
	if tlsEnanled { // tls 监听
		tlsCert := config.ServerConfig.TLSCertPath
		tlskey := config.ServerConfig.TLSKeyPath
		err = api.Run(iris.TLS(addr, tlsCert, tlskey), iris.WithPostMaxMemory(apiPostMaxSize))
	} else { // 普通监听
		err = api.Run(iris.Addr(addr), iris.WithPostMaxMemory(apiPostMaxSize))
	}
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	return
}

//性能监控
func goPprof() {
	go func() {
		http.ListenAndServe("127.0.0.1:8081", nil)
	}()
}

// 注册 信号
func signalNotify() {
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-c
		fmt.Printf("[got signal: %s], exiting goblog now\n", s)
		shutdownAll()
	}()
}

// 关闭所有
func shutdownAll() {
	// 关闭日志
	db.CloseAllDB()
	log.CloseLog()
	fmt.Println("[All] Shutdown")
	os.Exit(0)
}

// 初始化
func init() {
	log.InitLog()                       // 初始化日志
	db.InitDB()                         // 初始化数据库
	mail.InitEmail()                    // 初始化邮件
	task.InitTaskTicks(1 * time.Second) // 初始化 定时器模块，2 秒扫描一次
	pay.InitPay()                       // 初始化支付模块
	spider.InitSpider()                 // 初始化 虚拟货币价格爬虫
	order.InitPoolRPC()                 // 初始化 后台矿机操作模块
	cron.InitCron()                     // 初始化定时任务
	// 刷新数据库配置到系统缓存
	set.InitSysSeting()          // 系统设置
	set.UpdateUserSafeSeting()   // 更新安全设置
	wallet.UpdateWalletConfirm() // 更新交易设置
}

// 系统开始，主函数
func main() {
	signalNotify() // 注册监听程序
	// goPprof() // 程序性能监控
	InitIris() // web 启动
}
