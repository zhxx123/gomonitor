package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/product"
)

func GetAllProduct(ctx iris.Context) {
	productJson := new(model.ProductsQueryJson)
	if err := ctx.ReadQuery(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "bad params"))
		return
	}
	if err := validate.Struct(productJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "bad request"))
		return
	}
	response, msg, status := product.GetProductsList(productJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

func GetProDetail(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(model.STATUS_SUCCESS, nil, "success"))
}
