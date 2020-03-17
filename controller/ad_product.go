package controller

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/admin"
)

// GetAdminGoodsList 获取商品列表
func GetAdminGoodsList(ctx iris.Context) {
	productJson := new(model.AProductsJson)
	if err := ctx.ReadQuery(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminProductList(productJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// CreateAdminGoods 新增商品
func CreateAdminGoods(ctx iris.Context) {
	productJson := new(model.AUpdateProductsJson)
	if err := ctx.ReadJSON(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.CreateAdminProducts(ctx, productJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminGood 更新商品信息
func UpdateAdminGood(ctx iris.Context) {
	productJson := new(model.AUpdateProductsJson)
	if err := ctx.ReadJSON(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminProducts(ctx, productJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminGoodStatus 更新商品状态
func UpdateAdminGoodStatus(ctx iris.Context) {
	productJson := new(model.AUpdateProductStatusJson)
	if err := ctx.ReadJSON(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminProductStatus(ctx, productJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// QueryAdminGood 查询单个商品详情信息
func QueryAdminGood(ctx iris.Context) {
	goods_id := ctx.Values().GetString("id")
	if goods_id == "" {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_FAILED, nil, "参数错误"))
	} else {
		response, msg, status := admin.GetAdminProductDetail(goods_id)
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(ApiResource(status, response, msg))
	}
}

// GetAdminFarmsList 获取 farm 列表
func GetAdminFarmsList(ctx iris.Context) {
	farmJson := new(model.AFarmServerJson)
	if err := ctx.ReadQuery(farmJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(farmJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminFarmsList(farmJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
