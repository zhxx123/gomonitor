package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/set"
)

// SetKeyValue 设置key, value
func SetKeyValue(ctx iris.Context) {
	kvData := new(model.KeyValueData)
	if err := ctx.ReadJSON(kvData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.ErrorCode.ERROR, nil, "参数无效"))
		return
	}
	response, msg, status := set.SetKeyValue(ctx, kvData)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
