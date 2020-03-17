package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/service/seo"
)

// PushToBaidu 百度链接提交
func PushToBaidu(ctx iris.Context) {
	msg, status := seo.PushToBaidu()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
