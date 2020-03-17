package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/user"
	// "gopkg.in/kataras/iris.v6"
)

// 用户登录
func Signin(ctx iris.Context) {
	aul := new(model.UserLogin)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数有误"))
		return
	}
	response, msg, status := user.Signin(ctx, aul, false)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// 用户注册
func Signup(ctx iris.Context) {
	aul := new(model.UserJson)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	msg, status := user.Signup(ctx, aul)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// 重置密码
func ResetPassword(ctx iris.Context) {
	passwdData := new(model.UserUpdatePwd)
	if err := ctx.ReadJSON(passwdData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.ErrorCode.ERROR, nil, "参数无效"))
		return
	}
	if err := validate.Struct(passwdData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.ErrorCode.ERROR, nil, "参数无效"))
		return
	}
	msg, status := user.ResetPassword(ctx, passwdData)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// 注销
func Signout(ctx iris.Context) {
	msg, status := user.Signout(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))

}

// GetUserAsset 获取用户资产
func GetUserAsset(ctx iris.Context) {
	response, msg, status := user.GetUserAllAmount(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAssetFlow 获取资产流水
func GetAssetFlow(ctx iris.Context) {
	assetFlow := new(model.AssetflowJson)
	if err := ctx.ReadQuery(assetFlow); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(assetFlow); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := user.GetUserAssetFlow(ctx, assetFlow)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// 获取用户个人信息
func GetUserProfile(ctx iris.Context) {
	response, msg, status := user.GetUserProfileInfo(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateInfo 更新用户信息
func UpdateUserName(ctx iris.Context) {
	users := new(model.UserUpdateJson)
	if err := ctx.ReadJSON(users); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	msg, status := user.UpdateUserName(ctx, users)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdatePassword 更新用户密码
func UpdatePassword(ctx iris.Context) {

	var userData model.PasswordUpdateData
	if err := ctx.ReadJSON(&userData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.ErrorCode.ERROR, nil, "参数无效"))
		return
	}
	if err := validate.Struct(&userData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.ErrorCode.ERROR, nil, "参数无效"))
		return
	}
	msg, status := user.UpdatePassword(ctx, &userData)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetMessageList 获取用户消息列表
func GetMessageList(ctx iris.Context) {
	userMessage := new(model.UserMessageJson)
	if err := ctx.ReadQuery(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	response, msg, status := user.GetUserMessages(ctx, userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetMessageDetail 获取消息详情信息
func GetMessageDetail(ctx iris.Context) {

	userMessage := new(model.MessageDetailJson)
	if err := ctx.ReadQuery(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	response, msg, status := user.GetUserMessageDetail(ctx, userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))

}

// UpdateMessage 更新用户消息状态
func UpdateMessage(ctx iris.Context) {
	userMessage := new(model.MessageDetailJson)
	if err := ctx.ReadQuery(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	msg, status := user.UpdateUserMessageStatus(ctx, userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// 获取用户登录日志
func GetLoginList(ctx iris.Context) {

	userOauthJson := new(model.UserOauthJson)
	if err := ctx.ReadQuery(userOauthJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	if err := validate.Struct(userOauthJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	response, msg, status := user.GetUserLoginList(ctx, userOauthJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
