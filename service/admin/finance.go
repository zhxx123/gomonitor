package admin

import (
	"errors"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/pay"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

// GetAdminSystemAccountList 获取系统报表
func GetAdminSystemAccountList() (model.MyMap, string, int) {
	var sysAccountList []model.SystemAccount
	count := 0
	if err := db.DB.Model(model.SystemAccount{}).Count(&count).
		Find(&sysAccountList).Error; err != nil {
		logrus.Errorf("GetAdminSysAccountListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(sysAccountList)
	response := model.MyMap{
		"data":   sysAccountList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminFinanceReports 获取报表列表
func GetAdminFinanceReports(reportJson *model.AReportJson) (model.MyMap, string, int) {
	var assetReportList []model.Assets
	count := 0
	dbs := db.DB
	if reportJson.IsStatus == true {
		dbs = dbs.Where("status = ?", reportJson.Status)
	}
	if len(reportJson.CoinType) > 0 {
		dbs = dbs.Where("coin_type = ?", reportJson.CoinType)
	}
	offset, err := db.GetOffset(reportJson.Page, reportJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := dbs.Model(model.Assets{}).Order("id desc").Count(&count).
		Offset(offset).Limit(reportJson.Limit).
		Find(&assetReportList).Error; err != nil {
		logrus.Errorf("GetAdminFinanaceReportFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	response := model.MyMap{
		"data":   assetReportList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UpdateAdminFinanaceReport 修改报表状态
func UpdateAdminFinanaceReport(ctx iris.Context, reportJson *model.AUpdateReportStatusJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	asset := new(model.Assets)
	if err := db.DB.Model(asset).Where("ID = ?", reportJson.ID).Updates(map[string]interface{}{"author_id": userId, "status": reportJson.Status}).Error; err != nil {
		logrus.Errorf("UpdateAdminFinanaceReportStatusToDB failed err: %s", err.Error())
		return "failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// CreateAdminFinanceReports 新增报表
func CreateAdminFinanceReports(ctx iris.Context, reportJson *model.AUpdateReportsJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	endTime, err := GetNextTime(reportJson.EndAt)

	if err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 指定类型的报表
	if len(reportJson.CoinType) > 0 {
		if err := CreateAdminReport(endTime, reportJson.CoinType, userId); err != nil {
			return err.Error(), model.STATUS_FAILED
		}
		return "success", model.STATUS_SUCCESS
	}
	// 未指定，依次生成
	// 1  添加 CNY报表
	coinType := "CNY"
	if err := CreateAdminReport(endTime, coinType, userId); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 2 添加 MGD 报表
	coinType = "MGD"
	if err := CreateAdminReport(endTime, coinType, userId); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 3 添加 ETH 报表
	coinType = "ETH"
	if err := CreateAdminReport(endTime, coinType, userId); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 4 添加 BTC 报表
	coinType = "BTC"
	if err := CreateAdminReport(endTime, coinType, userId); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// GetAdminDigtials 获取数字货币交易记录
func GetAdminDigtials(walletJson *model.AWalletTxJson) (model.MyMap, string, int) {
	count := 0
	dbs := db.DB
	if walletJson.Txid != "" {
		dbs = dbs.Where("tx_id = ?", walletJson.Txid)
	}
	offset, err := db.GetOffset(walletJson.Page, walletJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	switch walletJson.CoinType {
	case "ETH":
		var walletTxList []model.EthTx
		if err := dbs.Model(model.EthTx{}).Count(&count).
			Offset(offset).Limit(walletJson.Limit).
			Find(&walletTxList).Error; err != nil {
			logrus.Errorf("GetACoinPriceList failed err: %s", err.Error())
			return nil, "error", model.STATUS_FAILED
		}
		response := model.MyMap{
			"data":   walletTxList,
			"length": count,
		}
		return response, "success", model.STATUS_SUCCESS

	case "MGD":
		var walletTxList []model.MgdTx
		if err := dbs.Model(model.MgdTx{}).Count(&count).
			Offset(offset).Limit(walletJson.Limit).
			Find(&walletTxList).Error; err != nil {
			logrus.Errorf("GetACoinPriceList failed err: %s", err.Error())
			return nil, "error", model.STATUS_FAILED
		}
		response := model.MyMap{
			"data":   walletTxList,
			"length": count,
		}
		return response, "success", model.STATUS_SUCCESS

	case "BTC":
		var walletTxList []model.BtcTx
		if err := dbs.Model(model.BtcTx{}).Count(&count).
			Offset(offset).Limit(walletJson.Limit).
			Find(&walletTxList).Error; err != nil {
			logrus.Errorf("GetACoinPriceList failed err: %s", err.Error())
			return nil, "error", model.STATUS_FAILED
		}
		response := model.MyMap{
			"data":   walletTxList,
			"length": count,
		}
		return response, "success", model.STATUS_SUCCESS
	default:
		return nil, "error", model.STATUS_FAILED
	}
}

// UpdateAdminWalletRecord 更新数字货币记录表
func UpdateAdminWalletRecord(walletJson *model.AWalletRecordJson) (string, int) {
	message := new(model.WalletRecord)
	if err := db.DB.Model(message).Where("tx_id = ?", walletJson.Txid).
		Update("added", walletJson.Added).Error; err != nil {
		logrus.Errorf("UpdateAdminWalletRecordToDB falied Err:%s", err.Error())
		return "failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// GetAdminOffices 获取 数字法币交易记录
func GetAdminOffices(payTxJson *model.APayTxJson) (model.MyMap, string, int) {
	count := 0
	dbs := db.DB
	if payTxJson.OutTradeNo != "" {
		dbs = dbs.Where("out_trade_no = ?", payTxJson.OutTradeNo)
	}
	if payTxJson.TradeStatus != 0 {
		dbs = dbs.Where("trade_status = ?", payTxJson.TradeStatus)
	}
	offset, err := db.GetOffset(payTxJson.Page, payTxJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	var payTxList []model.PayTx

	payTx := new(model.PayTx)
	if err := dbs.Model(payTx).Count(&count).
		Offset(offset).Limit(payTxJson.Limit).
		Find(&payTxList).Error; err != nil {
		logrus.Errorf("GetACoinPriceList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminOfficeList out_trade_no: %s status: %d count: %d page: %d limit: %d offset: %d \n", payTxJson.OutTradeNo, payTxJson.TradeStatus, count, payTxJson.Page, payTxJson.Limit, offset)
	response := model.MyMap{
		"data":   payTxList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UpdateAdminOfficeRecordStatus 更新数字货币记录表
func UpdateAdminOfficeRecordStatus(payTxJson *model.AUpatePayTxJson) (string, int) {
	payTx := new(model.PayTx)
	if err := db.DB.Model(&payTx).Where("out_trade_no = ?", payTxJson.OutTradeNo).Update("status", payTxJson.TradeStatus).Error; err != nil {
		logrus.Errorf("UpdateAdminPayTxStatusToDB falied Err:%s", err.Error())
		return "failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 获取用户资产列表
func GetAdminUserAccount(userAssetJson *model.AUserAccountsJson) (model.MyMap, string, int) {
	var userAssetList []model.UserAccounts
	count := 0
	dbs := db.DB
	if userAssetJson.UserId != 0 {
		dbs = dbs.Where("user_id = ?", userAssetJson.UserId)
	}
	if userAssetJson.CoinType != "" {
		dbs = dbs.Where("coin_type = ?", userAssetJson.CoinType)
	}
	offset, err := db.GetOffset(userAssetJson.Page, userAssetJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := dbs.Model(model.UserAccounts{}).Count(&count).
		Offset(offset).Limit(userAssetJson.Limit).
		Find(&userAssetList).Error; err != nil {
		logrus.Errorf("GetAdminUserAssetListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminUserAssetListFromDB count: %d page: %d limit: %d offset: %d\n", count, userAssetJson.Page, userAssetJson.Limit, offset)
	response := model.MyMap{
		"data":   userAssetList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminUserAssetFlow 获取资产流水列表
func GetAdminUserAssetFlow(userAccountsJson *model.AUserAccountsJson) (model.MyMap, string, int) {
	var userAccountsList []model.UserAssetflow
	count := 0
	dbs := db.DB
	if userAccountsJson.UserId != 0 {
		dbs = dbs.Where("user_id = ?", userAccountsJson.UserId)
	}
	if userAccountsJson.CoinType != "" {
		dbs = dbs.Where("coin_type = ?", userAccountsJson.CoinType)
	}
	offset, err := db.GetOffset(userAccountsJson.Page, userAccountsJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	userAssetFlow := new(model.UserAssetflow)
	if err := dbs.Model(userAssetFlow).Count(&count).
		Offset(offset).Limit(userAccountsJson.Limit).
		Find(&userAccountsList).Error; err != nil {
		logrus.Errorf("GetAdminUserAssetFlowFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	response := model.MyMap{
		"data":   userAccountsList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminVirtuals 获取虚拟充值记录
func GetAdminVirtuals(virtualJson *model.AVirtualAccountJson) (model.MyMap, string, int) {
	var virtualAccountsList []model.VirtualRecharge
	count := 0
	dbs := db.DB
	if virtualJson.UserId != 0 {
		dbs = dbs.Where("user_id = ?", virtualJson.UserId)
	}
	if virtualJson.CoinType != "" {
		dbs = dbs.Where("coin_type = ?", virtualJson.CoinType)
	}
	offset, err := db.GetOffset(virtualJson.Page, virtualJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	virutalRecharge := new(model.VirtualRecharge)
	if err := dbs.Model(virutalRecharge).Count(&count).
		Offset(offset).Limit(virtualJson.Limit).
		Find(&virtualAccountsList).Error; err != nil {
		logrus.Errorf("GetAVirtualsList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAVirtualsList count: %d page: %d limit: %d offset: %d\n", count, virtualJson.Page, virtualJson.Limit, offset)

	response := model.MyMap{
		"data":   virtualAccountsList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// AddAdminUserVirtualAccount 添加虚拟资产到用户账户
func AddAdminUserVirtualAccount(ctx iris.Context, virUserAssetJson *model.AVirtualAssetJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	// 获取订单号
	outTradeNo := utils.GenOutTradeNo(model.COIN_CHANNEL_VIRTUAL, 0, 1, virUserAssetJson.UserId)
	createAt := utils.GetNowTime()
	// 获取用户资产充值记录
	userAccount, err := pay.GetUserAmount(virUserAssetJson.UserId, virUserAssetJson.CoinType)
	if err != nil {
		logrus.Errorf("UpdateAdminUserVirtualAssetToDB GetUserAmount err: %s", err)
		return "error", model.STATUS_FAILED
	}
	// 添加虚拟账户余额
	newAccountAmount := utils.PayStringAdd(userAccount.VirtualAmount, virUserAssetJson.RechargeAmount, 2)
	// 添加用户余额表
	// 用户余额表
	virUserAsset := &model.VirtualRecharge{
		UserId:         virUserAssetJson.UserId,
		OutTradeNo:     outTradeNo,
		CreateAt:       createAt,
		CoinType:       virUserAssetJson.CoinType,
		CoinAmount:     newAccountAmount,
		OperatorId:     virUserAssetJson.OperatorId,
		RechargeAmount: virUserAssetJson.RechargeAmount,
		Description:    virUserAssetJson.Description,
	}
	// 事务操作
	tx := db.DB.Begin()
	// 新建一条虚拟充值记录
	if err := tx.Create(virUserAsset).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateAdminUserVirtualAssetToDB Create VirtualRecharge err %s", err)
		return "error", model.STATUS_FAILED
	}
	userAccounts := new(model.UserAccounts)
	// 用户 余额表 需要更新
	if err := tx.Model(userAccounts).Where("user_id = ? AND coin_type = ?", virUserAssetJson.UserId, virUserAssetJson.CoinType).Update("virtual_amount", newAccountAmount).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateAdminUserVirtualAssetToDB Create VirtualRecharge err %s", err)
		return "error", model.STATUS_FAILED
	}
	tx.Commit()
	// 添加用户账户流水
	userAssetflow := &model.UserAssetflow{
		UserId:      virUserAssetJson.UserId,
		OutTradeNo:  outTradeNo,
		TradeType:   model.TRADE_TYPE_INPUT,
		CreateAt:    createAt,
		CoinType:    virUserAssetJson.CoinType,
		Amount:      virUserAssetJson.RechargeAmount,
		TotalAmount: newAccountAmount,
		Description: "系统充值",
	}
	if err := db.DB.Create(userAssetflow).Error; err != nil {
		logrus.Errorf("UpdateWalletRecord service.DB.Create err %s, coinType: %s, userid: %d", err.Error(), userAssetflow.CoinType, userAssetflow.UserId)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 工具函数
func GetNextTime(timeStr string) (int64, error) {
	if len(timeStr) == 0 {
		return 0, errors.New("请选择报表截止时间")
	}
	nowTime := time.Now().Unix()
	endTimeStr, err := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	if err != nil {
		return 0, err
	}
	endTimeUnix := endTimeStr.AddDate(0, 0, 1).Unix()
	if endTimeUnix >= nowTime {
		return 0, errors.New(timeStr + ",结束时间不能大于当天零点时间")
	}
	return endTimeUnix, nil
}

// 创建报表
func CreateAdminReport(endTime int64, coinType string, authorId int) error {
	startTime, totalAmount, err := GetStartTimeAndTotalAmount(coinType)
	if err != nil {
		return err
	}
	if startTime >= endTime {
		return errors.New("截止时间小于起始时间")
	}
	res, count, err := GetAdminAssetFlowTimeIntervalFromDB(startTime, endTime, coinType)
	if err != nil {
		return err
	}
	increaceAmount := "0"
	reduceAmount := "0"
	totalTime := endTime - startTime
	flowId := uint(0)
	pointNumber := 8
	if coinType == "CNY" {
		pointNumber = 2
	}
	if count != 0 {
		flowId = res[count-1].ID
		for _, value := range res {
			if value.TradeType == model.TRADE_TYPE_INPUT {
				increaceAmount = utils.PayStringAdd(increaceAmount, value.Amount, pointNumber)
			} else if value.TradeType == model.TRADE_TYPE_SYS_OUT {
				reduceAmount = utils.PayStringAdd(reduceAmount, value.Amount, pointNumber)
			}
		}
		totalAmount = utils.PayStringAdd(totalAmount, increaceAmount, pointNumber)
		totalAmount, err = utils.PayStringSub(totalAmount, reduceAmount, pointNumber)
		if err != nil {
			return errors.New("可用余额为负,生成报表出错")
		}
	}
	createdAt := utils.GetNowTime()
	assets := &model.Assets{
		FlowId:         flowId,
		CoinType:       coinType,
		Count:          count,
		IncreaceAmount: increaceAmount,
		ReduceAmount:   reduceAmount,
		TotalAmount:    totalAmount,
		StartAt:        startTime,
		EndAt:          endTime,
		TotalTime:      totalTime,
		AuthorId:       authorId,
		CreateAt:       createdAt,
		Status:         false,
	}
	if err := db.DB.Create(assets).Error; err != nil {
		logrus.Errorf("AddAdminFinanaceReportToDB failed err: %s", err.Error())
		return err
	}
	sysAccount := &model.SystemAccount{
		CoinType:   coinType,
		CoinAmount: totalAmount,
		UpdateAt:   createdAt,
	}
	// 更新系统资产
	if err := CheckExistAndSysAccount(sysAccount); err != nil {
		logrus.Errorf("UpdateSysAccountToDB CheckExistAndSysAccount failed err: %s", err.Error())
		return err
	}
	sysAccountInfo := new(model.SystemAccount)
	if err := db.DB.Model(sysAccountInfo).Where("coin_type = ?", sysAccount.CoinType).Updates(map[string]interface{}{"coin_amount": sysAccount.CoinAmount, "update_at": sysAccount.UpdateAt}).Error; err != nil {
		logrus.Errorf("UpdateSysAccountToDB updates failed err: %s", err.Error())
	}
	logrus.Infof("UpdateSysAccountToDB coin_type: %s amount: %s \n", sysAccount.CoinType, sysAccount.CoinAmount)
	return nil
}

// 获取最后一次报表统计时间以及总金额
func GetStartTimeAndTotalAmount(coinType string) (int64, string, error) {
	asset, err := GetLastFinanceReportFromDB(coinType)
	// 如果未找到，设置初试时间为0，找到就设置为初始时间
	if err != nil {
		return 0, "", err
	}
	if err == nil && asset == nil {
		return 0, "0", nil
	}
	return asset.EndAt, asset.TotalAmount, nil
}

// 获取最后一条报表数据
func GetLastFinanceReportFromDB(coinType string) (*model.Assets, error) {
	asset := new(model.Assets)
	if err := db.DB.Where("coin_type = ?", coinType).Last(&asset).RecordNotFound(); err == true {
		logrus.Infof("GetLastFinanceReportFromDB  not found cointype: %s", coinType)
		return nil, nil
	}
	if err := db.DB.Last(&asset).Error; err != nil {
		logrus.Errorf("GetLastFinanceReportFromDB  Err: %s", err.Error())
		return nil, err
	}
	if asset.Status != true {
		return nil, errors.New("存在未审核报表")
	}
	return asset, nil
}

// 获取资产流水结果
func GetAdminAssetFlowTimeIntervalFromDB(startTime, endTime int64, coinType string) ([]model.UserAssetflow, int, error) {
	var userAccountsList []model.UserAssetflow
	count := 0
	dbs := db.DB
	if startTime != 0 {
		dbs = dbs.Where("create_at >= ?", startTime)
	}
	if endTime != 0 {
		dbs = dbs.Where("create_at < ?", endTime)
	}
	userAssetFlow := new(model.UserAssetflow)
	if err := dbs.Model(userAssetFlow).Where("coin_type = ?", coinType).
		Find(&userAccountsList).Error; err != nil {
		logrus.Errorf("GetAdminAssetFlowTimeIntervalFromDB failed err: %s", err.Error())
		return nil, count, errors.New("no data")
	}
	count = len(userAccountsList)
	logrus.Debugf("GetAdminAssetFlowTimeIntervalFromDB start_at: %d  end_at: %d coin_type: %s count: %d\n", startTime, endTime, coinType, count)
	return userAccountsList, count, nil
}

// 创建系统资产表
func CheckExistAndSysAccount(sysAccount *model.SystemAccount) error {
	sysAccountInfo := new(model.SystemAccount)
	if err := db.DB.Where("coin_type = ?", sysAccount.CoinType).First(sysAccountInfo).RecordNotFound(); err == true {
		// 创建 SystemAccount
		if err := db.DB.Create(sysAccount).Error; err != nil {
			logrus.Errorf("CheckExistAndSysAccount SystemAccount failed err: %s", err.Error())
			return err
		}
		return errors.New("aleady found")
	}
	return nil
}
