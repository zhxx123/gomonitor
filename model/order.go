package model

import "github.com/jinzhu/gorm"

const (
	COIN_CHANNEL_ALI     = iota + 1 // 1
	COIN_CHANNEL_WECHAT             // 2
	COIN_CHANNEL_DIGITAL            // 3
	COIN_CHANNEL_VIRTUAL
)
const (
	COIN_PAYTYPE_CNY = iota + 1 // 1
	COIN_PAYTYPE_MGD            // 2
	COIN_PAYTYPE_ETH            // 3
	COIN_PAYTYPE_BTC            // 4
)

const (
	ORDER_STATUS_SUCCESS  = "SUCCESS"   //（订单创建成功）
	ORDER_STATUS_FAILED   = "FAILED"    //（订单失败）
	ORDER_STATUS_WAIT     = "WAIT_PAY"  //（订单未付款）
	ORDER_STATUS_CLOSE    = "CLOSE"     //（订单已关闭）
	ORDER_STATUS_EXPIRED  = "EXPIRE"    // (订单已过期)
	ORDER_STATUS_FINISHED = "FINISH"    //（订单已完成)
	ORDER_STATUS_PROCESS  = "PROCESS"   //（订单处理中）
	ORDER_STATUS_NOTFOUND = "NOT_FOUND" // (订单未找到)
)
const (
	ORDER_SUCCESS  = 1101 //（订单创建成功）
	ORDER_FAILED   = 1102 //（订单失败）
	ORDER_WAIT     = 1103 //（订单未付款）
	ORDER_CLOSE    = 1104 //（订单已关闭）
	ORDER_EXPIRED  = 1105 // (订单已过期)
	ORDER_FINISHED = 1106 //（订单已完成)
	ORDER_PROCESS  = 1107 //（订单处理中）
	ORDER_NOTFOUND = 1108 // (订单未找到)
)

// 订单类型折扣
const (
	ORDER_RADIO_ONE_YEAR    = "0.90"
	ORDER_RADIO_NINETY_DAYS = "0.95"
	ORSER_RADIO_THIRTY_DAYS = "1"
)
const (
	ORDER_TRADE_TYPE_INPUT  = 1
	ORDER_TRADE_TYPE_OUTPUT = 2
)

const (
	MinerOrderUnknown   = 0 // 未知
	MinerOrderCreating  = 1 // 正在创建
	MinerOrderCreated   = 2 // 已运行
	MinerOrderReSetting = 3 // 修改了挖矿配置
	MinerOrderRunning   = 4 // 订单运行中
	MinerOrderCompleted = 5 // 已完成
)

// 订单状态数据库
// Orders
type Orders struct {
	gorm.Model
	UserId int `gorm:"not null; default 0; type:int(10)" json:"user_id"`
	//OrderIdstring`gorm:"notnulldefault'';type:varchar(20)"json:"order_id"`//临时使用，前端用户提交
	PayType        int    `gorm:"not null; default 0; type:int(10)" json:"pay_type"`              //支付类型1:cny2:mgd3:eth4:BTC
	PayAmount      string `gorm:"not null; default ''; type:varchar(30)" json:"pay_amount"`       //支付数量(金额)
	ExPrice        string `gorm:"not null; default ''; type:varchar(30)" json:"ex_price"`         //兑换价格
	OrderType      int    `gorm:"not null; default 0; type:int(1)" json:"order_type"`             //订单类型，1:云算力2:矿机
	TradeType      int    `gorm:"not null; default 0; type:int(5)" json:"trade_type"`             //1:购买，2:退款
	TradeNo        string `gorm:"not null; default ''; type:varchar(64)" json:"trade_no"`         //交易ID,唯一
	OutTradeNo     string `gorm:"not null; default ''; type:varchar(64)" json:"out_trade_no"`     //数据库中,可能不唯一,有可能有退款,存在相同id,多条记录
	TotalAmount    string `gorm:"not null; default ''; type:varchar(30)" json:"total_amount"`     //商品支付总价，人民币
	TotalTime      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"total_time"`         //订单总时长秒计
	RentalType     int    `gorm:"not null; default 0; type:int(10)" json:"rental_type"`           //365:一年90:90天30:30天10:10天
	CreateAt       int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"create_at"`          //创建时间
	OrderSubject   string `gorm:"not null; default ''; type:varchar(32)" json:"order_subject"`    //订单描述
	OrderStatus    int    `gorm:"not null; default 0; type:int(10)" json:"order_status"`          //订单状态、
	GoodsId        string `gorm:"not null; default ''; type:varchar(32)" json:"goods_id"`         //商品编号
	GoodsName      string `gorm:"not null; default ''; type:varchar(32)" json:"goods_name"`       //商品名称
	MinerGoodsType string `gorm:"not null; default ''; type:varchar(32)" json:"miner_goods_type"` //矿场使用的商品ID
	FarmID         string `gorm:"not null; default ''; type:varchar(32)" json:"farm_i_d"`         //矿场的ID
	GoodsPrice     string `gorm:"not null; default ''; type:varchar(30)" json:"goods_price"`      //商品单价
	GoodsQuantity  int    `gorm:"not null; default 0; type:int(10)" json:"goods_quantity"`        //购买数量
	GoodsUnit      string `gorm:"not null; default ''; type:varchar(10)" json:"goods_unit"`       //商品单位
	//OrderDetailstring`gorm:"notnulldefault'';type:text"json:"order_detail"`//订单其他详情
}

// 矿场租用订单详情
type MinerOrder struct {
	gorm.Model
	GoodsID       string `gorm:"not null; default ''; type:varchar(32)" json:"goods_i_d"`       //商品ID
	MinerOrderID  string `gorm:"not null; default ''; type:varchar(64)" json:"miner_order_i_d"` //矿场订单
	FarmID        string `gorm:"not null; default ''; type:varchar(64)" json:"farm_i_d"`        //矿场id
	MinerID       string `gorm:"not null; default ''; type:varchar(64)" json:"miner_i_d"`       //机器id
	CreateAt      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"create_at"`         //创建时间
	UpdateAt      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"`         //更新时间
	Status        int    `gorm:"not null; default 0; type:int(10)" json:"status"`               //当前状态
	GoodsType     string `gorm:"not null; default ''; type:varchar(32)" json:"goods_type"`      //类型
	GoodsPrice    uint64 `gorm:"not null; default ''; type:BIGINT(32)" json:"goods_price"`      //价格
	RentTime      uint64 `gorm:"not null; default ''; type:BIGINT(32)" json:"rent_time"`        //租用时长秒数
	MinerPool     string `gorm:"not null; default ''; type:varchar(128)" json:"miner_pool"`     //矿池
	MinerUsername string `gorm:"not null; default ''; type:varchar(64)" json:"miner_username"`  //矿工用户名
	MinerWorker   string `gorm:"not null; default ''; type:varchar(64)" json:"miner_worker"`    //矿工名
	OutTradeNo    string `gorm:"not null; default ''; type:varchar(80)" json:"out_trade_no"`    //订单号
}

//获取订单列表
type OrderListJson struct {
	StartAt     int64 `json:"startAt"`     //开始时间
	EndAt       int64 `json:"endAt"`       //结束时间
	PayType     int   `json:"payType"`     // 支付类型
	OrderStatus int   `json:"orderStatus"` // 订单状态
	Page        int   `json:"page" validate:"required,number,min=1,max=100"`
	Limit       int   `json:"limit" validate:"required,number,min=1,max=20"`
}
type OrderListRes struct {
	// PayChannel   int    `json:"pay_channel"`
	PayType int `json:"pay_type"`
	// TradeType    int    `json:"trade_type"`
	// OrderId      string `json:"order_id"`
	// TradeNo      string `json:"trade_no"`     // 交易ID,唯一
	OutTradeNo   string `json:"out_trade_no"` // 数据库中,可能不唯一,有可能有退款,存在相同id,多条记录
	TotalAmount  string `json:"total_amount"`
	TotalTime    int64  `json:"total_time"`
	CreateAt     int64  `json:"create_at"`
	OrderSubject string `json:"order_subject"` // 订单描述
	OrderStatus  int    `json:"order_status"`
	OrderDetail  string `json:"order_detail"`
}

type PreTradeJson struct {
	// UserId    int `json:"auth_user_id"`                                      // 服务端自动添加
	OrderType int `json:"orderType" validate:"required,number,min=1,max=2"` // 1: 云算力, 2:矿机
	// PayType   int    `json:"pay_type" validate:"required,number,min=1,max=3"`   // 1: cny,2:mgd,3:eth
	OrderId string `json:"orderId" validate:"required,gte=30,lte=32"`
}
type PreTradeRes struct {
	OutTradeNo string `json:"out_trade_no"` // 生产的订单id号
	// TotalAmount string `json:"total_amount"` // 账户余额
}

type PreCreateJson struct {
	// UserId      int    `json:"auth_user_id"`                                      // 服务端自动添加
	OrderType   int    `json:"orderType" validate:"required,number,min=1,max=2"` // 2 订单类型,1:云算力，2:矿机
	TradeType   int    `json:"tradeType" validate:"required,number,min=1,max=2"` // 1 交易类型，1:购买
	OrderId     string `json:"orderId" validate:"required,gte=6,lte=32"`         // uuid 前端传入订单号，主要是为了查找缓存中生成的订单ID
	OutTradeNo  string `json:"outTradeNo" validate:"required,gte=6,lte=32"`      // "xxxx" 商家订单号 23位
	GoodsId     string `json:"goodsId" validate:"required"`                      // "xxx" 商品ID
	Description string `json:"description" validate:"required"`                  // "xxxx"  商品描述(BTC云算力)
	CurPrice    string `json:"curPrice" validate:"required"`                     // "10.2" 商品当前价格
	TotalAmount int    `json:"totalAmount" validate:"required"`                  // 100 购买数量
	TotalPrice  string `json:"totalPrice" validate:"required"`                   // "100" 商品总价
	RentalType  int    `json:"rentalType" validate:"required"`                   // 365 租用时长类型 365 90 30 10 天数
	MinerPool   string `json:"minerPool" validate:"required"`                    // "xxx" 矿池链接
	MinerAddr   string `json:"minerAddr" validate:"required"`                    // "xxx.xxx" 挖矿收益地址
	// TotalTime   int64  `json:"total_time" validate:"required"`                 // 总时长 （按秒计算）
}

// 商品购买订单确认支付
type PreCreatePayJson struct {
	// UserId     int    `json:"auth_user_id"`                                // 服务端自动添加
	OutTradeNo string `json:"outTradeNo" validate:"required,gte=6,lte=32"`    // "xxxx" 商家订单号 23位
	PayType    int    `json:"payType" validate:"required,number,min=1,max=3"` // 1 实际支付类型,1:CNY,2:MGD,3:ETH,4:BTC
	PayAmount  string `json:"payAmount" validate:"required"`                  // "10.89" 实际支付数量(金额)
	ExPrice    string `json:"exPrice" validate:"required"`                    // "0.89" 兑换 cny 价格
	// MinerAddr  string `json:"miner_addr"` // 挖矿地址
	// MinerPool  string `json:"miner_pool"` // 矿池
}

// GoodsDetail 商品详情
type GoodsDetail struct {
	GoodsType     int    `json:"googs_type"`     // 商品类型，云算力商品，矿机商品...
	GoodsId       string `json:"goods_id"`       // 商品编号
	GoodsName     string `json:"goods_name"`     // 商品名称
	GoodsPrice    string `json:"goods_price"`    // 商品单价
	GoodsQuantity int    `json:"goods_quantity"` // 购买数量
	GoodsUnit     string `json:"goods_unit"`     // 商品单位
}

// // OrderRes resonse
// type OrderRes struct {
// 	OrderStatus int `json:"order_status"`
// }

// OrderQueryJson 订单查询数据结构
type OrderQueryJson struct {
	// PayChannel int    `json:"pay_channel" validate:"required,number,min=0,max=10"`
	// PayType    int    `json:"pay_type" validate:"required,number,min=0,max=10"`
	// UserId int `json:"auth_user_id"` // 服务端自动添加
	// OrderId    string `json:"order_id" validate:"lte=20"`
	OutTradeNo string `json:"out_trade_no" validate:"required,gte=15,lte=30"`
}

// OrderCancelJson 取消 前端传输数据
type OrderCancelJson struct {
	// PayChannel int    `json:"pay_channel" validate:"required,number,min=0,max=10"`
	// PayType    int    `json:"pay_type" validate:"required,number,min=0,max=10"`
	// UserId int `json:"auth_user_id"` // 服务端自动添加
	// OrderId    string `json:"order_id" validate:"lte=20"`
	OutTradeNo string `json:"out_trade_no"`
}

// OrderRefundJson 退款
type OrderRefundJson struct {
	PayChannel   int    `json:"pay_channel" validate:"required,number,min=0,max=10"`
	PayType      int    `json:"pay_type" validate:"required,number,min=1,max=3"`
	TradeType    int    `json:"trade_type" validate:"required,number,min=1,max=2"`
	UserId       int    `json:"auth_user_id"` //服务端自动添加
	OrderId      string `json:"order_id" validate:"lte=20"`
	OutTradeNo   string `json:"out_trade_no" validate:"required,gte=15,lte=30"`
	RefundReason string `json:"refund_reason" validate:"ltc=50"`
}

// OrderCloseJson 关闭 前端传输数据
type OrderCloseJson struct {
	// PayChannel int    `json:"pay_channel" validate:"required,number,min=0,max=10"`
	// PayType    int    `json:"pay_type" validate:"required,number,min=0,max=10"`
	// UserId int `json:"auth_user_id"` //服务端自动添加
	// OrderId    string `json:"order_id" validate:"lte=20"`
	OutTradeNo string `json:"out_trade_no" validate:"required,gte=15,lte=30"`
	// TradeNo    string `json:"trade_no"`
	// OperatorId string `json:"operator_id"`
}

// OrderUpdateJson 订单更新数据结构
type OrderUpdateJson struct {
	// PayChannel  int    `json:"pay_channel" validate:"required,number,min=0,max=10"`
	// PayType     int    `json:"pay_type" validate:"required,number,min=0,max=10"`
	UserId int `json:"auth_user_id"` // 服务端自动添加
	// OrderId     string `json:"order_id" validate:"lte=20"`
	OutTradeNo  string `json:"out_trade_no" validate:"required,gte=15,lte=30"`
	OrderStatus string `gorm:"not null default '';type:varchar(32)"`
}

// miner order
type MinerOrderInfo struct {
	FarmID        string
	GoodsType     string
	Price         uint64
	Time          uint64
	MinerPool     string
	MinerUsername string
	MinerWorker   string
}

// admin
// ********************************************************
// 所有订单列表
type AOrderListJson struct {
	OutTradeNo  string `json:"outTradeNo"`
	OrderType   int    `json:"orderType"`
	OrderStatus int    `json:"orderStatus"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}
type AOrdersRes struct {
	ID            uint   `json:"id"`
	UserId        int    `json:"user_id"`
	PayType       int    `json:"pay_type"`     // 支付类型 1: cny 2:mgd 3:eth
	PayAmount     string `json:"pay_amount"`   // 支付数量(金额)
	ExPrice       string `json:"ex_price"`     // 兑换价格
	OrderType     int    `json:"order_type"`   // 订单类型，1: 云算力 2:矿机
	TradeType     int    `json:"trade_type"`   // 1:购买，2:退款
	TradeNo       string `json:"trade_no"`     // 交易ID,唯一
	OutTradeNo    string `json:"out_trade_no"` // 数据库中,可能不唯一,有可能有退款,存在相同id,多条记录
	TotalAmount   string `json:"total_amount"` // 商品支付总价，人民币
	TotalTime     int64  `json:"total_time"`
	CreateAt      int64  `json:"create_at"`
	OrderSubject  string `json:"order_subject"` // 订单描述
	OrderStatus   int    `json:"order_status"`
	GoodsId       string `json:"goods_id"`       // 商品编号
	GoodsName     string `json:"goods_name"`     // 商品名称
	GoodsPrice    string `json:"goods_price"`    // 商品单价
	GoodsQuantity int    `json:"goods_quantity"` // 购买数量
	GoodsUnit     string `json:"goods_unit"`     // 商品单位
	// Detail       GoodsDetail `json:"order_detail"`
}

// 订单状态
// OrderUpdateJson 订单更新数据结构
type AOrderUpdateStatusJson struct {
	UserId      int    `json:"userId"`
	OutTradeNo  string `json:"outTradeNo" validate:"required,gte=15,lte=30"`
	OrderStatus int    `json:"orderStatus"`
}

// 订单详情
