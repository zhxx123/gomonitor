package alipay

// https://docs.open.alipay.com/194/103296#s5 支付宝扫码支付,异步回调通知
type TradeNotification struct {
	NotifyTime     string `json:"notify_time"`      // 通知时间
	NotifyType     string `json:"notify_type"`      // 通知类型
	NotifyId       string `json:"notify_id"`        // 通知校验ID
	SignType       string `json:"sign_type"`        // 签名类型
	Sign           string `json:"sign"`             // 签名
	TradeNo        string `json:"trade_no"`         // 支付宝交易号
	AppId          string `json:"app_id"`           // 开发者的app_id
	OutTradeNo     string `json:"out_trade_no"`     // 商户订单号
	OutBizNo       string `json:"out_biz_no"`       // 商户业务号
	BuyerId        string `json:"buyer_id"`         // 买家支付宝用户号
	BuyerLogonId   string `json:"buyer_logon_id"`   // 买家支付宝账号
	SellerId       string `json:"seller_id"`        // 卖家支付宝用户号
	SellerEmail    string `json:"seller_email"`     // 卖家支付宝账号
	TradeStatus    string `json:"trade_status"`     // 交易状态
	TotalAmount    string `json:"total_amount"`     // 订单金额
	ReceiptAmount  string `json:"receipt_amount"`   // 实收金额
	InvoiceAmount  string `json:"invoice_amount"`   // 开票金额
	BuyerPayAmount string `json:"buyer_pay_amount"` // 付款金额
	PointAmount    string `json:"point_amount"`     // 集分宝金额
	RefundFee      string `json:"refund_fee"`       // 总退款金额
	SendBackFee    string `json:"send_back_fee"`    // 实际退款金额
	Subject        string `json:"subject"`          // 总退款金额
	Body           string `json:"body"`             // 该订单的备注、描述、明细等
	GmtCreate      string `json:"gmt_create"`       // 交易创建时间
	GmtPayment     string `json:"gmt_payment"`      // 交易买家付款时间
	GmtRefund      string `json:"gmt_refund"`       // 交易退款时间
	GmtClose       string `json:"gmt_close"`        // 交易结束时间
	FundBillList   string `json:"fund_bill_list"`   // 支付成功的各个渠道金额信息
}
