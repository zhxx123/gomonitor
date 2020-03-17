package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/stats"
)

// 添加访问统计
func AddPV(ctx iris.Context) {
	pvData := new(model.ClientInfo)
	if err := ctx.ReadQuery(pvData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.ErrorCode.ERROR, nil, "参数无效"))
		return
	}
	msg, status := stats.PV(ctx, pvData)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
