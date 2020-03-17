package controller

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/admin"
)

// GetAdminOrderList 获取订单列表
func GetAdminOrderList(ctx iris.Context) {

	orderJson := new(model.AOrderListJson)
	if err := ctx.ReadQuery(orderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(orderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminOrderLists(orderJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminOrderStatus 更新订单状态
func UpdateAdminOrderStatus(ctx iris.Context) {
	orderJson := new(model.AOrderUpdateStatusJson)
	if err := ctx.ReadJSON(orderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(orderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminOrderStatus(orderJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// QueryAdminOrder
func QueryAdminOrder(ctx iris.Context) {
	response, msg, status := admin.GetAdminOrderDetail(ctx)
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminMinerList 获取所有 矿场订单 orders
func GetAdminMinerList(ctx iris.Context) {
	orderListJson := new(model.AOrderListJson)
	if err := ctx.ReadQuery(orderListJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(orderListJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := admin.GetAdminMinerList(orderListJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
