package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/service/user"
	// "gopkg.in/kataras/iris.v6"
)

func Login(ctx iris.Context) {

	msg, status := user.Login()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
func Logout(ctx iris.Context) {

	msg, status := user.Logout()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

func ArtilceList(ctx iris.Context) {

	msg, status := user.ArtilceList()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

func ArticleInfo(ctx iris.Context) {

	msg, status := user.ArticleInfo()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
