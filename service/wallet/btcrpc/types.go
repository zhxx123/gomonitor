package btcrpc

import (
	"bytes"
	"math/big"
	"net/http"

	"github.com/zhxx123/gomonitor/service/wallet/utils"
)

// type httpClient interface {
// 	Post(url string, contentType string, body io.Reader) (*http.Response, error)
// }

type logger interface {
	Println(v ...interface{})
}

// WithHttpClient set custom http client
func WithBTCHttpClient(client *http.Client) func(rpc *BtcRPC) {
	return func(rpc *BtcRPC) {
		rpc.client = client
	}
}

// WithLogger set custom logger
func WithBTCLogger(l logger) func(rpc *BtcRPC) {
	return func(rpc *BtcRPC) {
		rpc.log = l
	}
}

// WithDebug set debug flag
func WithBTCDebug(enabled bool) func(rpc *BtcRPC) {
	return func(rpc *BtcRPC) {
		rpc.Debug = enabled
	}
}

type hexInt int

func (i *hexInt) UnmarshalJSON(data []byte) error {
	result, err := utils.ParseInt(string(bytes.Trim(data, `"`)))
	*i = hexInt(result)

	return err
}

type hexBig big.Int

func (i *hexBig) UnmarshalJSON(data []byte) error {
	result, err := utils.ParseBigInt(string(bytes.Trim(data, `"`)))
	*i = hexBig(result)

	return err
}

// type strFloat string
// func (w *strFloat) UnmarshalJSON(value []byte) error {
// 	*w = strFloat(string(value))
// 	return nil
// }
// BlockChainInfo - object with blockchain data info
type WalletInfo struct {
	Version         int     `json:"version"`
	ProtocolVersion int     `json:"protocolversion"`
	WalletVersion   int     `json:"walletversion"`
	Balance         float64 `json:"balance"` // float64 .8
	Blocks          int64   `json:"blocks"`
	TimeOffset      int64   `json:"timeoffset"`
	Connections     int64   `json:"connections"`
	Proxy           string  `json:"proxy"`
	Difficulty      float64 `json:"difficulty"` // float 64 .8
	Testnet         bool    `json:"testnet"`
	KeypoololDest   int64   `json:"keypoololdest"`
	KeypoolSize     int     `json:"keypoolsize"`
	PayTxFee        float64 `json:"paytxfee"` // float64 .8
	RelayFee        float64 `json:"relayfee"` // float64 .8
	Erros           string  `json:"errors"`
}

type MiningInfo struct {
	Blocks           int64   `json:"blocks"`
	CurrentBlockSize int64   `json:"currentblocksize"`
	CurrentBlockTx   int64   `json:"currentblocktx"`
	Difficulty       float64 `json:"difficulty"` // float64
	Erros            string  `json:"errors"`
	GenprocLimit     int     `json:"genproclimit"`
	NetworkHashps    float64 `json:"networkhashps"` // float64
	PoolDtx          int64   `json:"pooledtx"`
	Testnet          bool    `json:"testnet"`
	Chain            string  `json:"chain"`
	Generate         bool    `json:"generate"`
}
type BlockChainInfo struct {
	Chain      string `json:"chain"`
	Blocks     int    `json:"blocks"`
	Headers    int    `json:"headers"`
	Mediantime int    `json:"mediantime"`
}
type ListTransactions struct {
	Abandoned         bool     `json:"abandoned"`
	Account           string   `json:"account"`
	Address           string   `json:"address,omitempty"`
	Amount            float64  `json:"amount"`
	BIP125Replaceable string   `json:"bip125-replaceable,omitempty"`
	BlockHash         string   `json:"blockhash,omitempty"`
	BlockIndex        *int64   `json:"blockindex,omitempty"`
	BlockTime         int64    `json:"blocktime,omitempty"`
	Category          string   `json:"category"`
	Confirmations     int64    `json:"confirmations"`
	Fee               *float64 `json:"fee,omitempty"`
	Generated         bool     `json:"generated,omitempty"`
	InvolvesWatchOnly bool     `json:"involveswatchonly,omitempty"`
	Time              int64    `json:"time"`
	TimeReceived      int64    `json:"timereceived"`
	Trusted           bool     `json:"trusted"`
	TxID              string   `json:"txid"`
	Vout              uint32   `json:"vout"`
	WalletConflicts   []string `json:"walletconflicts"`
	Comment           string   `json:"comment,omitempty"`
	OtherAccount      string   `json:"otheraccount,omitempty"`
}

type TransactionDetails struct {
	Account           string   `json:"account"`
	Address           string   `json:"address,omitempty"`
	Amount            float64  `json:"amount"`
	Category          string   `json:"category"`
	InvolvesWatchOnly bool     `json:"involveswatchonly,omitempty"`
	Fee               *float64 `json:"fee,omitempty"`
	Vout              uint32   `json:"vout"`
}

// Transaction models the data from the gettransaction command.
type Transaction struct {
	Amount          float64              `json:"amount"`
	Fee             float64              `json:"fee,omitempty"`
	Confirmations   int64                `json:"confirmations"`
	BlockHash       string               `json:"blockhash"`
	BlockIndex      int64                `json:"blockindex"`
	BlockTime       int64                `json:"blocktime"`
	TxID            string               `json:"txid"`
	WalletConflicts []string             `json:"walletconflicts"`
	Time            int64                `json:"time"`
	TimeReceived    int64                `json:"timereceived"`
	Details         []TransactionDetails `json:"details"`
	Hex             string               `json:"hex"`
}
