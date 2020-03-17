package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/zhxx123/gomonitor/service/wallet/utils"
)

// EthError - ethereum error
type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err EthError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// EthRPC - Ethereum rpc client
type EthRPC struct {
	url    string
	client httpClient
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func New(url string, options ...func(rpc *EthRPC)) *EthRPC {
	rpc := &EthRPC{
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
func NewEthRPC(url string, options ...func(rpc *EthRPC)) *EthRPC {
	return New(url, options...)
}

func (rpc *EthRPC) call(method string, target interface{}, params ...interface{}) error {
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
func (rpc *EthRPC) URL() string {
	return rpc.url
}

// Call returns raw response of method call
func (rpc *EthRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
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
		rpc.log.Println(fmt.Sprintf("%s\nRequest: %s\nResponse: %s\n", method, body, data))
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *EthRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

// NetVersion returns the current network protocol version.
func (rpc *EthRPC) NetVersion() (string, error) {
	var version string

	err := rpc.call("net_version", &version)
	return version, err
}

// NetPeerCount returns number of peers currently connected to the client.
func (rpc *EthRPC) NetPeerCount() (int, error) {
	var response string
	if err := rpc.call("net_peerCount", &response); err != nil {
		return 0, err
	}

	return utils.ParseInt(response)
}

// EthSyncing returns an object with data about the sync status or false.
func (rpc *EthRPC) EthSyncing() (*Syncing, error) {
	result, err := rpc.RawCall("eth_syncing")
	if err != nil {
		return nil, err
	}
	syncing := new(Syncing)
	if bytes.Equal(result, []byte("false")) {
		return syncing, nil
	}
	err = json.Unmarshal(result, syncing)
	return syncing, err
}

// EthGasPrice returns the current price per gas in wei.
func (rpc *EthRPC) EthGasPrice() (big.Int, error) {
	var response string
	if err := rpc.call("eth_gasPrice", &response); err != nil {
		return big.Int{}, err
	}

	return utils.ParseBigInt(response)
}

// EthAccounts returns a list of addresses owned by client.
func (rpc *EthRPC) EthAccounts() ([]string, error) {
	accounts := []string{}

	err := rpc.call("eth_accounts", &accounts)
	return accounts, err
}

// PersonalNewAccount returns a account address and encrypt with password
func (rpc *EthRPC) PersonalNewAccount(password string) (string, error) {
	var account string

	err := rpc.call("personal_newAccount", &account, password)
	return account, err
}

// EthBlockNumber returns the number of most recent block.
func (rpc *EthRPC) EthBlockNumber() (int, error) {
	var response string
	if err := rpc.call("eth_blockNumber", &response); err != nil {
		return 0, err
	}

	return utils.ParseInt(response)
}

// EthGetBalance returns the balance of the account of given address in wei.
func (rpc *EthRPC) EthGetBalance(address, block string) (big.Int, error) {
	var response string
	if err := rpc.call("eth_getBalance", &response, address, block); err != nil {
		return big.Int{}, err
	}

	return utils.ParseBigInt(response)
}

// EthGetTransactionCount returns the number of transactions sent from an address.
// func (rpc *EthRPC) EthGetTransactionCount(address, block string) (int, error) {
// 	var response string

// 	if err := rpc.call("eth_getTransactionCount", &response, address, block); err != nil {
// 		return 0, err
// 	}

// 	return ParseInt(response)
// }

// EthSendTransaction creates new message call transaction or a contract creation, if the data field contains code.
func (rpc *EthRPC) EthSendTransaction(transaction T) (string, error) {
	var hash string

	err := rpc.call("eth_sendTransaction", &hash, transaction)
	return hash, err
}

// EthEstimateGas makes a call or transaction, which won't be added to the blockchain and returns the used gas, which can be used for estimating the used gas.
func (rpc *EthRPC) EthEstimateGas(transaction T) (int, error) {
	var response string

	err := rpc.call("eth_estimateGas", &response, transaction)
	if err != nil {
		return 0, err
	}

	return utils.ParseInt(response)
}

func (rpc *EthRPC) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	result, err := rpc.RawCall(method, params...)
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

// EthGetBlockTransactionCountByNumber returns the number about a blockheight.
func (rpc *EthRPC) EthGetBlockTransactionCountByNumber(number int) (big.Int, error) {
	var response string
	if err := rpc.call("eth_getBlockTransactionCountByNumber", &response, number); err != nil {
		return big.Int{}, err
	}
	return utils.ParseBigInt(response)
}

// EthGetBlockByHash returns information about a block by hash.
func (rpc *EthRPC) EthGetBlockByHash(hash string, withTransactions bool) (*Block, error) {
	return rpc.getBlock("eth_getBlockByHash", withTransactions, hash, withTransactions)
}

// EthGetBlockByNumber returns information about a block by block number.
func (rpc *EthRPC) EthGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	return rpc.getBlock("eth_getBlockByNumber", withTransactions, utils.IntToHex(number), withTransactions)
}

func (rpc *EthRPC) getTransaction(method string, params ...interface{}) (*Transaction, error) {
	transaction := new(Transaction)

	err := rpc.call(method, transaction, params...)
	return transaction, err
}

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (rpc *EthRPC) EthGetTransactionByHash(hash string) (*Transaction, error) {
	return rpc.getTransaction("eth_getTransactionByHash", hash)
}

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (rpc *EthRPC) EthGetTransactionByBlockNumberAndIndex(hash string) (*Transaction, error) {
	return rpc.getTransaction("eth_getTransactionByBlockNumberAndIndex", hash)
}

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
func (rpc *EthRPC) Eth1() *big.Int {
	return Eth1()
}

// Eth1 returns 1 ethereum value (10^18 wei)
func Eth1() *big.Int {
	return big.NewInt(1000000000000000000)
}
