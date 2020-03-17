package model

import "github.com/jinzhu/gorm"

const (
	PAYS_TRADE_STATUS_SUCCESS        = "SUCCESS"        //（交易支付成功）
	PAYS_TRADE_STATUS_FAILED         = "FAILED"         // 订单失败
	PAYS_TRADE_STATUS_WAIT_BUYER_PAY = "WAIT_BUYER_PAY" //（交易创建，等待买家付款）
	PAYS_TRADE_STATUS_CLOSED         = "CLOSED"         //（未付款交易超时关闭，或支付完成后全额退款）
	PAYS_TRADE_STATUS_EXPIRED        = "EXPIRE"         // (已过期)
	PAYS_TRADE_STATUS_FINISHED       = "TRADE_FINISHED" //（交易结束，不可退款）
	PAYS_TRADE_STATUS_NOT_PAY        = "NOT_PAY"
)
const (
	TRADE_SUCCESS  = 1001
	TRADE_FAILED   = 1002
	TRADE_WAIT_PAY = 1003
	TRADE_CLOSED   = 1004
	TRADE_EXPIRED  = 1005
	TRADE_FINISHED = 1006
	TRADE_NOT_PAY  = 1007
)
const (
	PAYS_TRADE_TYPE_INPUT  = 1
	PAYS_TRADE_TYPE_OUTPUT = 2
)
const (
	PAYS_ORDER_TYPE_CLOUD = 1
	PAYS_ORDER_TYPE_MINER = 2
)
const (
	PAYS_TRADE_CLOSE    = "close"    // 关闭
	PAY_TRADE_CANCEL    = "cancle"   // 取消
	PAYS_TRADE_REFUND   = "refund"   // 退款
	PAYS_TRADE_RECHARGE = "recharge" // 充值
)

const (
	ALI_PAY_WAIT_BUYER_PAY = "WAIT_BUYER_PAY"
	ALI_PAY_TRADE_CLOSED   = "TRADE_CLOSED"
	ALI_PAY_TRADE_SUCCESS  = "TRADE_SUCCESS"
	ALI_PAY_TRADE_FINISHED = "TRADE_FINISHED"
)

// 用户资产流水记录表
type UserAssetflow struct {
	gorm.Model
	UserId      int    `gorm:"not null; default 0; type:int(10)" json:"user_id"`           //用户id
	OutTradeNo  string `gorm:"not null; default ''; type:varchar(64)" json:"out_trade_no"` //系统生成订单号
	TradeType   int    `gorm:"not null; default 0; type:int(1)" json:"trade_type"`         //交易类型1:收入2支出
	CreateAt    int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"create_at"`      //创建时间
	CoinType    string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`    //充值资产类型，1:cny;2:mgd;3:eth4:btc
	Amount      string `gorm:"not null; default ''; type:varchar(30)" json:"amount"`       //账户金额
	TotalAmount string `gorm:"not null; default ''; type:varchar(30)" json:"total_amount"` //当前账户金额,充值之后
	Description string `gorm:"not null; default ''; type:varchar(64)" json:"description"`  //充值字符串描述
}

// PayTx 法币交易 数据库
type PayTx struct {
	gorm.Model
	PayChannel    int    `gorm:"not null; default 0; type:int(5)" json:"pay_channel"`
	PayType       int    `gorm:"not null; default 0; type:int(5)" json:"pay_type"`
	TradeType     int    `gorm:"not null; default 0; type:int(1)" json:"trade_type"`
	OrderId       string `gorm:"not null; default ''; type:varchar(32)" json:"order_id"`
	UserId        int    `gorm:"not null; default 0; type:int(10)" json:"user_id"`
	TradeNo       string `gorm:"not null; default ''; type:varchar(64)" json:"trade_no"`     //交易ID,唯一
	OutTradeNo    string `gorm:"not null; default ''; type:varchar(30)" json:"out_trade_no"` //数据库中,可能不唯一,有可能有退款,存在相同id,多条记录
	TotalAmount   string `gorm:"not null; default ''; type:varchar(30)" json:"total_amount"`
	RemainAmount  string `gorm:"not null; default ''; type:varchar(30)" json:"remain_amount"`
	RecvTime      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"recv_time"`
	TradeStatus   int    `gorm:"not null; default 0; type:int(10)" json:"trade_status"` //订单状态
	BuyerUserId   string `gorm:"not null; default ''; type:varchar(20)" json:"buyer_user_id"`
	BuyerLogonId  string `gorm:"not null; default ''; type:varchar(64)" json:"buyer_logon_id"`
	SendPayDate   string `gorm:"not null; default ''; type:varchar(32)" json:"send_pay_date"`  //交易时间
	OrderDescript string `gorm:"not null; default ''; type:varchar(64)" json:"order_descript"` //订单描述
	TradeDetail   string `gorm:"not null; default ''; type:text" json:"trade_detail"`
}

// 虚拟充值记录表
type VirtualRecharge struct {
	gorm.Model
	UserId         int    `gorm:"not null; default 0; type:int(10)" json:"user_id"`              //用户id
	OutTradeNo     string `gorm:"not null; default ''; type:varchar(64)" json:"out_trade_no"`    //系统生成订单号
	CreateAt       int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"create_at"`         //创建时间
	CoinType       string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`       //充值资产类型，0:cny;1:mgd;2:eth
	CoinAmount     string `gorm:"not null; default ''; type:varchar(30)" json:"coin_amount"`     //充值之后，账户总金额
	OperatorId     int    `gorm:"not null; default 0; type:int(10)" json:"operator_id"`          //操作员id
	RechargeAmount string `gorm:"not null; default ''; type:varchar(30)" json:"recharge_amount"` //充值金额
	Description    string `gorm:"not null; default ''; type:varchar(64)" json:"description"`     //充值字符串描述
}
type PayIDJson struct {
	// UserId     int    `json:"auth_user_id"` // 服务端自动添加
	PayChannel int    `json:"payChannel" validate:"required,number,min=1,max=10"` // 1 支付渠道,1:支付宝，2:微信
	PayType    int    `json:"payType" validate:"required,number,min=1,max=10"`    // 1 交易类型，1:充值,2:退款
	OrderId    string `json:"orderId" validate:"required,gte=6,lte=36"`           // "xxxx" uuid,32位
}

//PayCreateJson 创建 前端传输数据
type PayCreateJson struct {
	PayChannel  int    `json:"payChannel" validate:"required,number,min=1,max=10"` // 1 支付渠道，1:支付宝，2:微信
	PayType     int    `json:"payType" validate:"required,number,min=1,max=10"`    // 1 支付类型,1:网页支付
	TradeType   int    `json:"tradeType" validate:"required,number,min=1,max=2"`   // 1 交易类型，1:充值,2:退款
	OutTradeNo  string `json:"outTradeNo" validate:"required,gte=6,lte=32"`        // "xxxx" 商家订单号 23位
	OrderId     string `json:"orderId"  validate:"required"`                       // uuid 订单id
	Subject     string `json:"subject" validate:"required"`                        // "账户充值" 订单描述，可选
	TotalAmount string `json:"totalAmount" validate:"required"`                    // "100" 订单金额
	// GoodsId     string `json:"goods_id"`                         // "商品id"
}

// PayQueryJson 查询 前端传输数据
type PayQueryJson struct {
	// PayChannel int    `json:"pay_channel" validate:"required,number,min=0,max=10"`
	// PayType    int    `json:"pay_type" validate:"required,number,min=0,max=10"`
	// UserId int `json:"auth_user_id"` // 服务端自动添加
	// OrderId    string `json:"order_id" validate:"lte=20"`
	OutTradeNo string `json:"outtradeno" validate:"required,gte=10,lte=32"`
}

// PayCancelJson 取消 前端传输数据
type PayCancelJson struct {
	PayChannel int `json:"pay_channel" validate:"required,number,min=1,max=10"`
	PayType    int `json:"pay_type" validate:"required,number,min=1,max=10"`
	// UserId     int `json:"auth_user_id"` //服务端自动添加
	// OrderId    string `json:"order_id" validate:"lte=20"`
	OutTradeNo string `json:"out_trade_no" validate:"required,gte=15,lte=30"`
}

type APayTxJson struct {
	TradeStatus int    `json:"tradeStatus"`
	OutTradeNo  string `json:"outTradeNo"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type AUpatePayTxJson struct {
	TradeStatus int    `json:"tradeStatus"`
	OutTradeNo  string `json:"outTradeNo"`
}

type AVirtualAccountJson struct {
	UserId   int    `json:"userId"`
	CoinType string `json:"coinType"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type AVirtualAssetJson struct {
	UserId         int    `json:"userId"`         // 用户id
	CoinType       string `json:"coinType"`       // 充值资产类型，0: cny; 1: mgd; 2: eth
	OperatorId     int    `json:"operatorId"`     // 操作员 id
	RechargeAmount string `json:"rechargeAmount"` // 充值金额
	Description    string `json:"description"`
}
