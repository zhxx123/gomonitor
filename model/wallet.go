package model

import "github.com/jinzhu/gorm"

type MgdTx struct {
	MyModel
	TxIndex   int    `gorm:"not null; default 0; type:int(10)" json:"tx_index"`
	TxId      string `gorm:"not null; default ''; type:varchar(80)" json:"tx_id"`
	ToAddress string `gorm:"not null; default ''; type:varchar(80)" json:"to_address"`
	TxType    string `gorm:"not null; default ''; type:varchar(10)" json:"tx_type"`
	ToAmount  string `gorm:"not null; default '0'; type:varchar(30)" json:"to_amount"`
	Fee       string `gorm:"not null; default '0'; type:varchar(30)" json:"fee"`
	Time      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"time"`
	ServerIP  string `gorm:"not null; default ''; type:varchar(50)" json:"server_ip"`
	Detail    string `gorm:"not null; default ''; type:varchar(4500)" json:"detail"`
}
type BtcTx struct {
	MyModel
	TxIndex   int    `gorm:"not null; default 0; type:int(10)" json:"tx_index"`
	TxId      string `gorm:"not null; default ''; type:varchar(80)" json:"tx_id"`
	ToAddress string `gorm:"not null; default ''; type:varchar(80)" json:"to_address"`
	TxType    string `gorm:"not null; default ''; type:varchar(10)" json:"tx_type"`
	ToAmount  string `gorm:"not null; default '0'; type:varchar(30)" json:"to_amount"`
	Fee       string `gorm:"not null; default '0'; type:varchar(30)" json:"fee"`
	Time      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"time"`
	ServerIP  string `gorm:"not null; default ''; type:varchar(50)" json:"server_ip"`
	Detail    string `gorm:"not null; default ''; type:varchar(4500)" json:"detail"`
}
type WalletRecord struct {
	gorm.Model
	TxId     string `gorm:"unique;not null; default ''; comment:'交易哈希'; type:varchar(80)"`
	CoinType string `gorm:"not null; default ''; comment:'币种类型'; type:varchar(10)"`
	Added    bool   `gorm:"not null; default 0; comment:'是否添加'; type:tinyint(1)"`
}

// 钱包地址
type WalletAddress struct {
	gorm.Model
	Indexs    int    `gorm:"not null; default 0; comment:'序号'; type:int(10)" json:"indexs"` //这里不能使用index，关键字，否则插入数据库报错
	CoinType  string `gorm:"not null; default ''; comment:'币种类型'; type:varchar(10)" json:"coinType"`
	Account   string `gorm:"not null; default ''; comment:'币种标签'; type:varchar(64)" json:"account"`
	Address   string `gorm:"unique;not null; default ''; comment:'币种地址'; type:varchar(80)" json:"address"`
	UserId    int    `gorm:"not null; default 0; comment:'用户id'; type:int(10)" json:"userId"`
	Allocated bool   `gorm:"not null; default 0; comment:'是否分配'; type:tinyint(1)" json:"allocated"`
}

// 钱包基础信息
type WalletBasic struct {
	gorm.Model
	Name     string `gorm:"not null; default '' comment:'名称'; type:varchar(20)" json:"name"`        // getinfo 名称
	Version  string `gorm:"not null; default 0 comment:'版本号'; type:varchar(10)" json:"version"`     // getinfo 版本号
	NetModel string `gorm:"not null; default '' comment:'网络类型'; type:varchar(10)" json:"net_model"` // getinfo ，test 或者 main getnetworkinfo
	UpdateAt int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"`
}

type WalletSimple struct {
	gorm.Model
	Name        string `gorm:"not null; default ''; comment:'名称'; type:varchar(20)" json:"name"`                 // 名称
	Balance     string `gorm:"not null; default ''; comment:'钱包余额'; type:varchar(32)" json:"balance"`            // getbalance 余额
	BlockHeight int64  `gorm:"not null; default 0; comment:'区块高度'; type:BIGINT(20)" json:"block_height"`         // getmininginfo 区块高度
	Difficulty  string `gorm:"not null; default ''; comment:'网络难度'; type:varchar(32)" json:"difficulty"`         // getmininginfo 难度
	NetworkHash string `gorm:"not null; default ''; comment:'网络hashrate'; type:varchar(32)" json:"network_hash"` // getmininginfo 全网算力
	UpdateAt    int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"`
}

//****** ETH 系列
// eth 同步状态表
type WalletSync struct {
	gorm.Model
	CoinType   string `gorm:"not null;default ''; comment:'币种类型'; type:varchar(10)" json:"coin_type"` // 货币类型
	StartBlock int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"start_block"`                // 开始更新区块高度
	LastBlock  int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"last_block"`                 // 最后更新的区块高度
	UpdateAt   int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"`                  // 更新时间
}
type EthTx struct {
	MyModel
	// Indexs      int    `gorm:"not null default 0;type:int(10)" json:"indexs"` //数据库序号
	BlockHeight int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"block_height"`     // 区块高度
	TxHash      string `gorm:"unique;not null; default ''; type:varchar(70)" json:"tx_hash"` // 交易hash
	TxIndex     int    `gorm:"not null; default 0; type:int(10)" json:"tx_index"`            // 交易序号
	CoinFrom    string `gorm:"not null; default ''; type:varchar(64)" json:"coin_from"`      // 交易支出地址
	CoinTo      string `gorm:"not null; default ''; type:varchar(64)" json:"coin_to"`        // 交易收入地址
	TxType      string `gorm:"not null; default ''; type:varchar(64)" json:"tx_type"`        // 交易类型
	Nonce       int    `gorm:"not null; default 0; type:int(12)" json:"nonce"`               // nonce
	Value       string `gorm:"not null; default '0'; type:varchar(30)" json:"value"`         // 交易值
	Gas         int    `gorm:"not null; default '0'; type:int(10)" json:"gas"`               // 花费 gas
	GasPrice    string `gorm:"not null; default '0'; type:varchar(30)" json:"gas_price"`     // gas 价格
	ServerIP    string `gorm:"not null; default ''; type:varchar(30)" json:"server_ip"`
	Time        int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"time"` // 时间
}
type CoinMarket struct {
	MyModel
	Name  string `gorm:"not null; default ''; type:varchar(12)" json:"name"`
	Url   string `gorm:"not null; default ''; type:varchar(32)" json:"url"`
	Price string `gorm:"not null; default ''; type:varchar(32)" json:"price"`
	Time  int64  `gorm:"not null; default  0; type:BIGINT(20)" json:"time"`
}
type CoinPrice struct {
	gorm.Model
	Name       string `gorm:"not null; default ''; type:varchar(12)" json:"name"`
	Price      string `gorm:"not null; default ''; type:varchar(32)" json:"price"`
	Discount   string `gorm:"not null; default ''; type:varchar(12)" json:"discount"`
	AutoUpdate bool   `gorm:"not null; default 1; type:tinyint(1)" json:"auto_update"`
	Time       int64  `gorm:"not null; default  0; typeBIGINT(20)" json:"time"`
}
type CoinPriceRes struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	Time  int64  `json:"time"`
}

type CoinPriceJsonRes struct {
	Name     string `json:"name"`
	Price    string `json:"price"`
	Discount string `json:"discount"`
}

// ********************************************************
// admin
// ********************************************************

type AWalletTxJson struct {
	Txid     string `json:"txId"`
	CoinType string `json:"coinType"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type AWalletRecordJson struct {
	Txid  string `json:"txId"`
	Added bool   `json:"added"`
}

type AWalletInfoJson struct {
	Name  string `json:"name"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
type AWalletBasicDataRes struct {
	Data   *[]WalletBasic `json:"data"`
	Length int            `json:"length"`
}
type AWalletInfoDataRes struct {
	Data   *[]WalletSimple `json:"data"`
	Length int             `json:"length"`
}

// 地址列表查询
type AWalletAddressJson struct {
	CoinType string `json:"coinType"`
	Status   bool   `json:"status"`
	IsStatus bool   `json:"isStatus"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type AUpdateWalletAddressJson struct {
	Address []WalletAddress `json:"address"`
}

//用户地址分配更新
type AUserCoinAddressJson struct {
	UserId   int    `json:"userId"`
	CoinType string `json:"coinType"`
}

type ACoinMarketJson struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ACoinMarketRes struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Price string `json:"price"`
	Time  int64  `json:"time"`
}

type AUpdateCoinPriceJson struct {
	AuthorId  int    `json:"authorId"`
	CoinType  string `json:"name" validate:"required"`
	CoinPrice string `json:"price"`
	Discount  string `json:"discount"`
}

type AUpdateStatusCoinPriceJson struct {
	CoinType string `json:"coinType" validate:"required"`
	Status   bool   `json:"status"`
}
