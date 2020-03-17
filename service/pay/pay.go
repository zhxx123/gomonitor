package pay

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/cache"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/pay/gopay"
	"github.com/zhxx123/gomonitor/service/task"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

// CheckAndGetPayID 获取订单号
func CheckAndGetPayID(ctx iris.Context, payID *model.PayIDJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	outTradeNo, err := checkAndGetCache(userId, payID)
	if err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}
	if len(outTradeNo) == 0 {
		outTradeNo := utils.GenOutTradeNo(payID.PayChannel, payID.PayType, model.PAYS_TRADE_TYPE_INPUT, userId)
		// 将用户提交的订和 uuid 存入缓存
		cache.OC.Set(payID.OrderId, outTradeNo, cache.CacheDefaultExpiration)
	}
	response := model.MyMap{
		"out_trade_no": outTradeNo,
	}
	return response, "success", model.STATUS_SUCCESS
}

// CheckAndPayCreate 校验并创建预支付订单
func CheckAndPayCreate(ctx iris.Context, payPreCreate *model.PayCreateJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	if payPreCreate.OrderId == "" {
		return nil, "订单不存在或已过期", model.STATUS_FAILED
	}
	if payPreCreate.OutTradeNo == "" {
		return nil, "订单不存在或已过期", model.STATUS_FAILED
	}

	if userId == 0 {
		return nil, "error", model.STATUS_FAILED
	}
	if payPreCreate.TradeType != model.PAYS_TRADE_TYPE_INPUT {
		return nil, "订单错误", model.STATUS_FAILED
	}
	value, found := cache.OC.Get(payPreCreate.OrderId)
	if found == false || value.(string) != payPreCreate.OutTradeNo {
		return nil, "订单不存在", model.STATUS_FAILED
	}
	// 待定，检查缓存订单状态，只返回订单状态 ！！！
	if _, err := checkPayOrderCache(payPreCreate.OutTradeNo); err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}

	// 创建订单
	// 判断支付类型,如果不是人民币支付,返回失败
	if payPreCreate.PayChannel != model.COIN_CHANNEL_ALI && payPreCreate.PayChannel == model.COIN_CHANNEL_WECHAT {
		return nil, "参数错误", model.STATUS_FAILED
	}
	payCreateRes, err := officeCreate(payPreCreate, userId)
	if err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}
	// 加入缓存,只缓存订单状态
	cache.OC.Set(payPreCreate.OutTradeNo, model.TRADE_WAIT_PAY, cache.CacheDefaultExpiration)

	return payCreateRes, "success", model.STATUS_SUCCESS
}

// CheckAndPayQuery 校验并查询交易
func CheckAndPayQuery(ctx iris.Context, payQuery *model.PayQueryJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	// 检验订单号
	if utils.CheckOutTradeNo(0, 0, userId, payQuery.OutTradeNo) <= 0 {
		return nil, "订单错误", model.STATUS_FAILED
	}
	// 先查找缓存,找不到,查找数据库
	if payQuery.OutTradeNo == "" {
		return nil, "订单不存在", model.STATUS_FAILED
	}
	value, found := cache.OC.Get(payQuery.OutTradeNo)
	if found { //在缓存中找到
		order_status := (value).(int)
		response := model.MyMap{
			"order_status": order_status,
		}
		return response, "success", model.STATUS_SUCCESS
	}
	// 从数据库 查询订单状态, 数据库中存的已经是 int 类型，无需转换
	status, _ := QueryDBPayTx(payQuery.OutTradeNo, userId)
	// if success != true {
	// 	return nil, "failed", model.STATUS_FAILED
	// }
	response := model.MyMap{
		"order_status": status,
	}
	// 添加缓存
	cache.OC.Set(payQuery.OutTradeNo, status, cache.CacheDefaultExpiration)

	return response, "success", model.STATUS_SUCCESS
}

// CheckAndPayCancel 校验并取消交易
func CheckAndPayCancel(ctx iris.Context, payCancel *model.PayCancelJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	// 检验订单号
	if utils.CheckOutTradeNo(payCancel.PayChannel, payCancel.PayType, userId, payCancel.OutTradeNo) != 3 {
		return "订单错误", model.STATUS_FAILED
	}
	//先查找缓存,找不到,查询数据库(可能存在,用户连续多次点击请求取消)
	// 先查找缓存,找不到,查找数据库
	if payCancel.OutTradeNo != "" {
		_, found := cache.OC.Get(payCancel.OutTradeNo)
		if found {
			// order_status := value.(int)
			return "success", model.STATUS_SUCCESS
		}
	}
	// 先添加缓存, 将订单状态设为关闭
	cache.OC.Set(payCancel.OutTradeNo, model.PAYS_TRADE_STATUS_CLOSED, cache.CacheDefaultExpiration)
	// 创建数据库
	chargeCancel := &gopay.ChargeCancel{
		PayChannel: payCancel.PayChannel,
		PayType:    payCancel.PayType,
		OutTradeNo: payCancel.OutTradeNo,
	}
	cancelPayOrder(chargeCancel)
	return "success", model.STATUS_SUCCESS
}

// CheckAndDelatePayOrder 删除支付订单
func CheckAndDelatePayOrder(ctx iris.Context, payCloseJson *model.PayQueryJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	// 检验订单号
	if utils.CheckOutTradeNo(0, 0, userId, payCloseJson.OutTradeNo) <= 1 {
		return "订单错误", model.STATUS_FAILED
	}
	// 从数据库 查询订单状态
	status, success := QueryDBPayTx(payCloseJson.OutTradeNo, userId)
	if success != true {
		return "error", model.STATUS_FAILED
	}
	if status == model.TRADE_EXPIRED || status == model.TRADE_FINISHED || status == model.TRADE_CLOSED {
		// 当订单状态为 过期,或者已经完成时候,可以删除订单
		payTx := new(model.PayTx)
		if err := db.DB.Where("out_trade_no = ? AND user_id = ? ", payCloseJson.OutTradeNo, userId).Delete(payTx).Error; err != nil {
			logrus.Errorf("DeleteDBPayTx error :%s", err)
			return "error", model.STATUS_FAILED
		}
		return "success", model.STATUS_SUCCESS
	}
	return "error", model.STATUS_SUCCESS
}

// 人民币 支付, 生成预创建订单
func officeCreate(payPreCreate *model.PayCreateJson, userID int) (model.MyMap, error) {

	timeoutExpress := config.PayConfig.Timeout //交易超时关闭时间, 默认5分钟之后二维码失效,无法付款
	goodsId := fmt.Sprintf("u%d", userID)
	subject := fmt.Sprintf("MassGrid 云算力-%s", payPreCreate.Subject)
	charge := &gopay.ChargePreCreate{
		PayChannel:           payPreCreate.PayChannel,
		PayType:              payPreCreate.PayType,
		TradeType:            payPreCreate.TradeType,
		OutTradeNo:           payPreCreate.OutTradeNo,
		UserId:               userID,
		Subject:              subject,
		TotalAmount:          payPreCreate.TotalAmount,
		GoodsId:              goodsId,
		TimeOutExpress:       timeoutExpress,
		QrCodeTimeOutExpress: timeoutExpress,
	}
	rsp, err := gopay.DoTradePreCreate(charge)
	if err != nil {
		return nil, err
	}

	orderPaySatus := model.TRADE_WAIT_PAY
	chargeQuery := map[string]interface{}{
		"pay_channel":  charge.PayChannel,
		"pay_type":     charge.PayType,
		"trade_type":   charge.TradeType,
		"user_id":      charge.UserId,
		"out_trade_no": orderPaySatus,
	}
	pollTime := config.PayConfig.PollTime
	if pollTime < 3 || pollTime > 10 {
		pollTime = 5
	}
	timeDuation := time.Duration(pollTime) * time.Second
	logrus.Debugf("officeCreate channel: %d type: %d user_id: %d out_trade_no: %s", charge.PayChannel, charge.PayType, charge.UserId, charge.OutTradeNo)
	// 3s 查询一次，查询 五分钟
	task.AddTimerWithDeadLine(timeDuation, timeDuation*60, chargeQuery, PayChargeQueryTimer)
	payCreateRes := model.MyMap{
		"out_trade_no": rsp.OutTradeNo,
		"express_time": rsp.QRCode,
		"qr_code":      timeoutExpress,
		"order_status": orderPaySatus,
	}
	return payCreateRes, nil
}

// 论询支付订单
func PayChargeQueryTimer(data map[string]interface{}) bool {

	charge := new(gopay.ChargeQuery)
	err := utils.MapToStruct(data, charge)
	if err != nil { // 说明参数有问题，直接返回true，结束轮询
		logrus.Errorf("PayChargeQueryTimer MapToStruct err:%s ", err.Error())
		return true
	}
	return payForQuery(charge)
}

// 查询支付订单状态
func payForQuery(charge *gopay.ChargeQuery) bool {

	var err error
	var rsp *gopay.ChargeQueryRsp

	rsp, err = gopay.DoTradeQuery(charge)
	if err != nil {
		return false
	}
	pay_status := chargeAliPayStatus(rsp.TradeStatus)
	cache_status := chargePayStatus(pay_status)
	rsp.TradeStatus = pay_status
	//更新缓存 前端提交订单id, 状态, 超时时间
	cache.OC.Set(charge.OutTradeNo, cache_status, cache.CacheDefaultExpiration)
	if rsp.TradeStatus == model.PAYS_TRADE_STATUS_SUCCESS { //交易支付成功
		// 写入数据库
		if err := createPayTx(charge, rsp); err != nil {
			logrus.Infof("payForQuery pay failed OutTradeNo: %s  err: %s", charge.OutTradeNo, err.Error())
		}
		logrus.Infof("payForQuery pay succefully OutTradeNo: %s ", charge.OutTradeNo)
		// 如果 err != nil 表示 数据库操作出错,可能需要 发送报警邮件
		return true
	}
	if rsp.TradeStatus == model.PAYS_TRADE_STATUS_FINISHED || rsp.TradeStatus == model.PAYS_TRADE_STATUS_CLOSED {
		logrus.Warnf("payForQuery pay failed OutTradeNo: %s  TradeStatus %s", charge.OutTradeNo, rsp.TradeStatus)
		chargeCancel := &gopay.ChargeCancel{
			PayChannel: charge.PayChannel,
			PayType:    charge.PayType,
			OutTradeNo: charge.OutTradeNo,
		}
		cancelPayOrder(chargeCancel)
		return true
	}
	// 执行到此处表示, 没有查询到数据, 因此直接返回 [可能原因: 用户可能未付款]，继续加入队列，定时查询
	logrus.Infof("payForQuery pay succefully OutTradeNo: %s status: %s -> %d", charge.OutTradeNo, pay_status, cache_status)
	return false
}

// 关闭异常订单
func cancelPayOrder(charge *gopay.ChargeCancel) bool {
	var err error
	var rsp *gopay.ChargeCancelRsp

	rsp, err = gopay.DoTradeCancel(charge)
	if err != nil {
		return false
	}
	if rsp.RetryFlag == "N" { // 取消支付成功,// 撤销交易成功
		logrus.Infof("payForCancelTs pay succefully OutTradeNo: %s", charge.OutTradeNo)
		return true
	}
	if rsp.Action == model.PAYS_TRADE_CLOSE || rsp.Action == model.PAYS_TRADE_REFUND { // 交易失败,已关闭,超时,等等
		logrus.Errorf("payForCancelTs pay failed OutTradeNo: %s  Action %s", charge.OutTradeNo, rsp.Action)
		return true
	}
	// 执行到此处表示, 没有查询到数据, 因此直接返回
	logrus.Warnf("payForCancelTs pay failed OutTradeNo: %s timout or pay_not_found", charge.OutTradeNo)
	return false
}

// 订单查询结果 转换 数据库存储结构
func createPayTx(chargeQuery *gopay.ChargeQuery, chargeQueryRsp *gopay.ChargeQueryRsp) error {
	// payTx := new(model.PayTx)
	nowTime := utils.GetNowTime()
	tradeStatus := chargePayStatus(chargeQueryRsp.TradeStatus)
	tradeDetail := utils.StructToStr(chargeQueryRsp)
	payTx := &model.PayTx{
		PayChannel:    chargeQuery.PayChannel,
		PayType:       chargeQuery.PayType,
		TradeType:     chargeQuery.TradeType,
		UserId:        chargeQuery.UserId,
		TradeNo:       chargeQueryRsp.TradeNo,
		OutTradeNo:    chargeQueryRsp.OutTradeNo,
		TotalAmount:   chargeQueryRsp.TotalAmount,
		RemainAmount:  chargeQueryRsp.TotalAmount,
		RecvTime:      nowTime,
		TradeStatus:   tradeStatus,
		BuyerUserId:   chargeQueryRsp.BuyerUserId,
		BuyerLogonId:  chargeQueryRsp.BuyerLogonId,
		SendPayDate:   chargeQueryRsp.SendPayDate,
		OrderDescript: model.PAYS_TRADE_RECHARGE,
		TradeDetail:   tradeDetail,
	}
	err := CreateDBPayTx(payTx)
	return err
}

// 用户资产添加
func CreateDBPayTx(payTx *model.PayTx) error {
	// 添加用户余额表
	// 获取用户余额表
	coinType := "CNY"
	userAccount, err := GetUserAmount(payTx.UserId, coinType)
	if err != nil {
		logrus.Errorf("getPayReFundAmount GetUserAmount err: %s", err)
		return errors.New("erro for user amount")
	}

	// 用户资产余额增加
	newAccountAmount := utils.PayStringAdd(userAccount.CoinAmount, payTx.TotalAmount, 2)
	// 事务操作
	tx := db.DB.Begin()
	// 新建一条充值记录
	if err := tx.Create(payTx).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("CreateDBPayTx service.DB.Create err %s", err)
		return err
	}
	userAccounts := new(model.UserAccounts)

	// 用户 余额表 需要更新
	if err := tx.Model(&userAccounts).Where("user_id = ? AND coin_type = ?", payTx.UserId, coinType).Update("coin_amount", newAccountAmount).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("CreateDBPayTx service.DB.Model err %s", err)
		return err
	}

	// 添加用户账户流水
	userAssetflow := &model.UserAssetflow{
		UserId:      payTx.UserId,
		OutTradeNo:  payTx.OutTradeNo,
		TradeType:   model.TRADE_TYPE_INPUT,
		CreateAt:    payTx.RecvTime,
		CoinType:    userAccount.CoinType,
		Amount:      payTx.TotalAmount,
		TotalAmount: newAccountAmount,
		Description: "用户充值",
	}
	if err := tx.Create(userAssetflow).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateWalletRecord service.DB.Create err %s, coinType: %s, userid: %d", err.Error(), userAssetflow.CoinType, userAssetflow.UserId)
		return err
	}
	tx.Commit()

	logrus.Infof("CreateDBPayTx AddNew PayTx OutTradeNo: %s, CoinType: %s Now AccountAmount %s ", payTx.OutTradeNo, userAccount.CoinType, newAccountAmount)
	return nil
}

// 返回用户,该币种的余额
func GetUserAmount(userId int, coinType string) (*model.UserAccounts, error) {
	// 获取当前用户余额
	userAccount := new(model.UserAccounts)
	if err := db.DB.Where("user_id = ? AND coin_type = ? ", userId, coinType).First(&userAccount).Error; err != nil {
		logrus.Errorf("getUserAmount service.DB get user accountAmount err %s", err)
		return nil, err
	}
	logrus.Debugf("getUserAmount  userTd: %d coinType: %s  coinAmount: %s", userId, coinType, userAccount.CoinAmount)
	return userAccount, nil
}

// 查询订单状态
func QueryDBPayTx(outTradeNo string, userId int) (int, bool) {
	var ptx model.PayTx
	if err := db.DB.Where("out_trade_no = ? AND user_id = ? ", outTradeNo, userId).First(&ptx).Error; err != nil {
		logrus.Infof("QueryDBPayTx error :%s", err)
		return model.TRADE_NOT_PAY, false
	}
	return ptx.TradeStatus, true
}

// 检查订单缓存
func checkAndGetCache(userId int, payID *model.PayIDJson) (string, error) {
	if payID.OrderId == "" {
		return "", errors.New("订单错误")
	}
	if userId == 0 {
		return "", errors.New("error")
	}
	value, found := cache.OC.Get(payID.OrderId)
	if found { // 当前orderId已经存在,可能是用户重复请求
		logrus.Debugf("checkAndGetCache OrderId: %s found: %t value: %s\n", payID.OrderId, found, value.(string))
		return value.(string), nil
	}
	return "success", nil
}

// 订单缓存
func checkPayOrderCache(outTradeNo string) (int, error) {
	if outTradeNo == "" {
		return 0, errors.New("订单不存在或已过期")
	}
	value, found := cache.OC.Get(outTradeNo)
	if found {
		if status, ok := (value).(int); ok {
			return status, errors.New("订单已经创建")
		}
		return 0, errors.New("订单错误")
	}
	return 0, nil
}

// 支付状态 消息转换
func chargeAliPayStatus(value string) string {
	switch value {
	case model.ALI_PAY_TRADE_SUCCESS: //交易支付成功
		return model.PAYS_TRADE_STATUS_SUCCESS
	case model.ALI_PAY_WAIT_BUYER_PAY: //交易创建，等待买家付款
		return model.PAYS_TRADE_STATUS_WAIT_BUYER_PAY
	case model.ALI_PAY_TRADE_CLOSED: // 未付款交易超时关闭，或支付完成后全额退款
		return model.PAYS_TRADE_STATUS_CLOSED
	case model.ALI_PAY_TRADE_FINISHED: //交易结束，不可退款
		return model.PAYS_TRADE_STATUS_FINISHED
	default: //自定义状态 未支付,或者 未找到
		return model.PAYS_TRADE_STATUS_NOT_PAY
	}
}

// 支付 消息转换, 转换为 int 类型
func chargePayStatus(value interface{}) int {
	switch value {
	case model.PAYS_TRADE_STATUS_SUCCESS: //交易支付成功
		return model.TRADE_SUCCESS
	case model.PAYS_TRADE_STATUS_WAIT_BUYER_PAY: //交易创建，等待买家付款
		return model.TRADE_WAIT_PAY
	case model.PAYS_TRADE_STATUS_CLOSED: // 未付款交易超时关闭，或支付完成后全额退款
		return model.TRADE_CLOSED
	case model.PAYS_TRADE_STATUS_FINISHED: //交易结束，不可退款
		return model.TRADE_FINISHED
	case model.PAYS_TRADE_STATUS_FAILED:
		return model.TRADE_FAILED
	default: //自定义状态 未支付,或者 未找到
		return model.TRADE_NOT_PAY
	}
}
