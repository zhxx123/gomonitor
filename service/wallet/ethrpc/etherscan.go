package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zhxx123/gomonitor/service/wallet/utils"
)

// EthError - ethereum error
type EtherError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err EtherError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type etherResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Error   *EtherError     `json:"error"`
}

type etherRequest struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Module  string `json:"module"`
	Action  string `json:"action"`
	// APIKey  string                 `json:"apikey"`
	Params map[string]interface{} `json:"params"`
}

// EthRPC - Ethereum rpc client
type EtherRPC struct {
	url    string
	client *http.Client
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func NewEther(url string, options ...func(rpc *EtherRPC)) *EtherRPC {
	rpc := &EtherRPC{
		url:    url,
		client: http.DefaultClient,
		log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// NewEthRPC create new rpc client with given url
func NewEtherRPC(url string, options ...func(rpc *EtherRPC)) *EtherRPC {
	return NewEther(url, options...)
}

func (rpc *EtherRPC) call(query map[string]interface{}, target interface{}) error {
	result, err := rpc.Call(query)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// // URL returns client url
// func (rpc *EtherRPC) URL() string {
// 	return rpc.url
// }

// Call returns raw response of method call
func NewRequest(method, url string, params url.Values) (*http.Request, error) {
	var m = strings.ToUpper(method)
	var body io.Reader
	if m == "GET" || m == "HEAD" {
		if len(params) > 0 {
			if strings.Contains(url, "?") {
				url = url + "&" + params.Encode()
			} else {
				url = url + "?" + params.Encode()
			}
		}
	} else {
		body = strings.NewReader(params.Encode())
	}
	return http.NewRequest(m, url, body)
}
func (rpc *EtherRPC) Call(params map[string]interface{}) (json.RawMessage, error) {
	req, err := http.NewRequest("GET", rpc.url, nil)
	if params != nil {
		p := req.URL.Query()
		for key, val := range params {
			if val != "" {
				p.Add(key, val.(string))
			}
		}
		req.URL.RawQuery = p.Encode()
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
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
		rpc.log.Println(fmt.Sprintf("GET\nHeader:%v\nRequest: %v\nResponse: %s\n", req.Header, req.URL.String(), data))
	}

	resp := new(etherResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *EtherRPC) RawCall(params map[string]interface{}) (json.RawMessage, error) {

	return rpc.Call(params)
}

// EthBlockNumber returns the number of most recent block.
func (rpc *EtherRPC) EthBlockNumber() (int64, error) {
	query := map[string]interface{}{
		"module": "proxy",
		"action": "eth_blockNumber",
	}
	var response string
	if err := rpc.call(query, &response); err != nil {
		return 0, err
	}
	return utils.ParseInt64(response)
}

// EthGetTransactionCount returns the number of transactions sent from an address.
// func (rpc *EthRPC) EthGetTransactionCount(address, block string) (int, error) {
// 	var response string

// 	if err := rpc.call("eth_getTransactionCount", &response, address, block); err != nil {
// 		return 0, err
// 	}

// 	return ParseInt(response)
// }

// // EthSendTransaction creates new message call transaction or a contract creation, if the data field contains code.
// func (rpc *EthRPC) EthSendTransaction(transaction T) (string, error) {
// 	var hash string

// 	err := rpc.call("eth_sendTransaction", &hash, transaction)
// 	return hash, err
// }

// EthEstimateGas makes a call or transaction, which won't be added to the blockchain and returns the used gas, which can be used for estimating the used gas.
// func (rpc *EthRPC) EthEstimateGas(transaction T) (int, error) {
// 	var response string

// 	err := rpc.call("eth_estimateGas", &response, transaction)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return utils.ParseInt(response)
// }

func (rpc *EtherRPC) getBlock(method string, withTransactions bool, params string) (*Block, error) {
	trans := "false"
	if withTransactions {
		trans = "true"
	}
	query := map[string]interface{}{
		"module":  "proxy",
		"action":  "eth_getBlockByNumber",
		"tag":     params,
		"boolean": trans,
	}
	result, err := rpc.RawCall(query)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(result, []byte("null")) {
		return nil, nil
	}

	var response proxyBlock
	if withTransactions {
		response = new(proxyBlockWithTransactions)
	} else {
		response = new(proxyBlockWithoutTransactions)
	}

	err = json.Unmarshal(result, response)
	if err != nil {
		return nil, err
	}

	block := response.toBlock()
	return &block, nil
}

// EthGetBlockByNumber returns information about a block by block number.
func (rpc *EtherRPC) EthGetBlockByNumber(number int64, withTransactions bool) (*Block, error) {

	return rpc.getBlock("eth_getBlockByNumber", withTransactions, utils.Int64ToHex(number))
}

// EthGetBlockTransactionCountByNumber returns the number about a blockheight.
// func (rpc *EthRPC) EthGetBlockTransactionCountByNumber(number int) (big.Int, error) {
// 	var response string
// 	if err := rpc.call("eth_getBlockTransactionCountByNumber", &response, number); err != nil {
// 		return big.Int{}, err
// 	}
// 	return utils.ParseBigInt(response)
// }

// EthGetBlockByHash returns information about a block by hash.
func (rpc *EtherRPC) EthGetBlockByHash(hash string, withTransactions bool) (*Block, error) {
	return rpc.getBlock("eth_getBlockByHash", withTransactions, hash)
}

// func (rpc *EthRPC) getTransaction(method string, params ...interface{}) (*Transaction, error) {
// 	transaction := new(Transaction)

// 	err := rpc.call(method, transaction, params...)
// 	return transaction, err
// }

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
// func (rpc *EthRPC) EthGetTransactionByHash(hash string) (*Transaction, error) {
// 	return rpc.getTransaction("eth_getTransactionByHash", hash)
// }

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
// func (rpc *EthRPC) EthGetTransactionByBlockNumberAndIndex(hash string) (*Transaction, error) {
// 	return rpc.getTransaction("eth_getTransactionByBlockNumberAndIndex", hash)
// }

// EthNewFilter creates a new filter object.
// func (rpc *EthRPC) EthNewFilter(params FilterParams) (string, error) {
// 	var filterID string
// 	err := rpc.call("eth_newFilter", &filterID, params)
// 	return filterID, err
// }

// EthNewBlockFilter creates a filter in the node, to notify when a new block arrives.
// To check if the state has changed, call EthGetFilterChanges.
// func (rpc *EthRPC) EthNewBlockFilter() (string, error) {
// 	var filterID string
// 	err := rpc.call("eth_newBlockFilter", &filterID)
// 	return filterID, err
// }

// EthNewPendingTransactionFilter creates a filter in the node, to notify when new pending transactions arrive.
// To check if the state has changed, call EthGetFilterChanges.
// func (rpc *EthRPC) EthNewPendingTransactionFilter() (string, error) {
// 	var filterID string
// 	err := rpc.call("eth_newPendingTransactionFilter", &filterID)
// 	return filterID, err
// }

// EthUninstallFilter uninstalls a filter with given id.
// func (rpc *EthRPC) EthUninstallFilter(filterID string) (bool, error) {
// 	var res bool
// 	err := rpc.call("eth_uninstallFilter", &res, filterID)
// 	return res, err
// }

// // EthGetFilterChanges polling method for a filter, which returns an array of logs which occurred since last poll.
// func (rpc *EthRPC) EthGetFilterChanges(filterID string) ([]Log, error) {
// 	var logs = []Log{}
// 	err := rpc.call("eth_getFilterChanges", &logs, filterID)
// 	return logs, err
// }

// EthGetFilterLogs returns an array of all logs matching filter with given id.
// func (rpc *EthRPC) EthGetFilterLogs(filterID string) ([]Log, error) {
// 	var logs = []Log{}
// 	err := rpc.call("eth_getFilterLogs", &logs, filterID)
// 	return logs, err
// }

// EthGetLogs returns an array of all logs matching a given filter object.
// func (rpc *EthRPC) EthGetLogs(params FilterParams) ([]Log, error) {
// 	var logs = []Log{}
// 	err := rpc.call("eth_getLogs", &logs, params)
// 	return logs, err
// }

// Eth1 returns 1 ethereum value (10^18 wei)
func (rpc *EtherRPC) Eth1() *big.Int {
	return Ether1()
}

// Eth1 returns 1 ethereum value (10^18 wei)
func Ether1() *big.Int {
	return big.NewInt(1000000000000000000)
}
