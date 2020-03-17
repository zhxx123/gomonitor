package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/admin"
)

// 获取系统资产
func GetAdminSummay(ctx iris.Context) {
	response, msg, status := admin.GetAdminSystemAccountList()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// 获取资产报表列表
func GetAdminFinanceReportList(ctx iris.Context) {
	reportJson := new(model.AReportJson)
	if err := ctx.ReadQuery(reportJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(reportJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminFinanceReports(reportJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// AddAdminFinanceReport 添加一条资产报表
func AddAdminFinanceReport(ctx iris.Context) {
	reportJson := new(model.AUpdateReportsJson)
	if err := ctx.ReadJSON(reportJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(reportJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.CreateAdminFinanceReports(ctx, reportJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminReportStatus 修改报表状态
func UpdateAdminReportStatus(ctx iris.Context) {
	reportJson := new(model.AUpdateReportStatusJson)
	if err := ctx.ReadJSON(reportJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(reportJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminFinanaceReport(ctx, reportJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminDigtialList 获取数字交易货币列表
func GetAdminDigtialList(ctx iris.Context) {
	walletJson := new(model.AWalletTxJson)
	if err := ctx.ReadQuery(walletJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminDigtials(walletJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminDigtial 更新数字货币交易记录
func UpdateAdminDigtial(ctx iris.Context) {
	walletJson := new(model.AWalletRecordJson)
	if err := ctx.ReadJSON(walletJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(walletJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminWalletRecord(walletJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminOfficeList 获取法定货币列表
func GetAdminOfficeList(ctx iris.Context) {
	payTxJson := new(model.APayTxJson)
	if err := ctx.ReadQuery(payTxJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(payTxJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminOffices(payTxJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminOffice 更新法币列表
func UpdateAdminOffice(ctx iris.Context) {
	payTxJson := new(model.AUpatePayTxJson)
	if err := ctx.ReadJSON(payTxJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(payTxJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminOfficeRecordStatus(payTxJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminUserAssetList 获取用户资产列表
func GetAdminUserAssetList(ctx iris.Context) {
	userAssetJson := new(model.AUserAccountsJson)
	if err := ctx.ReadQuery(userAssetJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userAssetJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminUserAccount(userAssetJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminAssetFlow 获取用户资产流水
func GetAdminAssetFlow(ctx iris.Context) {
	userAssetJson := new(model.AUserAccountsJson)
	if err := ctx.ReadQuery(userAssetJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(userAssetJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminUserAssetFlow(userAssetJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// GetAdminVirtualList 获取虚拟充值列表
func GetAdminVirtualList(ctx iris.Context) {
	virTxJson := new(model.AVirtualAccountJson)
	if err := ctx.ReadQuery(virTxJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(virTxJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminVirtuals(virTxJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// AddAdminAssetVirtual 添加用户虚拟充值
func AddAdminAssetVirtual(ctx iris.Context) {
	virAssetJson := new(model.AVirtualAssetJson)
	if err := ctx.ReadJSON(virAssetJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(virAssetJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.AddAdminUserVirtualAccount(ctx, virAssetJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
