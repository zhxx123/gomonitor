package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/service/db"

	"github.com/zhxx123/gomonitor/model"
	"github.com/sirupsen/logrus"
)

// GetAdminOrderLists 获取订单列表
func GetAdminOrderLists(orderJson *model.AOrderListJson) (model.MyMap, string, int) {
	var orderList []model.Orders
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(orderJson.Page, orderJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if len(orderJson.OutTradeNo) > 0 {
		dbs = dbs.Where("out_trade_no = ?", orderJson.OutTradeNo)
	}
	if orderJson.OrderType != 0 {
		dbs = dbs.Where("order_type = ?", orderJson.OrderType)
	}
	if orderJson.OrderStatus != 0 {
		dbs = dbs.Where("order_status = ?", orderJson.OrderStatus)
	}
	orders := new(model.Orders)
	if err := dbs.Model(orders).Count(&count).
		Offset(offset).Limit(orderJson.Limit).
		Find(&orderList).Error; err != nil {
		logrus.Errorf("GetAdminOrderLists failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminOrderLists count: %d page: %d limit: %d offset: %d\n", count, orderJson.Page, orderJson.Limit, offset)

	response := model.MyMap{
		"data":   orderList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminMinerList 获取矿场订单列表
func GetAdminMinerList(orderJson *model.AOrderListJson) (model.MyMap, string, int) {
	var orderList []model.MinerOrder
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(orderJson.Page, orderJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if len(orderJson.OutTradeNo) > 0 {
		dbs = dbs.Where("out_trade_no = ?", orderJson.OutTradeNo)
	}
	if orderJson.OrderStatus != 0 {
		dbs = dbs.Where("order_status = ?", orderJson.OrderStatus)
	}
	orders := new(model.MinerOrder)
	if err := dbs.Model(orders).Count(&count).
		Offset(offset).Limit(orderJson.Limit).
		Find(&orderList).Error; err != nil {
		logrus.Errorf("GetAdminMinerList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminMinerList count: %d page: %d limit: %d offset: %d\n", count, orderJson.Page, orderJson.Limit, offset)

	response := model.MyMap{
		"data":   orderList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UpdateAdminOrderStatus 更新订单状态
func UpdateAdminOrderStatus(orderJson *model.AOrderUpdateStatusJson) (string, int) {
	product := new(model.Orders)
	if err := db.DB.Model(product).Where("out_trade_no = ?", orderJson.OutTradeNo).
		Update("order_status", orderJson.OrderStatus).Error; err != nil {
		logrus.Errorf("UpdateAdminOrderStatus failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// GetAdminOrderDetail 查询订单详情
func GetAdminOrderDetail(ctx iris.Context) (*model.Orders, string, int) {
	outTradeNo := ctx.Values().GetString("id")
	if outTradeNo == "" {
		return nil, "参数错误", model.STATUS_FAILED
	}
	var order model.Orders
	if err := db.DB.Where("out_trade_no = ?", outTradeNo).First(&order).Error; err != nil {
		logrus.Errorf("GetAdminOrderDetail failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	return &order, "success", model.STATUS_SUCCESS
}
