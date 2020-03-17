package poolrpc

import (
	"net/http"
	"time"
)

// type httpClient interface {
// 	Post(url string, contentType string, body io.Reader) (*http.Response, error)
// }

type logger interface {
	Println(v ...interface{})
}

// WithHttpClient set custom http client
func WithPoolHttpClient(client *http.Client) func(rpc *PoolRPC) {
	return func(rpc *PoolRPC) {
		rpc.client = client
	}
}

// WithLogger set custom logger
func WithPoolLogger(l logger) func(rpc *PoolRPC) {
	return func(rpc *PoolRPC) {
		rpc.log = l
	}
}

// WithDebug set debug flag
func WithPoolDebug(enabled bool) func(rpc *PoolRPC) {
	return func(rpc *PoolRPC) {
		rpc.Debug = enabled
	}
}

// type strFloat string
// func (w *strFloat) UnmarshalJSON(value []byte) error {
// 	*w = strFloat(string(value))
// 	return nil
// }
const (
	PoolOrderUnknown   = 0 // 未知
	PoolOrderCreating  = 1
	PoolOrderCreated   = 2 // 已运行
	PoolOrderReSetting = 3 // 修改了挖矿配置
	PoolOrderRunning   = 4 // 订单运行中
	PoolOrderCompleted = 5
)

// orderInfo 订单详情 创建
type OrderInfo struct {
	FarmID string       `json:"FarmID"` // 矿场ID
	Type   string       `json:"Type"`   // 机器类型
	Price  MachinePrice `json:"Price"`  // 定价单
	Config MinerConfig  `json:"config"` // 配置单
}

type MachinePrice struct {
	Price uint64 `json:"Price"` // 价格，人民币
	Time  uint64 `json:"Time"`  // 秒数
}
type MinerConfig struct {
	MinerPool     string `json:"MinerPool"`
	MinerUserName string `json:"MinerUsername"`
	MinerWorker   string `json:"MinerWorker"`
}

// 订单结果 查询返回，需要传入订单号
type OrderRes struct {
	ID        string       `json:"id"`
	FarmID    string       `json:"farm_id"`
	MinerID   string       `json:"miner_id"`
	CreatedAt time.Time    `json:"create_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Status    int          `json:"status"`
	Request   OrderRequest `json:"request"`
}

type OrderRequest struct {
	FarmID string       `json:"FarmID"` //订单的 farmid
	Type   string       `json:"Type"`   //订单 矿机类型
	Price  MachinePrice `json:"Price"`  //订单价格
	Config MinerConfig  `json:"Config"` //订单矿机设置
}

// 矿场机器详情,可以添加 查询返回
type FarmServer struct {
	ID        string               `json:"id"` //矿场名称
	UpdatedAt time.Time            `json:"updated_at"`
	Miners    map[string]MinerInfo `json:"miners"` //价格表
	Orders    map[string]OrderRes  `json:"orders,omitempty"`
}

type MinerInfo struct {
	Type           string
	PriceList      []MachinePrice
	AvailableCount uint64
}
