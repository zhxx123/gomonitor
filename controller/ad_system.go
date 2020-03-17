package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/admin"
)

// GetAdminSysBasicListInfo 获取系统基本信息列表
func GetAdminSysBasicListInfo(ctx iris.Context) {
	sysInfoJson := new(model.SysInfoJson)
	if err := ctx.ReadQuery(sysInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(sysInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.CheckAndAdminBasicListData(sysInfoJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminSysSimpleInfo 获取系统状态信息
func GetAdminSysSimpleInfo(ctx iris.Context) {
	sysInfoJson := new(model.SysInfoJson)
	if err := ctx.ReadQuery(sysInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(sysInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.CheckAndAdminSimpleData(sysInfoJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminWalletBasicInfo 获取钱包基本信息
func GetAdminWalletBasicInfo(ctx iris.Context) {
	walletInfoJson := new(model.AWalletInfoJson)
	if err := ctx.ReadQuery(walletInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminWalletBasic(walletInfoJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminWalletSimpleInfo 获取钱包状态信息
func GetAdminWalletSimpleInfo(ctx iris.Context) {
	walletInfoJson := new(model.AWalletInfoJson)
	if err := ctx.ReadQuery(walletInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletInfoJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
	}
	response, msg, status := admin.GetAdminWalletSimple(walletInfoJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminWalletAddressRecordList 获取钱包地址分配表
func GetAdminWalletAddressRecordList(ctx iris.Context) {
	walletAddressJson := new(model.AWalletAddressJson)
	if err := ctx.ReadQuery(walletAddressJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletAddressJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminWalletAddressRecord(walletAddressJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminWalletAddress 上传钱包地址
func UpdateAdminWalletAddress(ctx iris.Context) {
	walletAddressJson := new(model.AUpdateWalletAddressJson)
	if err := ctx.ReadJSON(walletAddressJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletAddressJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminWalletAddress(walletAddressJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminUserCoinAddress 更新用户钱包地址分配
func UpdateAdminUserCoinAddress(ctx iris.Context) {
	walletAddressJson := new(model.AUserCoinAddressJson)
	if err := ctx.ReadJSON(walletAddressJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletAddressJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminUserNewCoinAddress(walletAddressJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminCoinPrice 获取虚拟货币价格
func GetAdminCoinPrice(ctx iris.Context) {
	response, msg, status := admin.GetAdminCoinPriceList()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminCoinMarketList 获取货币价格市场
func GetAdminCoinMarketList(ctx iris.Context) {
	coinMarketJson := new(model.ACoinMarketJson)
	if err := ctx.ReadQuery(coinMarketJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(coinMarketJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminCoinMarket(coinMarketJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminCoinPrice 手动更新货币价格
func UpdateAdminCoinPrice(ctx iris.Context) {
	coinPriceJson := new(model.AUpdateCoinPriceJson)
	if err := ctx.ReadJSON(coinPriceJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(coinPriceJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminCoinPrices(ctx, coinPriceJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminCoinPriceAutoUp 更新价格状态
func UpdateAdminCoinPriceAutoUp(ctx iris.Context) {
	coinMarketJson := new(model.AUpdateStatusCoinPriceJson)
	if err := ctx.ReadJSON(coinMarketJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(coinMarketJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminCoinPriceStatus(coinMarketJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminSettingsList 获取系统设置
func GetAdminSettingsList(ctx iris.Context) {
	settingJson := new(model.ASettingsJson)
	if err := ctx.ReadQuery(settingJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(settingJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminUserSettingsList(settingJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminSettings 更新系统设置
func UpdateAdminSettings(ctx iris.Context) {
	settingJson := new(model.ASettingsUpdateJson)
	if err := ctx.ReadJSON(settingJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(settingJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminSettings(ctx, settingJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminLogList 获取其他系统日志记录
func GetAdminLogList(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(model.STATUS_SUCCESS, nil, "success"))
}

// GetAdminUserLoginLogs 获取当前用户登录记录
func GetAdminUserLoginLogs(ctx iris.Context) {
	userOauthJson := new(model.UserOauthJson)
	if err := ctx.ReadQuery(userOauthJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userOauthJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetUserAdminLoginList(ctx, userOauthJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// 获取 github 提交记录
func WebHookLog(ctx iris.Context) {
	queryData := new(model.QueryHookJson)
	if err := ctx.ReadQuery(queryData); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_FAILED, nil, "参数无效"))
		return
	}
	response, msg, status := admin.WebHookLog(queryData)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}
