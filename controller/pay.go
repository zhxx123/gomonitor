package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/pay"
)

// GetPayID 获取订单ID
func GetPayID(ctx iris.Context) {
	payIDJson := new(model.PayIDJson)
	if err := ctx.ReadJSON(payIDJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(payIDJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := pay.CheckAndGetPayID(ctx, payIDJson)

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// PayCreate 创建订单
func PayCreate(ctx iris.Context) {
	payCreateJson := new(model.PayCreateJson)

	if err := ctx.ReadJSON(payCreateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(payCreateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
	}
	response, msg, status := pay.CheckAndPayCreate(ctx, payCreateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// PayQuery 查询订单
func PayQuery(ctx iris.Context) {
	payQueryJson := new(model.PayQueryJson)

	if err := ctx.ReadQuery(payQueryJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(payQueryJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := pay.CheckAndPayQuery(ctx, payQueryJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// PayCancel 取消第三方支付订单
func PayCancel(ctx iris.Context) {

	payCancelJson := new(model.PayCancelJson)
	if err := ctx.ReadJSON(payCancelJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(payCancelJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := pay.CheckAndPayCancel(ctx, payCancelJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// PayUpdate 删除订单
func PayUpdate(ctx iris.Context) {

	payCloseJson := new(model.PayQueryJson)

	if err := ctx.ReadJSON(payCloseJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(payCloseJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := pay.CheckAndDelatePayOrder(ctx, payCloseJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
