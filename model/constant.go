package model

// redis相关常量, 为了防止从redis中存取数据时key混乱了，在此集中定义常量来作为各key的名字
const (
	// SendCode 验证码发送
	SendCode = "sendCode"

	// ResetTime 生成重置密码的链接
	ResetTime = "resetTime"

	// LoginUser 用户信息
	LoginUser = "loginUser"

	// 工单缓存
	WorkOrderKey = "workOrder"

	// 商品订单
	ProductOrder = "productOrder"

	// 支付订单
	PayOrder = "payOrder"

	// 商品详情
	GoodsKey = "goods"
)

type MyMap map[string]interface{}
