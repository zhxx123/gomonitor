package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/common"
	"github.com/zhxx123/gomonitor/service/user"
	"github.com/zhxx123/gomonitor/service/wallet"
)

// GetCaptchaCode 获取图片验证码
func GetCaptchaCode(ctx iris.Context) {
	verifyCode := new(model.CaptchaCodeJson)
	if err := ctx.ReadJSON(verifyCode); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(verifyCode); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	response, msg, status := common.CreateCaptchaCode(verifyCode)
	ctx.JSON(ApiResource(status, response, msg))
}

// CheckCaptchaCode 校验验证码，主要用于测试
func CheckCaptchaCode(ctx iris.Context) {
	verifyCode := new(model.CodeJson)
	if err := ctx.ReadJSON(verifyCode); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(verifyCode); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	msg, status := common.CheckCaptchaCodes(verifyCode)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetEmailPhoneCode 获取邮箱或者手机验证码
func GetEmailPhoneCode(ctx iris.Context) {
	verifyCode := new(model.VerifyCodeJson)
	if err := ctx.ReadJSON(verifyCode); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(verifyCode); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	ctx.StatusCode(iris.StatusOK)
	msg, status := user.CreateEmailPhoneCode(verifyCode)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UploadImage 上传图片
func UploadImage(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	response, msg, status := common.UploadImage(ctx)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetWorkOrderID 获取workorder ID
func GetWorkOrderID(ctx iris.Context) {
	workOrderJson := new(model.WorkOrderIDJson)
	if err := ctx.ReadJSON(workOrderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(workOrderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	response, msg, status := common.GetWorkID(ctx, workOrderJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))

}

// WorkOrder 创建工单信息
func WorkOrder(ctx iris.Context) {
	workOrderJson := new(model.WorkOrderJson)
	if err := ctx.ReadJSON(workOrderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	if err := validate.Struct(workOrderJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数错误"))
		return
	}
	msg, status := common.UploadWorkOrder(ctx, workOrderJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetCoinPrice 获取币种价格
func GetCoinPrice(ctx iris.Context) {
	coinName := ctx.Params().GetString("id")
	if coinName == "" {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, "参数无效"))
		return
	}
	response, msg, status := wallet.GetCoinPriceWithCoinName(coinName)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetCoinPriceList 表示获取所有币种价格
func GetCoinPriceList(ctx iris.Context) {
	response, msg, status := wallet.GetCoinPrices()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
