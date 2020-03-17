package controller

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/admin"
	"github.com/zhxx123/gomonitor/service/user"
)

// 管理员账户登录
func AdminSignin(ctx iris.Context) {
	aul := new(model.UserLogin)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数有误"))
		return
	}
	response, msg, status := user.Signin(ctx, aul, true)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminUserList 获取用户列表
func GetAdminUserList(ctx iris.Context) {
	aul := new(model.AUserJson)

	if err := ctx.ReadQuery(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.CheckAndGetUserList(aul)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminUserInfo 获取用户个人信息
func GetAdminUserInfo(ctx iris.Context) {
	response, msg, status := admin.GetAdminUserInfo(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminSpecificUserInfo 获取指定用户详细信息
func GetAdminSpecificUserInfo(ctx iris.Context) {
	user, msg, status := admin.GetAdminUserInfomation(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, user, msg))
}

// 更新用户信息
func UpdateAdminUserStatusRole(ctx iris.Context) {
	aul := new(model.AUserUpdateStatusJson)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminUserStatusRole(aul)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// DeleteAdminUser 删除用户
func DeleteAdminUser(ctx iris.Context) {
	msg, status := admin.CheckAndDeleteUser(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// LogoutAdminUser 下线指定用户
func LogoutAdminUser(ctx iris.Context) {
	msg, status := admin.CheckAndLogout(ctx)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminLoginLogs 获取用户登录日志
func GetAdminLoginLogs(ctx iris.Context) {
	logLogJson := new(model.ALoginLogsJson)
	if err := ctx.ReadQuery(logLogJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(logLogJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminUserLoginLogsList(logLogJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UserAdminLogout 管理员账号退出登录
func UserAdminLogout(ctx iris.Context) {
	msg, status := admin.AdminUserLogout(ctx)
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminMessageList 获取用户消息列表
func GetAdminMessageList(ctx iris.Context) {
	userMessage := new(model.AUserMessageJson)
	if err := ctx.ReadQuery(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminUserMessages(userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// AddAdminUserMessage 添加用户消息
func AddAdminUserMessage(ctx iris.Context) {
	userMessage := new(model.AEditUserMessageJson)
	if err := ctx.ReadJSON(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.AddAdminUserMessages(ctx, userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminMessageStatus 更新用户消息状态
func UpdateAdminMessageStatus(ctx iris.Context) {
	userMessage := new(model.AUpdateUserMessageJson)
	if err := ctx.ReadJSON(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "bad request"))
		return
	}
	msg, status := admin.UpdateAdminUserMessageStatus(userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminMessage 更新用户消息
func UpdateAdminMessage(ctx iris.Context) {

	userMessage := new(model.AEditUserMessageJson)
	if err := ctx.ReadJSON(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userMessage); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminUserMessage(ctx, userMessage)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminUserMangerList 获取管理员列表
func GetAdminUserMangerList(ctx iris.Context) {
	aul := new(model.AUserJson)
	if err := ctx.ReadQuery(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminUserMangerList(aul)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))

}
