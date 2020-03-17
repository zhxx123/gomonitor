package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/order/poolrpc"
	"github.com/zhxx123/gomonitor/service/task"
	"github.com/zhxx123/gomonitor/utils"

	"github.com/zhxx123/gomonitor/config"
)

var poolRpc *poolrpc.PoolRPC

// 初始化 poolrpc
func InitPoolRPC() {
	rpcHost := config.RPCConfig.PoolHost
	rpcUser := config.RPCConfig.PoolUser
	rpcPassword := config.RPCConfig.PoolPassword

	poolRpc = poolrpc.NewPoolRPC(rpcHost, rpcUser, rpcPassword, poolrpc.WithPoolDebug(false))
	fmt.Printf("init pool rpc %s %s\n", rpcHost, rpcUser)
	// 添加定时任务 cron, 矿场机器 Farms 同步
	task.AddTimer(time.Minute*6, SyncPoolFarms)
	logrus.Debug("InitSpider AddTimer SyncPoolFarms, time duation 5 mins, repated")
	// 矿场订单状态同步
	task.AddTimer(time.Hour*12, SyncPoolMinerOrder)
	logrus.Debug("InitSpider AddTimer SyncPoolMinerOrder, time duation 12 hours, repated")
}
func SyncPoolFarms(param map[string]interface{}) bool {
	res, err := GetFarmsInfo()
	if err != nil {
		logrus.Error(err)
		return false
	}
	if len(res) == 0 {
		logrus.Debugf("SyncPoolFarms no machine")
		return false
	}
	for _, value := range res {
		UpdateFarmsInfoToDB(value)
	}
	return false
}
func SyncPoolMinerOrder(param map[string]interface{}) bool {
	// 从数据库获取所有订单
	var minerOrder []model.MinerOrder
	if err := db.DB.Model(&model.MinerOrder{}).Where("status < ?", model.MinerOrderCompleted).Find(&minerOrder).Error; err != nil {
		logrus.Error(err)
		return false
	}
	// 更新订单状态
	for _, value := range minerOrder {
		// 只更新未完成的订单
		if value.Status != model.MinerOrderCompleted {
			if status := UpdateMinerOrderInfoToDB(value.MinerOrderID); status != true {
				logrus.Debug("UpdateMinerOrderInfoToDB err")
			}
		}
	}
	return false
}

// 更新矿场机器信息
func UpdateFarmsInfoToDB(res model.FarmServer) error {
	var farms model.FarmServer
	if err := db.DB.Where("farm_id = ? AND miner_type = ?", res.FarmID, res.MinerType).First(&farms).RecordNotFound(); err == true {
		if err := db.DB.Create(&res).Error; err != nil {
			logrus.Errorf("updateFarmsInfoToDB service.DB.Create err %s", err.Error())
			return err
		}
		return nil
	}
	if err := db.DB.Model(farms).Where("farm_id = ? AND miner_type = ?", res.FarmID, res.MinerType).
		Updates(map[string]interface{}{"available_count": res.AvailableCount, "create_at": res.CreateAt}).Error; err != nil {
		logrus.Errorf("updateFarmsInfoToDB  update failed err: %s", err.Error())
		return err
	}
	// 更新类型价格列表
	for _, value := range res.PriceList {
		UpdateFarmsMinerPrice(farms.ID, value)
	}
	logrus.Infof("updateFarmsInfoToDB farm_id: %s miner_type: %s AvailableCount: %d Time: %s", res.FarmID, res.MinerType, res.AvailableCount, utils.TimeStFormat(res.CreateAt))
	return nil
}

// 更新矿场机器价格
func UpdateFarmsMinerPrice(id uint, res model.MinerPriceList) error {
	var farms model.MinerPriceList
	if err := db.DB.Where("farm_server_id = ? AND value = ?", id, res.Value).First(&farms).RecordNotFound(); err == true {
		if err := db.DB.Create(&res).Error; err != nil {
			logrus.Errorf("updateFarmsMinerPrice service.DB.Create err %s", err.Error())
			return err
		}
		return nil
	}
	if err := db.DB.Model(farms).Where("farm_server_id = ? AND value = ?", id, res.Value).
		Updates(map[string]interface{}{"time": res.Time}).Error; err != nil {
		logrus.Errorf("updateFarmsMinerPrice  update failed err: %s", err.Error())
		return err
	}
	logrus.Debugf("updateFarmsInfoToDB farm_id: %d value: %d time: %d", res.FarmServerID, res.Value, res.Time)
	return nil
}

// 创建订单
func CreateOrder(minerOrders *model.MinerOrderInfo) (string, error) {
	price := poolrpc.MachinePrice{
		Price: minerOrders.Price,
		Time:  minerOrders.Time,
	}
	config := poolrpc.MinerConfig{
		MinerPool:     minerOrders.MinerPool,
		MinerUserName: minerOrders.MinerUsername,
		MinerWorker:   minerOrders.MinerWorker,
	}
	orderInfo := &poolrpc.OrderInfo{
		FarmID: minerOrders.FarmID,
		Type:   minerOrders.GoodsType,
		Price:  price,
		Config: config,
	}

	if poolRpc != nil {
		return poolRpc.CreateOrder(orderInfo)
	}
	return "", errors.New("poolrpc GetFarmsInfo error")
}

// 获取 订单详情
func GetOrderInfo(orderId string) (*model.MinerOrder, error) {
	if poolRpc == nil {
		return nil, errors.New("poolrpc GetOrderInfo error")
	}
	orderRes, err := poolRpc.GetOrderInfo(orderId)
	if err != nil {
		return nil, err
	}
	// orderInfo := &model.MinerOrder{
	// 	MinerOrderID:  orderRes.ID,
	// 	FarmID:        orderRes.FarmID,
	// 	MinerID:       orderRes.MinerID,
	// 	CreateAt:      orderRes.CreatedAt.Unix(),
	// 	UpdateAt:      orderRes.UpdatedAt.Unix(),
	// 	Status:        orderRes.Status,
	// 	GoodsType:     orderRes.Request.Type,
	// 	RentTime:      orderRes.Request.Price.Time,
	// 	GoodsPrice:    orderRes.Request.Price.Price,
	// 	MinerPool:     orderRes.Request.Config.MinerPool,
	// 	MinerUsername: orderRes.Request.Config.MinerUserName,
	// 	MinerWorker:   orderRes.Request.Config.MinerWorker,
	// }
	return orderRes, nil
}

// 获取 矿场所有机器详情
func GetFarmsInfo() ([]model.FarmServer, error) {
	if poolRpc == nil {
		return nil, errors.New("poolrpc GetFarmsInfo error")
	}
	res, err := poolRpc.GetFarmsInfo()
	if err != nil {
		return nil, err
	}
	// var farmServers []model.FarmServer
	// for farmID, value := range res {
	// 	miners := value.Miners
	// 	for _, mValue := range miners {
	// 		priceList := GetMinerPrice(mValue.Type, mValue.PriceList)
	// 		farmServers = append(farmServers, model.FarmServer{
	// 			FarmID:         farmID,
	// 			MinerType:      mValue.Type,
	// 			PriceList:      priceList,
	// 			AvailableCount: mValue.AvailableCount,
	// 		})
	// 	}
	// }
	return res, nil
}

// 获取矿机价格
// func GetMinerPrice(minerType string, priceList []poolrpc.MachinePrice) []model.MinerPriceList {
// 	var minerPrice []model.MinerPriceList
// 	for _, value := range priceList {
// 		minerPrice = append(minerPrice, model.MinerPriceList{
// 			Value: value.Price,
// 			Time:  value.Time,
// 		})
// 	}
// 	return minerPrice
// }

// 查询指定 id 机器详情
func GetFarmInfo(farmID string) ([]model.FarmServer, error) {
	if poolRpc == nil {
		return nil, errors.New("poolrpc GetFarmsInfo error")
	}
	// var farmServers []model.FarmServer
	res, err := poolRpc.GetFarmInfo(farmID)
	if err != nil {
		return nil, err
	}
	// for _, mValue := range res.Miners {
	// 	priceList := GetMinerPrice(mValue.Type, mValue.PriceList)
	// 	farmServers = append(farmServers, model.FarmServer{
	// 		FarmID:         farmID,
	// 		MinerType:      mValue.Type,
	// 		PriceList:      priceList,
	// 		AvailableCount: mValue.AvailableCount,
	// 	})
	// }
	return res, nil
}

// 创建订单 定时任务订单回调函数
func CreateServiceToFarm(data map[string]interface{}) bool {

	outTradeNos, ok := data["out_trade_no"]
	if ok != true {
		return true
	}
	outTradeNo := outTradeNos.(string)

	farmID, ok := data["farm_id"]
	if ok != true {
		return true
	}
	goodsType, ok := data["miner_goods_type"]
	if ok != true {
		return true
	}
	price, ok := data["price"]
	if ok != true {
		return true
	}
	times, ok := data["time"]
	if ok != true {
		return true
	}
	minerPool, ok := data["miner_pool"]
	if ok != true {
		return true
	}
	minerUsername, ok := data["miner_username"]
	if ok != true {
		return true
	}
	minerWorker, ok := data["miner_worker"]
	if ok != true {
		return true
	}
	// 获取 farm 信息
	farmid := farmID.(string)
	res, err := GetFarmInfo(farmid)
	if err != nil {
		return false // 等待下次创建任务
	}
	if len(res) == 0 {
		return false // 等待下次创建任务
	}

	minerOrders := &model.MinerOrderInfo{
		FarmID:        farmid,
		GoodsType:     goodsType.(string),
		Price:         price.(uint64),
		Time:          times.(uint64),
		MinerPool:     minerPool.(string),
		MinerUsername: minerUsername.(string),
		MinerWorker:   minerWorker.(string),
	}
	orderId, err := CreateOrder(minerOrders)
	if err != nil {
		return false
	}
	orderInfo := &model.MinerOrder{
		MinerOrderID:  orderId,
		FarmID:        farmid,
		CreateAt:      utils.GetNowTime(),
		UpdateAt:      utils.GetNowTime(),
		Status:        model.MinerOrderCreating,
		GoodsType:     goodsType.(string),
		RentTime:      times.(uint64),
		GoodsPrice:    price.(uint64),
		MinerPool:     minerPool.(string),
		MinerUsername: minerUsername.(string),
		MinerWorker:   minerWorker.(string),
		OutTradeNo:    outTradeNo,
	}
	// 添加订单进入数据库
	if err := db.DB.Create(orderInfo).Error; err != nil {
		logrus.Errorf("CreateServiceToFarm %+v %s", orderInfo, err.Error())
		return false
	}
	// 更新数据库订单id，以及订单状态
	chargeParams := map[string]interface{}{
		"order_id": orderId,
	}
	timeDuation := 12 * time.Second
	// 添加定时任务 查询订单状态任务，并且第一次查询是在 15秒后，并更新数据库
	task.AddTimerWithDeadLine(timeDuation, time.Second*2, chargeParams, UpdateMinerOrderInfo)
	//只执行一次，因此返回 true
	return true
}

func UpdateMinerOrderInfo(data map[string]interface{}) bool {
	orderId, ok := data["order_id"]
	if ok != true {
		return true
	}
	UpdateMinerOrderInfoToDB(orderId.(string))
	return false
}
func UpdateMinerOrderInfoToDB(orderId string) bool {
	res, err := GetOrderInfo(orderId)
	if err != nil {
		return false
	}
	if err := db.DB.Model(model.MinerOrder{}).Update(res).Error; err != nil {
		logrus.Errorf("UpdateMinerOrderInfo %s %s", orderId, err.Error())
	}
	return true
}
