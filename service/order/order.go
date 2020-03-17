package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/cache"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/task"
	"github.com/zhxx123/gomonitor/service/wallet"
	"github.com/zhxx123/gomonitor/utils"

	"github.com/sirupsen/logrus"
)

// model.UserFromRedis 从redis中取出用户信息
func OrderFromRedis(outTradeNo string) (model.Orders, error) {
	orderKey := fmt.Sprintf("%s%s", model.ProductOrder, outTradeNo)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	orderBytes, err := redis.Bytes(RedisConn.Do("GET", orderKey))
	if err != nil {
		logrus.Error(err)
		return model.Orders{}, errors.New("订单不存在")
	}
	var order model.Orders
	bytesErr := json.Unmarshal(orderBytes, &order)
	if bytesErr != nil {
		logrus.Error(bytesErr)
		return order, errors.New("订单有误")
	}
	return order, nil
}

// UserToRedis 将用户信息存到redis
func OrderToRedis(order model.Orders) error {
	orderBytes, err := json.Marshal(order)
	if err != nil {
		logrus.Error(err)
		return errors.New("error")
	}
	orderKey := fmt.Sprintf("%s%s", model.ProductOrder, order.OutTradeNo)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	if _, redisErr := RedisConn.Do("SET", orderKey, orderBytes, "EX", config.ServerConfig.OrderMaxAge); redisErr != nil {
		logrus.Errorf("order redis set failed: %s", redisErr.Error())
		return errors.New("error")
	}
	return nil
}

// model.UserFromRedis 从redis中取出用户信息
func GoodsFromRedis(goodsID string) (model.Products, error) {
	orderKey := fmt.Sprintf("%s%s", model.GoodsKey, goodsID)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	goodsBytes, err := redis.Bytes(RedisConn.Do("GET", orderKey))
	if err != nil {
		logrus.Error(err)
		return model.Products{}, errors.New("商品不存在")
	}
	var goods model.Products
	bytesErr := json.Unmarshal(goodsBytes, &goods)
	if bytesErr != nil {
		logrus.Error(bytesErr)
		return goods, errors.New("商品有误")
	}
	return goods, nil
}

// UserToRedis 将用户信息存到redis
func GoodsToRedis(goods model.Products) error {
	goodsBytes, err := json.Marshal(goods)
	if err != nil {
		logrus.Error(err)
		return errors.New("error")
	}
	goodsKey := fmt.Sprintf("%s%s", model.GoodsKey, goods.GoodsId)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	if _, redisErr := RedisConn.Do("SET", goodsKey, goodsBytes, "EX", config.ServerConfig.OrderMaxAge); redisErr != nil {
		logrus.Errorf("goods redis set failed: %s", redisErr.Error())
		return errors.New("error")
	}
	return nil
}

// CheckAndGetOrderID 获取订单号
func CheckAndGetOrderID(ctx iris.Context, preOrderID *model.PreTradeJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var outTradeNo string
	if len(preOrderID.OrderId) <= 0 {
		return nil, "订单号错误", model.STATUS_FAILED
	}
	if userId == 0 {
		return nil, "error", model.STATUS_FAILED
	}
	value, found := cache.OC.Get(preOrderID.OrderId)
	if found {
		outTradeNo = value.(string)
		logrus.Debugf("CheckOrderAndGetCache OrderId: %s found: %t value: %s\n", preOrderID.OrderId, found, outTradeNo)
	}

	if len(outTradeNo) == 0 { // 重新生成订单号
		outTradeNo = utils.GenOrderOutTradeNo(preOrderID.OrderType, model.ORDER_TRADE_TYPE_INPUT, userId)
		// 将用户提交的订和 uuid 存入缓存
		cache.OC.Set(preOrderID.OrderId, outTradeNo, cache.CacheDefaultExpiration)
	}
	response := model.MyMap{
		"out_trade_no": outTradeNo,
	}
	return response, "success", model.STATUS_SUCCESS
}

// CheckAndPreCreate 创建订单
func CheckAndPreCreate(ctx iris.Context, preCreate *model.PreCreateJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var err error
	if len(preCreate.OrderId) <= 0 {
		return "订单号错误", model.STATUS_FAILED
	}
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	if preCreate.TradeType != model.ORDER_TRADE_TYPE_INPUT {
		return "订单异常", model.STATUS_FAILED
	}
	cValue, cFound := cache.OC.Get(preCreate.OrderId) // 订单状态
	if cFound == false || cValue.(string) != preCreate.OutTradeNo {
		return "订单未找到", model.STATUS_FAILED
	}
	err = preCreateOrder(userId, preCreate)
	if err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// CheckAndPayOrder 确认支付订单
func CheckAndPayOrder(ctx iris.Context, orderPayJson *model.PreCreatePayJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	orders, err := CheckPreOrderPayCache(userId, orderPayJson.OutTradeNo)
	if err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 检查订单状态，并支付
	if err := CheckOrderPayList(orders, orderPayJson); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// OrderCreateQuery 查询订单
func OrderCreateQuery(ctx iris.Context, orderQuery *model.OrderQueryJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	// 检验订单号
	if utils.CheckOutTradeNo(0, 0, userId, orderQuery.OutTradeNo) <= 0 {
		return nil, "订单错误", model.STATUS_FAILED
	}
	// 先查找缓存,找不到,查找数据库
	if orderQuery.OutTradeNo == "" {
		return nil, "订单不存在", model.STATUS_FAILED
	}
	// 获取缓存订单信息
	// value, found := cache.OC.Get(orderQuery.OutTradeNo)
	// if found {
	// 	orders := value.(model.Orders)
	// 	if orders.UserId != userId {
	// 		return nil, "订单不存在", model.STATUS_FAILED
	// 	}
	// 	preOrderRes := model.MyMap{
	// 		"order_status": orders.OrderStatus,
	// 	}
	// 	return preOrderRes, "success", model.STATUS_SUCCESS
	// }
	if orders, err := OrderFromRedis(orderQuery.OutTradeNo); err == nil {
		if orders.UserId != userId {
			return nil, "订单不存在", model.STATUS_FAILED
		}
		preOrderRes := model.MyMap{
			"order_status": orders.OrderStatus,
		}
		return preOrderRes, "success", model.STATUS_SUCCESS
	}

	// 从数据库 查询订单状态
	var orders model.Orders
	if err := db.DB.Where("out_trade_no = ? AND user_id = ? ", orderQuery.OutTradeNo, userId).First(&orders).Error; err != nil {
		logrus.Errorf(" error:%s", err)
		return nil, "订单不存在", model.STATUS_FAILED
	}
	// 添加缓存
	// cache.OC.Set(orderQuery.OutTradeNo, orders, cache.CacheDefaultExpiration)
	if err := OrderToRedis(orders); err != nil {
		logrus.Error(err)
	}
	// 返回数据
	preOrderRes := model.MyMap{
		"order_status": orders.OrderStatus,
	}
	return preOrderRes, "success", model.STATUS_SUCCESS
}

// OrderCancel 取消订单
func OrderCancel(ctx iris.Context, orderCancel *model.OrderCancelJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	// 检验订单号
	if utils.CheckOutTradeNo(0, 0, userId, orderCancel.OutTradeNo) <= 0 {
		return "订单错误", model.STATUS_FAILED
	}
	// 先查找缓存,找不到,查找数据库
	if orderCancel.OutTradeNo == "" {
		return "订单不存在", model.STATUS_FAILED
	}
	// 获取缓存订单信息
	// value, found := cache.OC.Get(orderCancel.OutTradeNo)
	// if found {
	// 	orders := value.(model.Orders)
	// 	if orders.UserId != userId {
	// 		return "订单不存在", model.STATUS_FAILED
	// 	}
	// 	if orders.OrderStatus != model.ORDER_WAIT {
	// 		return "订单不可取消", model.STATUS_FAILED
	// 	}
	// 	return CancleOrder(orders.OutTradeNo, orders.UserId)
	// }
	if orders, err := OrderFromRedis(orderCancel.OutTradeNo); err == nil {
		if orders.UserId != userId {
			return "订单不存在", model.STATUS_FAILED
		}
		if orders.OrderStatus != model.ORDER_WAIT {
			return "订单不可取消", model.STATUS_FAILED
		}
		return CancleOrder(orders.OutTradeNo, orders.UserId)
	}

	// 从数据库 查询订单状态
	orders, err := QueryDBOrderList(orderCancel.OutTradeNo, userId)
	if err != nil {
		return "订单不存在", model.STATUS_FAILED
	}
	if orders.OrderStatus != model.ORDER_WAIT {
		return "订单不可取消", model.STATUS_FAILED
	}
	return CancleOrder(orders.OutTradeNo, orders.UserId)
}

// OrderDelete 删除订单
func OrderDelete(ctx iris.Context, orderClose *model.OrderCloseJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	// 从数据库 查询订单状态
	ordres, err := QueryDBOrderList(orderClose.OutTradeNo, userId)
	if err != nil {
		return "error order not found", model.STATUS_FAILED
	}
	if ordres.OrderStatus == model.ORDER_EXPIRED || ordres.OrderStatus == model.ORDER_FINISHED || ordres.OrderStatus == model.ORDER_CLOSE {
		// 当订单状态为 过期,或者已经完成时候,可以删除订单
		if err := db.DB.Where("out_trade_no = ? AND user_id = ? ", orderClose.OutTradeNo, userId).Delete(&model.Orders{}).Error; err != nil {
			return "error", model.STATUS_FAILED
		}
		return "success", model.STATUS_SUCCESS
	}
	return "failed", model.STATUS_FAILED
}

// GetUserOrders 查询用户所有订单
func GetUserOrders(ctx iris.Context, orderListJson *model.OrderListJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	if userId == 0 {
		return nil, "error", model.STATUS_FAILED
	}
	// 从数据库 查询订单状态
	var orderList []model.Orders
	var orderListRes []model.OrderListRes
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(orderListJson.Page, orderListJson.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if orderListJson.StartAt != 0 {
		dbs = dbs.Where("create_at >= ?", orderListJson.StartAt)
	}
	if orderListJson.EndAt != 0 {
		dbs = dbs.Where("create_at < ?", orderListJson.EndAt)
	}
	if orderListJson.PayType != 0 {
		dbs = dbs.Where("pay_type = ?", orderListJson.PayType)
	}
	var orders model.Orders
	if err := dbs.Model(&orders).Where("user_id = ?", userId).Count(&count).
		Offset(offset).Limit(orderListJson.Limit).
		Find(&orderList).Scan(&orderListRes).Error; err != nil {
		logrus.Errorf("GetUserById failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("QueryUserOrderList useId: %d page: %d limit %d offset: %d\n", userId, orderListJson.Page, orderListJson.Limit, offset)

	response := model.MyMap{
		"data":   orderListRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS

}

// 取消订单
func CancleOrder(outTradeNo string, userId int) (string, int) {
	// 当订单状态为 等待付款状态时,方可取消
	orders := new(model.Orders)
	if err := db.DB.Model(orders).Where("out_trade_no = ? AND user_id = ? ", outTradeNo, userId).Update("order_status", model.ORDER_CLOSE).Error; err != nil {
		logrus.Errorf("preCreateCancel service.DB.Model err %s", err)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 校验订单
func CheckOrderPayList(orders *model.Orders, orderPayJson *model.PreCreatePayJson) error {
	if orders.OrderStatus != model.ORDER_WAIT {
		return errors.New("订单异常")
	}
	// 1. 从获取订单支付类型
	coinType := getOrderCoinPayType(orderPayJson.PayType)

	// 2. 展示 CNY->MGD, CNY->ETH  CNY-> BTC 兑换汇率
	coinPriceRes, _, status := wallet.GetCoinPriceWithCoinName(coinType)
	if status != model.STATUS_SUCCESS || coinPriceRes == nil {
		return errors.New("error coin price")
	}
	// preCreate.PayAmount // 表示商品原始总价
	// 3. 订单金额校验
	exPrice := coinPriceRes.Price
	if exPrice != orderPayJson.ExPrice {
		logrus.Errorf("err,cal exprice: %s order exprice: %s", exPrice, orderPayJson.ExPrice)
	}
	calAmount := utils.FloatMutiplyStrPoint(exPrice, orderPayJson.PayAmount)
	realAmount := orders.TotalAmount
	diffRes := greaterOrEqualAmount(calAmount, realAmount)
	logrus.Infof("CheckOrderPayList  userTd: %d coinType: %s orderAmount(cny): %s calAmount:%s -> %s * %s diff: %t", orders.UserId, coinType, realAmount, calAmount, exPrice, orderPayJson.PayAmount, diffRes)
	if diffRes != true {
		return errors.New("订单金额有误")
	}
	// 4. 检验当前用户余额并扣除
	// 获取当前用户余额
	userAccount := new(model.UserAccounts)
	if err := db.DB.Where("user_id = ? AND coin_type = ? ", orders.UserId, coinType).First(&userAccount).Error; err != nil {
		logrus.Errorf("CreateService DB get user accountAmount err %s", err)
		return errors.New("get db userAccount error")
	}
	// 扣除余额,并创建服务
	remainAmount, err := utils.PayStringSub(userAccount.CoinAmount, orderPayJson.PayAmount, 2)
	if err != nil {
		logrus.Errorf("CreateService PayStringSub err %s , userAmount: %s", err, remainAmount)
		return errors.New("余额不足")
	}

	// 资产流水记录
	createAt := utils.GetNowTime()
	userAssetflow := &model.UserAssetflow{
		UserId:      orders.UserId,
		OutTradeNo:  orders.OutTradeNo,
		TradeType:   1,
		CreateAt:    createAt,
		CoinType:    coinType,
		Amount:      orderPayJson.PayAmount, // 用户支付金额
		TotalAmount: remainAmount,
		Description: "用户购买",
	}
	// ************************************* 事务操作
	tx := db.DB.Begin()
	// 扣除余额
	if err := tx.Model(&userAccount).Where("user_id = ? AND coin_type = ?", orders.UserId, coinType).Update("coin_amount", remainAmount).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("CreateService DB.Model update DB  remain Amount %s, err %s", remainAmount, err)
		return errors.New("error")
	}
	// 添加资产流水
	if err := tx.Create(userAssetflow).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateWalletRecord service.DB.Create err %s, coinType: %s, userid: %d", err.Error(), userAssetflow.CoinType, userAssetflow.UserId)
		return errors.New("error")
	}
	// 更新订单状态
	updateData := map[string]interface{}{
		"pay_type":     orderPayJson.PayType,
		"pay_amount":   orderPayJson.PayAmount,
		"pay_exprice":  exPrice,
		"order_status": model.ORDER_PROCESS,
	}
	// 更新订单数据库
	product := new(model.Orders)
	if err := tx.Model(product).Where("user_id = ? AND out_trade_no = ?", orders.UserId, orders.OutTradeNo).Updates(updateData).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateAdminOrderDetailToDB failed err: %s", err.Error())
		return errors.New("error")
	}
	// 减少商品库存
	// // 获取数据库商品信息
	products, err := GetProductDetailFromDB(orders.OrderType, orders.GoodsId)
	if err != nil || products.Quantity < orders.GoodsQuantity {
		tx.Rollback()
		return errors.New("商品不足")
	}
	// 更新
	newQuantity := products.Quantity - orders.GoodsQuantity
	if err := tx.Model(products).Where("goods_type = ? AND goods_id = ?", orders.OrderType, orders.GoodsId).Update("quantity", newQuantity).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateProductsQuantityToDB failed err: %s", err.Error())
		return errors.New("error")
	}
	tx.Commit()

	// ************************************* 以上 事务操作
	// 更新缓存
	orders.PayType = orderPayJson.PayType
	orders.PayAmount = orderPayJson.PayAmount
	orders.ExPrice = exPrice
	orders.OrderStatus = model.ORDER_PROCESS

	// cache.OC.Set(orders.OutTradeNo, orders, cache.CacheDefaultExpiration)
	if err := OrderToRedis(*orders); err != nil {
		logrus.Error("OrderToRedis error", err.Error())
		return err
	}
	// 获取挖矿收益地址 miner username worker pool
	queryMap := map[string]interface{}{
		"user_id": orders.UserId,
		"status":  1,
	}
	var minerAccounts model.MinerAccounts
	if err := db.DB.Where(queryMap).Last(&minerAccounts).Error; err != nil {
		logrus.Errorf("addMinerAddrToDB service.DB.Create err %s", err)
		return err
	}
	minerPool := minerAccounts.MinerPool
	if strings.HasPrefix(minerPool, "stratum+tcp://") {
		minerPool = strings.Replace(minerPool, "stratum+tcp://", "", 1)
	}
	newAddr := strings.Split(minerAccounts.CoinAddr, ".")
	var userName, minerWorker string
	if len(newAddr) == 2 {
		userName = newAddr[0]
		minerWorker = newAddr[1]
	}
	createParams := map[string]interface{}{
		"out_trade_no":   orders.OutTradeNo,
		"farm_id":        orders.FarmID,
		"goods_type":     orders.MinerGoodsType,
		"price":          orders.RentalType, // 矿机时长类型
		"time":           orders.TotalTime,
		"miner_pool":     minerPool,
		"miner_username": userName,
		"miner_worker":   minerWorker,
	}
	timeDuation := 5 * time.Minute
	// 添加定时任务，到期自动关闭订单
	task.AddTimerWithDeadLine(timeDuation, timeDuation*2, createParams, CreateServiceToFarm)

	logrus.Infof("CreateService  userTd: %d coinType: %s update Amount: %s", orders.UserId, coinType, remainAmount)
	return nil
}

// 检查 订单号以及订单状态 缓存
func CheckPreOrderPayCache(userId int, outTradeNo string) (*model.Orders, error) {
	if outTradeNo == "" {
		return nil, errors.New("订单号错误")
	}
	if userId == 0 {
		return nil, errors.New("error")
	}
	cValue, cFound := cache.OC.Get(outTradeNo) // 订单状态
	if cFound {
		res := cValue.(model.Orders)
		// 验证订单所有者
		if res.UserId != userId {
			return nil, errors.New("订单不存在")
		}
		return &res, nil
	}
	// 如果缓存没有，就查找数据库，然后加入缓存
	var orders model.Orders
	if err := db.DB.Where("user_id = ? AND out_trade_no = ?", userId, outTradeNo).First(&orders).Error; err != nil {
		return nil, errors.New("订单不存在")
	}
	// 将订单加入加入缓存
	cache.OC.Set(outTradeNo, orders, cache.CacheDefaultExpiration)
	return &orders, nil
}

// 创建订单
func preCreateOrder(userId int, preCreate *model.PreCreateJson) error {
	// 首先检查用户的提交的矿工地址，以及矿工名
	preCreate.MinerPool = strings.TrimSpace(preCreate.MinerPool)
	preCreate.MinerPool = utils.AvoidXSS(preCreate.MinerPool)
	if ok := strings.HasPrefix(preCreate.MinerPool, "stratum+tcp://"); ok != true {
		return errors.New("矿池地址有误")
	}

	preCreate.MinerPool = strings.TrimSpace(preCreate.MinerPool)
	preCreate.MinerPool = utils.AvoidXSS(preCreate.MinerPool)
	if ind := strings.Index(preCreate.MinerPool, "."); ind == -1 {
		return errors.New("矿工名有误")
	}

	goods, err := CheckOrderList(preCreate)
	if err != nil {
		return err
	}
	// var goodsDetail = make(map[string]interface{})
	// if err := utils.StrToStruct(goods.Details, &goodsDetail); err != nil {
	// 	logrus.Errorf("precreateorder goodsdetail strToStruct err: %s", err.Error())
	// 	return errors.New("error goods cointype")
	// }
	// 更新挖矿地址到数据库
	// coinType := ""
	// if value, ok := goodsDetail["coin_type"]; ok {
	// 	coinType = value.(string)
	// }

	coinType := goods.Details.CoinType

	minerAccounts := &model.MinerAccounts{
		UserId:      userId,
		CoinType:    coinType,
		MinerPool:   preCreate.MinerPool,
		CoinAddr:    preCreate.MinerAddr,
		AddrType:    1, //表示属于用户
		CoinAmount:  "0",
		TotalAmount: "0",
		Status:      1,
	}
	// 更新数据库
	if err := db.DB.Create(minerAccounts).Error; err != nil {
		logrus.Errorf("addMinerAddrToDB service.DB.Create err %s", err)
		return err
	}
	// totalAmount := preCreate.TotalAmount
	// 生成内部订单号
	tradeNo := utils.GenOrderTradeNo(preCreate.OrderType, preCreate.TradeType, userId)
	orderStatus := model.ORDER_WAIT // 未付款
	// orderDetail := util.StructToStr(preCreate.Detail)
	totalTime := GetRentalTimeSeconds(preCreate.RentalType)
	createTime := utils.GetNowTime()
	charge := &model.Orders{
		UserId:         userId,
		PayType:        0,
		PayAmount:      "0",
		ExPrice:        "0",
		OrderType:      preCreate.OrderType,
		TradeType:      preCreate.TradeType,
		TradeNo:        tradeNo,
		OutTradeNo:     preCreate.OutTradeNo,
		TotalAmount:    preCreate.TotalPrice, // 人民币总价
		TotalTime:      totalTime,            // 租用总时长
		CreateAt:       createTime,
		OrderSubject:   preCreate.Description, // 商品描述
		RentalType:     preCreate.RentalType,
		OrderStatus:    orderStatus,
		GoodsId:        preCreate.GoodsId, // 商品id
		GoodsName:      goods.GoodsName,
		GoodsPrice:     preCreate.CurPrice,    //  商品价格
		GoodsQuantity:  preCreate.TotalAmount, // 购买数量
		GoodsUnit:      goods.Unit,            //商品单位
		FarmID:         goods.FarmID,          // 矿场id
		MinerGoodsType: goods.MinerGoodsType,  // 矿场商品标识
	}
	// 创建订单
	if err := db.DB.Create(charge).Error; err != nil {
		logrus.Errorf("createDBPreCreate service.DB.Create err %s", err)
		return err
	}
	// 订单信息写入缓存
	// cache.OC.Set(preCreate.OutTradeNo, charge, cache.CacheDefaultExpiration)
	if err := OrderToRedis(*charge); err != nil {
		logrus.Error("OrderToRedis error", err.Error())
		return err
	}

	logrus.Infof("preCreateOrder UserId: %d orderType: %d", userId, preCreate.OrderType)
	// 更新 商品数量 缓存
	newQuantity := goods.Quantity - preCreate.TotalAmount
	goods.Quantity = newQuantity
	// cache.OC.Set(preCreate.GoodsId, goods, cache.CacheDefaultExpiration*2)
	if err := GoodsToRedis(*goods); err != nil {
		logrus.Error("GoodsToRedis error", err.Error())
	}
	chargeParams := map[string]interface{}{
		"user_id":      userId,
		"out_trade_no": preCreate.OutTradeNo,
	}
	// 这里取消支付主要针对用户还未成功付款,关闭付款链接,实际上订单设置了超时时间,可以自动关闭,因此这里可以取消成功
	// chargeCancel := new(gopay.ChargeCancel)
	timeDuation := 5 * time.Minute
	// 添加定时任务，到期自动关闭订单
	task.AddTimerWithDeadLine(timeDuation, timeDuation*2, chargeParams, OrderForCancelTimer)
	logrus.Debugf("AddTimerWithDeadLine OrderForCancelTimer %ds", timeDuation/1000000000)
	return nil
}

// 定时任务订单回调函数
func OrderForCancelTimer(data map[string]interface{}) bool {
	userId, ok := data["user_id"]
	if ok != true {
		return true
	}
	outTradeNo, ok := data["out_trade_no"]
	if ok != true {
		return true
	}

	// logrus.Debugf("OrderForCancelTimer cancle order %+v", charge)
	res, status := CancleOrder(outTradeNo.(string), userId.(int))
	if status != model.STATUS_SUCCESS {
		logrus.Errorf("OrderForCancelTimer cancle order error: %s", res)
	}
	//只执行一次，因此返回 true
	return true
}

// 查询订单状态
// 支付成功 status 为 true
func QueryDBOrderList(outTradeNo string, userId int) (*model.Orders, error) {
	var ordList model.Orders
	if err := db.DB.Where("out_trade_no = ? AND user_id = ? ", outTradeNo, userId).First(&ordList).Error; err != nil {
		logrus.Errorf("QueryDBPayTx error:%s", err)
		return nil, err
	}
	return &ordList, nil
}

// 校验订单
func CheckOrderList(preCreate *model.PreCreateJson) (*model.Products, error) {
	// 1.检验订单参数,库存数量
	goods, err := GetGoodsDetail(preCreate.OrderType, preCreate.GoodsId)
	if err != nil {
		return nil, err
	}
	// 校验订单库存
	if goods.Quantity < preCreate.TotalAmount {
		// 表示当前商品库存数量小于购买数量
		return nil, errors.New("商品不足，请刷新重试")
	}
	// 校验价格以及数量
	calPrice := getTotalPrice(preCreate)
	totalPrice := preCreate.TotalPrice
	if res := greaterOrEqualAmount(totalPrice, calPrice); res != true {
		return nil, errors.New("价格错误")
	}
	// 检查地址是否填写，并保存数据库
	if preCreate.MinerAddr == "" {
		return nil, errors.New("收益地址错误")
	}
	return goods, nil
}

// 商品总价格
func getTotalPrice(preCreate *model.PreCreateJson) string {
	// 商品价格
	curPrice := preCreate.CurPrice
	// 商品数量
	amount := fmt.Sprintf("%d", preCreate.TotalAmount)

	acPrice := utils.FloatMutiplyStr(curPrice, amount)

	// 租用类型 比例
	rentalRadio := getRadioFromType(preCreate.RentalType)
	acPrice = utils.FloatMutiplyStr(acPrice, rentalRadio)
	return acPrice

}

// 商品详情
func GetGoodsDetail(goodsType int, goodsId string) (*model.Products, error) {
	// 首先查找缓存
	// value, found := cache.OC.Get(goodsId) // 从缓存中查找商品信息
	// if found {
	// 	res := value.(model.Products)
	// 	return &res, nil
	// }
	if value, err := GoodsFromRedis(goodsId); err != nil {
		return &value, nil
	}
	// 缓存没找到 查找数据库
	product, err := GetProductDetailFromDB(goodsType, goodsId)
	if err != nil {
		return nil, errors.New("商品不存在")
	}
	// 将商品信息添加到缓存
	// cache.OC.Set(goodsId, product, cache.CacheDefaultExpiration*2)
	if err := GoodsToRedis(*product); err != nil {
		logrus.Error("GetGoodsDetail GoodsToRedis error", err.Error())
	}
	return product, nil
}

// 获取商品信息
func GetProductDetailFromDB(goodsType int, goodsId string) (*model.Products, error) {
	var product model.Products
	if err := db.DB.Where("goods_type = ? AND goods_id = ?", goodsType, goodsId).First(&product).Error; err != nil {
		logrus.Errorf("UpdateAdminProductToDB failed err: %s", err.Error())
		return nil, err
	}
	// 获取商品详情
	if err := db.DB.Model(&product).Related(&product.Details).Error; err != nil {
		fmt.Printf("where failed")
	}
	return &product, nil
}

// 字符串 金额大小比较
func greaterOrEqualAmount(accountAmount, orderAmount string) bool {
	// 设置精确度为 0.00000001
	var a utils.Accuracy = func() float64 { return 0.00000001 }
	accAmount, err := utils.ParseStrToFloat(accountAmount)
	if err != nil {
		return false
	}
	ordAmount, err := utils.ParseStrToFloat(orderAmount)
	if err != nil {
		return false
	}

	return a.GreaterOrEqual(accAmount, ordAmount)
}

// 获取支付类型
func getOrderCoinPayType(paytype int) string {
	switch paytype {
	case model.COIN_PAYTYPE_MGD:
		return "MGD"
	case model.COIN_PAYTYPE_ETH:
		return "ETH"
	case model.COIN_PAYTYPE_BTC:
		return "BTC"
	default:
		return "CNY"
	}
}

// 获取订单类型折扣比例
func getRadioFromType(types int) string {
	switch types {
	case 1:
		return model.ORDER_RADIO_ONE_YEAR
	case 2:
		return model.ORDER_RADIO_NINETY_DAYS
	default:
		return model.ORSER_RADIO_THIRTY_DAYS
	}
}

// 获取租用秒数
func GetRentalTimeSeconds(types int) int64 {
	switch types {
	case 365:
		return 365 * 24 * 3600
	case 90:
		return 90 * 24 * 3600
	case 30:
		return 30 * 24 * 3600
	case 10:
		return 10 * 24 * 3600
	default:
		return 0
	}
}
