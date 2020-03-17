package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/hook"
)

func GithubHook(ctx iris.Context) {
	bodyContent, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	msg, status := hook.GithubHook(ctx, string(bodyContent))
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

func GithubHookLog(ctx iris.Context) {
	queryData := new(model.QueryHookJson)
	if err := ctx.ReadQuery(queryData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_FAILED, nil, "参数无效"))
		return
	}
	response, msg, status := hook.GithubHookLog(queryData)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
