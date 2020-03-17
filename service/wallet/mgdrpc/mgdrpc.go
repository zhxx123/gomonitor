package mgdrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// MgdError - massgrid error
type MgdError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err MgdError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type mgdResponse struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *MgdError       `json:"error"`
}

type mgdRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// MgdRPC - Mgdereum rpc client
type MgdRPC struct {
	url    string
	user   string
	passwd string
	client *http.Client
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func New(url, user, passwd string, options ...func(rpc *MgdRPC)) *MgdRPC {
	rpc := &MgdRPC{
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

// NewMgdRPC create new rpc client with given url
func NewMgdRPC(url, user, passwd string, options ...func(rpc *MgdRPC)) *MgdRPC {
	return New(url, user, passwd, options...)
}

func (rpc *MgdRPC) call(method string, target interface{}, params ...interface{}) error {
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
func (rpc *MgdRPC) URL() string {
	return rpc.url
}

// Call returns raw response of method call
func (rpc *MgdRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := mgdRequest{
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
		rpc.log.Println(fmt.Sprintf("%s\nHeader:%v\nRequest: %s\nResponse: %s\n", method, req.Header, body, data))
	}

	resp := new(mgdResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *MgdRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

// *** The above is the latest block information.

// NetVersion returns the current network protocol version.
func (rpc *MgdRPC) GetBlockChainInfo() (*BlockChainInfo, error) {

	result, err := rpc.RawCall("getblockchaininfo")
	if err != nil {
		return nil, err
	}
	blockChainInfo := new(BlockChainInfo)
	if bytes.Equal(result, []byte("false")) {
		return blockChainInfo, nil
	}
	err = json.Unmarshal(result, blockChainInfo)
	return blockChainInfo, err
}

func (rpc *MgdRPC) GetListTransactions(count, from, confirm int) ([]ListTransactions, error) {
	result, err := rpc.RawCall("listtransactions", "*", count, from, true, confirm)
	if err != nil {
		return nil, err
	}
	var listTransactions []ListTransactions

	if bytes.Equal(result, []byte("[]")) {
		return listTransactions, nil
	}
	err = json.Unmarshal(result, &listTransactions)
	return listTransactions, err
}

func (rpc *MgdRPC) GetTransaction(txid string) (*Transaction, error) {
	result, err := rpc.RawCall("gettransaction", txid)
	if err != nil {
		return nil, err
	}
	transaction := new(Transaction)
	if bytes.Equal(result, []byte("")) {
		return transaction, nil
	}
	err = json.Unmarshal(result, transaction)
	return transaction, err
}

func (rpc *MgdRPC) GetInfo() (*WalletInfo, error) {
	result, err := rpc.RawCall("getinfo")
	if err != nil {
		return nil, err
	}
	walletInfo := new(WalletInfo)
	if bytes.Equal(result, []byte("")) {
		return walletInfo, nil
	}
	err = json.Unmarshal(result, walletInfo)
	return walletInfo, err
}

func (rpc *MgdRPC) GetMiningInfo() (*MiningInfo, error) {
	result, err := rpc.RawCall("getmininginfo")
	if err != nil {
		return nil, err
	}
	walletInfo := new(MiningInfo)
	if bytes.Equal(result, []byte("")) {
		return walletInfo, nil
	}
	err = json.Unmarshal(result, walletInfo)
	return walletInfo, err
}

func (rpc *MgdRPC) GetBalance() (string, error) {
	result, err := rpc.RawCall("getbalance")
	if err != nil {
		return "-1", err
	}
	if bytes.Equal(result, []byte("")) {
		return "-1", nil
	}
	balance := string(result)
	return balance, err
}
