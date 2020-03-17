package poolrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zhxx123/gomonitor/model"
)

// rpcError - massgrid error
type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err RpcError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type rpcResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *RpcError       `json:"error"`
}

type rpcRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// PoolRPC - Poolereum rpc client
type PoolRPC struct {
	url    string
	user   string
	passwd string
	client *http.Client
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func New(url, user, passwd string, options ...func(rpc *PoolRPC)) *PoolRPC {
	rpc := &PoolRPC{
		url:    url,
		user:   user,
		passwd: passwd,
		client: http.DefaultClient,
		log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// NewPoolRPC create new rpc client with given url
func NewPoolRPC(url, user, passwd string, options ...func(rpc *PoolRPC)) *PoolRPC {
	return New(url, user, passwd, options...)
}

func (rpc *PoolRPC) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// URL returns client url
func (rpc *PoolRPC) URL() string {
	return rpc.url
}

// Call returns raw response of method call
func (rpc *PoolRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := rpcRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", rpc.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(rpc.user, rpc.passwd)
	response, err := rpc.client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if rpc.Debug {
		rpc.log.Println(fmt.Sprintf("\nURL:%s\nMethod:%s\nHeader:%v\nRequest: %s\nResponse: %s\n", rpc.url, method, req.Header, body, data))
	}

	resp := new(rpcResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *PoolRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

// *** The above is the latest block information.

// 创建订单
func (rpc *PoolRPC) CreateOrder(orderInfo *OrderInfo) (string, error) {

	result, err := rpc.RawCall("Order.Create", orderInfo)
	if err != nil {
		return "", err
	}
	var orderId string
	if bytes.Equal(result, []byte("false")) {
		return orderId, nil
	}
	err = json.Unmarshal(result, &orderId)
	return orderId, err
}

// 获取订单详情 old !!!
// func (rpc *PoolRPC) GetOrderInfo(orderId string) (*OrderRes, error) {
// 	result, err := rpc.RawCall("Order.GetOrderInfo")
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res OrderRes
// 	if bytes.Equal(result, []byte("false")) {
// 		return &res, nil
// 	}
// 	err = json.Unmarshal(result, &res)
// 	return &res, err
// }

func (rpc *PoolRPC) GetOrderInfo(orderId string) (*model.MinerOrder, error) {
	result, err := rpc.RawCall("Order.GetOrderInfo")
	if err != nil {
		return nil, err
	}
	var res model.MinerOrder
	if bytes.Equal(result, []byte("false")) {
		return &res, nil
	}
	err = json.Unmarshal(result, &res)
	return &res, err
}

// 获取矿场机器详情 old !!!
// func (rpc *PoolRPC) GetFarmsInfo() (map[string]FarmServer, error) {
// 	var params []map[string]interface{}
// 	result, err := rpc.RawCall("Order.GetFarmsInfo", params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var farmServer map[string]FarmServer

// 	if bytes.Equal(result, []byte("[]")) {
// 		return farmServer, nil
// 	}
// 	err = json.Unmarshal(result, &farmServer)
// 	return farmServer, err
// }

func (rpc *PoolRPC) GetFarmsInfo() ([]model.FarmServer, error) {
	var params []map[string]interface{}
	result, err := rpc.RawCall("Order.GetFarmsInfo", params)
	if err != nil {
		return nil, err
	}
	var farmServer []model.FarmServer

	if bytes.Equal(result, []byte("[]")) {
		return farmServer, nil
	}
	err = json.Unmarshal(result, &farmServer)
	return farmServer, err
}

// 获取指定 farminfo old !!!
// func (rpc *PoolRPC) GetFarmInfo(farmId string) (FarmServer, error) {
// 	result, err := rpc.RawCall("Order.GetFarmInfo", farmId)
// 	var farmServer FarmServer
// 	if err != nil {
// 		return farmServer, err
// 	}

// 	if bytes.Equal(result, []byte("[]")) {
// 		return farmServer, nil
// 	}
// 	err = json.Unmarshal(result, &farmServer)
// 	return farmServer, err
// }
func (rpc *PoolRPC) GetFarmInfo(farmId string) ([]model.FarmServer, error) {
	result, err := rpc.RawCall("Order.GetFarmInfo", farmId)
	var farmServer []model.FarmServer
	if err != nil {
		return farmServer, err
	}

	if bytes.Equal(result, []byte("[]")) {
		return farmServer, nil
	}
	err = json.Unmarshal(result, &farmServer)
	return farmServer, err
}
