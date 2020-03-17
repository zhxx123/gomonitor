package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/order"
)

// GetOrderID 获取订单 ID
func GetOrderID(ctx iris.Context) {

	preTradeJson := new(model.PreTradeJson)

	if err := ctx.ReadJSON(preTradeJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(preTradeJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := order.CheckAndGetOrderID(ctx, preTradeJson)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// OrderCreate 创建商品订单
func OrderCreate(ctx iris.Context) {

	preCreateJson := new(model.PreCreateJson)

	if err := ctx.ReadJSON(preCreateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(preCreateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := order.CheckAndPreCreate(ctx, preCreateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// OrderCreatePay 确认商品支付订单
func OrderCreatePay(ctx iris.Context) {
	preCreateJson := new(model.PreCreatePayJson)

	if err := ctx.ReadJSON(preCreateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(preCreateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := order.CheckAndPayOrder(ctx, preCreateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// OrderQuery 查询创建订单
func OrderQuery(ctx iris.Context) {
	orderQueryJson := new(model.OrderQueryJson)
	if err := ctx.ReadQuery(orderQueryJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(orderQueryJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := order.OrderCreateQuery(ctx, orderQueryJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// OrderCancel 取消未付款订单
func OrderCancel(ctx iris.Context) {
	orderCancelJson := new(model.OrderCancelJson)
	if err := ctx.ReadJSON(orderCancelJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(orderCancelJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := order.OrderCancel(ctx, orderCancelJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// OrderUpdate 更新订单,用户删除订单
func OrderUpdate(ctx iris.Context) {

	orderCloseJson := new(model.OrderCloseJson)

	if err := ctx.ReadJSON(orderCloseJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(orderCloseJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := order.OrderDelete(ctx, orderCloseJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetUserOrders 获取用户所有 orders
func GetUserOrderList(ctx iris.Context) {
	orderListJson := new(model.OrderListJson)
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
	response, msg, status := order.GetUserOrders(ctx, orderListJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
