package gopay

const (
	ALI_PAY = iota + 1
	WECHAT_PAY
)
const (
	ALI_WEB = iota + 1 // 支付宝网页
	ALI_APP            // 支付宝App
)
const (
	WECHAT_WEB          = iota + 1 // 微信网页
	WECHAT_APP                     // 微信App
	WECHAT_MINI_PROGRAM            // 微信小程序
)

//订单号生成格式
const (
	cTimeFormat = "01021504"
)

// PreCreate Charge 预创建支付
type ChargePreCreate struct {
	PayChannel int `json:"pay_channel"` // 支付渠道
	PayType    int `json:"pay_type"`    // 支付类型
	// OrderId              string `json:"order_id"`                  // 自定义订单ID,前端传入,防止重复出单,入单
	OutTradeNo           string `json:"out_trade_no"`              // MGD自定义订单号
	TradeType            int    `json:"trade_type"`                // 业务类型
	UserId               int    `json:"user_id"`                   // 用户注册 ID
	Subject              string `json:"subject"`                   // 订单名称描述
	TotalAmount          string `json:"total_amount"`              // 订单金额
	GoodsId              string `json:"goods_id"`                  // 商品编号
	TimeOutExpress       string `json:"timer_out_express"`         // 订单失效时间
	QrCodeTimeOutExpress string `json:"qr_code_timer_out_express"` // 预付款二维码失效时间
	CallbackURL          string `json:"callback_url,omitempty"`    // 回调路径
	ReturnURL            string `json:"return_url,omitempty"`      // 回调转发路径
	ShowURL              string `json:"show)_url,omitempty"`       // 展示路径
}

// ChargeQuery 查询支付
type ChargeQuery struct {
	PayChannel int    `json:"pay_channel"`  // 支付渠道
	PayType    int    `json:"pay_type"`     // 支付类型
	TradeType  int    `json:"trade_type"`   // 业务类型
	OutTradeNo string `json:"out_trade_no"` // 查询的订单自定义ID
	TradeNo    string `json:"trade_no"`     // 支付宝订单号, OutTradeNo 二选一
	UserId     int    `json:"user_id"`
}

// ChargeCancel 取消支付
type ChargeCancel struct {
	PayChannel int    `json:"pay_channel"`  // 支付渠道
	PayType    int    `json:"pay_type"`     // 支付类型
	OutTradeNo string `json:"out_trade_no"` // 查询的订单自定义ID
	TradeNo    string `json:"trade_no"`     // 支付宝订单号, OutTradeNo 二选一
}

// ChargeReFund 退款
type ChargeRefund struct {
	PayChannel   int    `json:"pay_channel"`   // 支付渠道
	PayType      int    `json:"pay_type"`      // 支付类型
	TradeType    int    `json:"trade_type"`    // 业务类型,此时应该为2, 可选,非全额退款,必填
	OutTradeNo   string `json:"out_trade_no"`  //如果不是全额退款,一定需要支付号,可选,非全额退款,必填
	UserId       int    `json:"user_id"`       // 用户注册 ID,可选,非全额退款,必填
	RefundAmount string `json:"refund_amount"` //退款金额
	RefundReason string `json:"refund_reason"` //退款原因
}

// ChargeClose 返回
type ChargeClose struct {
	PayChannel int    `json:"pay_channel"`  // 支付渠道
	PayType    int    `json:"pay_type"`     // 支付类型
	OutTradeNo string `json:"out_trade_no"` // 查询的订单自定义ID
	TradeNo    string `json:"trade_no"`     // 支付宝订单号, OutTradeNo 二选一
	OperatorId string `json:"operator_id"`  // 卖家端自定义的的操作员 ID,可选
}

// ChargePreCreateRsp 预创建支付返回
type ChargePreCreateRsp struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	SubCode    string `json:"sub_code"`
	SubMsg     string `json:"sub_msg"`
	OutTradeNo string `json:"out_trade_no"`
	QRCode     string `json:"qr_code"`
}

// ChargeQueryRsp 查询订单返回
type ChargeQueryRsp struct {
	Code         string `json:"code"`
	Msg          string `json:"msg"`
	SubCode      string `json:"sub_code"`
	SubMsg       string `json:"sub_msg"`
	TradeNo      string `json:"trade_no"`
	OutTradeNo   string `json:"out_trade_no"`
	BuyerLogonId string `json:"buyer_logon_id"`
	BuyerUserId  string `json:"buyer_user_id"`
	TradeStatus  string `json:"trade_status"`
	TotalAmount  string `json:"total_amount"`
	SendPayDate  string `json:"send_pay_date"`
}

// ChargeCancelRsp 取消订单返回
type ChargeCancelRsp struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	SubCode    string `json:"sub_code"`
	SubMsg     string `json:"sub_msg"`
	TradeNo    string `json:"trade_no"`
	OutTradeNo string `json:"out_trade_no"`
	RetryFlag  string `json:"retry_flag"`
	Action     string `json:"action"`
}

// ChargeRefundRsp 退款订单返回
type ChargeRefundRsp struct {
	Code         string `json:"code"`
	Msg          string `json:"msg"`
	SubCode      string `json:"sub_code"`
	SubMsg       string `json:"sub_msg"`
	TradeNo      string `json:"trade_no"`
	OutTradeNo   string `json:"out_trade_no"`
	BuyerLogonId string `json:"buyer_logon_id"`
	BuyerUserId  string `json:"buyer_user_id"`
	FundChange   string `json:"fund_change"`
	RefundFee    string `json:"refund_fee"`
	GmtRefundPay string `json:"gmt_refund_pay"`
	StoreName    string `json:"store_name"`
}

// ChargeCloseRsp 关闭订单返回
type ChargeCloseRsp struct {
	Code       string `json:"code"`
	Msg        string `json:"msg"`
	SubCode    string `json:"sub_code"`
	SubMsg     string `json:"sub_msg"`
	TradeNo    string `json:"trade_no"`
	OutTradeNo string `json:"out_trade_no"`
}

// //PayCallback 支付返回
// type PayCallback struct {
// 	Origin      string `json:"origin"`
// 	TradeNum    string `json:"trade_num"`
// 	OrderNum    string `json:"order_num"`
// 	CallBackURL string `json:"callback_url"`
// 	Status      int64  `json:"status"`
// }

// // CallbackReturn 回调业务代码时的参数
// type CallbackReturn struct {
// 	IsSucceed     bool   `json:"isSucceed"`
// 	OrderNum      string `json:"orderNum"`
// 	TradeNum      string `json:"tradeNum"`
// 	UserID        string `json:"userID"`
// 	MoneyFee      int64  `json:"moneyFee"`
// 	Sign          string `json:"sign"`
// 	ThirdDiscount int64  `json:"thirdDiscount"`
// }

// // BaseResult 支付结果
// type BaseResult struct {
// 	IsSucceed     bool   // 是否交易成功
// 	TradeNum      string // 交易流水号
// 	MoneyFee      int64  // 支付金额
// 	TradeTime     string // 交易时间
// 	ContractNum   string // 交易单号
// 	UserInfo      string // 支付账号信息(有可能有，有可能没有)
// 	ThirdDiscount int64  // 第三方优惠
// }
