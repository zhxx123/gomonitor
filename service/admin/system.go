package admin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/set"
	"github.com/zhxx123/gomonitor/service/user"
	"github.com/zhxx123/gomonitor/service/wallet"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

// 获取 basic 列表 信息
func CheckAndAdminBasicListData(simpleLogJson *model.SysInfoJson) (model.MyMap, string, int) {
	var basicInfo []model.SystemBasic
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(simpleLogJson.Page, simpleLogJson.Limit)
	if err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}
	if len(simpleLogJson.UID) > 0 {
		dbs = dbs.Where("uid = ?", simpleLogJson.UID)
	}
	sySimple := new(model.SystemBasic)
	if err := dbs.Model(sySimple).
		Offset(offset).Limit(simpleLogJson.Limit).
		Find(&basicInfo).Error; err != nil {
		logrus.Errorf("GetAdminSysBasicListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(basicInfo)
	logrus.Debugf("GetAdminSysBasicListFromDB count: %d uid: %s page: %d limit: %d\n", count, simpleLogJson.UID, simpleLogJson.Page, simpleLogJson.Limit)

	response := model.MyMap{
		"data":   basicInfo,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// CheckAndAdminSimpleData 获取系统 Simple 信息
func CheckAndAdminSimpleData(simpleLogJson *model.SysInfoJson) (model.MyMap, string, int) {

	var simpleInfo []model.SystemSimple
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(simpleLogJson.Page, simpleLogJson.Limit)
	if err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}
	if len(simpleLogJson.UID) > 0 {
		dbs = dbs.Where("uid = ?", simpleLogJson.UID)
	}
	sySimple := new(model.SystemSimple)
	if err := dbs.Model(sySimple).
		Offset(offset).Limit(simpleLogJson.Limit).
		Find(&simpleInfo).Error; err != nil {
		logrus.Errorf("GetAdminSysSimpleFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(simpleInfo)
	logrus.Debugf("GetAdminSysSimpleFromDB count: %d uid: %s %d page: %d limit: %d\n", count, simpleLogJson.UID, len(simpleLogJson.UID), simpleLogJson.Page, simpleLogJson.Limit)

	response := model.MyMap{
		"data":   simpleInfo,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminWalletBasic 获取钱包 basic信息
func GetAdminWalletBasic(walletJson *model.AWalletInfoJson) (model.MyMap, string, int) {
	var walletBasic []model.WalletBasic
	count := 0
	dbs := db.DB
	if len(walletJson.Name) > 0 {
		dbs = dbs.Where("name = ?", strings.ToUpper(walletJson.Name))
	}
	wBasic := new(model.WalletBasic)
	if err := dbs.Model(wBasic).
		Find(&walletBasic).Error; err != nil {
		logrus.Errorf("GetAdminWalletBasicFromDB failed err: %s Name: %s", err.Error(), walletJson.Name)
		return nil, "error", model.STATUS_FAILED
	}
	count = len(walletBasic)
	response := model.MyMap{
		"data":   walletBasic,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminWalletSimple 获取钱包状态信息
func GetAdminWalletSimple(walletJson *model.AWalletInfoJson) (model.MyMap, string, int) {
	var walletBasic []model.WalletBasic
	count := 0
	dbs := db.DB
	if len(walletJson.Name) > 0 {
		dbs = dbs.Where("name = ?", strings.ToUpper(walletJson.Name))
	}
	wBasic := new(model.WalletBasic)
	if err := dbs.Model(wBasic).
		Find(&walletBasic).Error; err != nil {
		logrus.Errorf("GetAdminWalletBasicFromDB failed err: %s Name: %s", err.Error(), walletJson.Name)
		return nil, "error", model.STATUS_FAILED
	}
	count = len(walletBasic)
	response := model.MyMap{
		"data":   walletBasic,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminWalletAddressRecord 获取钱包地址分配状况
func GetAdminWalletAddressRecord(walletJson *model.AWalletAddressJson) (model.MyMap, string, int) {
	var walletAddress []model.WalletAddress
	count := 0
	dbs := db.DB
	if walletJson.IsStatus == true {
		dbs = dbs.Where("allocated = ?", walletJson.Status)
	}
	if len(walletJson.CoinType) > 0 {
		dbs = dbs.Where("coin_type = ?", walletJson.CoinType)
	}
	offset, err := db.GetOffset(walletJson.Page, walletJson.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if err := dbs.Model(model.WalletAddress{}).Count(&count).
		Offset(offset).Limit(walletJson.Limit).
		Find(&walletAddress).Error; err != nil {
		logrus.Errorf("GetAdminWalletAddressRecordFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminWalletAddressRecordFromDB status: %t cointype: %s page: %d limit %d offset: %d\n", walletJson.Status, walletJson.CoinType, walletJson.Page, walletJson.Limit, offset)
	response := model.MyMap{
		"data":   walletAddress,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UpdateAdminWalletAddress 上传钱包地址
func UpdateAdminWalletAddress(walletJson *model.AUpdateWalletAddressJson) (string, int) {
	// insert mutil data
	time := time.Now()
	sqlStr := "INSERT INTO wallet_addresses (created_at,updated_at,deleted_at,indexs,coin_type,account,address,user_id,allocated) VALUES "
	valueArgs := []interface{}{}
	const rowSQL = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var inserts []string
	for _, elem := range walletJson.Address {
		inserts = append(inserts, rowSQL)
		valueArgs = append(valueArgs, time, time, nil, elem.Indexs, elem.CoinType, elem.Account, elem.Address, 0, 0)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	if err := db.DB.Exec(sqlStr, valueArgs...).Error; err != nil {
		return "failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 更新用户地址分配
func UpdateAdminUserNewCoinAddress(walletJson *model.AUserCoinAddressJson) (string, int) {
	coinAddr := user.GetCoinAddressFromDB(walletJson.CoinType)
	if coinAddr == "" {
		return "没有可用的地址,请稍后再试", model.STATUS_FAILED
	}
	tx := db.DB.Begin()
	userAcc := new(model.UserAccounts)
	if err := tx.Model(userAcc).Where("user_id = ? AND coin_type = ?", walletJson.UserId, walletJson.CoinType).Update("coin_addr", coinAddr).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateUserAccountAddress falied Err:%s", err.Error())
		return "error", model.STATUS_FAILED
	}
	// 更新地址分配记录
	wtAddress := new(model.WalletAddress)
	if err := tx.Model(wtAddress).Where("coin_type = ? AND address = ?  AND allocated = ?", walletJson.CoinType, coinAddr, false).
		Update(map[string]interface{}{"user_id": walletJson.UserId, "allocated": true}).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateCoinAddressToDB falied Err:%s", err.Error())
		return "error", model.STATUS_FAILED
	}
	tx.Commit()
	return "success", model.STATUS_SUCCESS
}

// 获取资产列表
func GetAdminCoinPriceList() (model.MyMap, string, int) {
	var coinPriceList []model.CoinPrice
	count := 0
	if err := db.DB.Model(model.CoinPrice{}).Count(&count).
		Find(&coinPriceList).Error; err != nil {
		logrus.Errorf("GetACoinPriceList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(coinPriceList)
	response := model.MyMap{
		"data":   coinPriceList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取数字货币市场价格表
func GetAdminCoinMarket(coinMarketJson *model.ACoinMarketJson) (model.MyMap, string, int) {
	var coinMarketList []model.CoinMarket
	count := 0
	offset, err := db.GetOffset(coinMarketJson.Page, coinMarketJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := db.DB.Model(model.CoinMarket{}).Order("id desc").
		Offset(offset).Limit(coinMarketJson.Limit).
		Find(&coinMarketList).Error; err != nil {
		logrus.Errorf("GetACoinMarketList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(coinMarketList)
	response := model.MyMap{
		"data":   coinMarketList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UpdateAdminCoinPrices 更新数字货币价格
func UpdateAdminCoinPrices(ctx iris.Context, coinPriceJson *model.AUpdateCoinPriceJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	coinPrice := new(model.CoinPrice)
	if err := db.DB.Model(coinPrice).Where("name = ?", coinPriceJson.CoinType).
		Updates(map[string]interface{}{"price": coinPriceJson.CoinPrice, "discount": coinPriceJson.Discount, "auto_update": false}).Error; err != nil {
		logrus.Errorf("UpdateAdminCoinPriceToDB falied Err:%s", err.Error())
		return "error", model.STATUS_FAILED
	}
	logrus.Debugf("UpdateAdminCoinPriceStatusToDB name: %s price: %s discount: %s auto_update: 0\n", coinPriceJson.CoinType, coinPriceJson.CoinPrice, coinPriceJson.Discount)
	// 添加更新记录
	nowTime := utils.GetNowTime()
	url := fmt.Sprintf("admin [id:%d]", userId)
	coinMarket := &model.CoinMarket{
		Name:  coinPriceJson.CoinType,
		Url:   url,
		Price: coinPriceJson.CoinPrice,
		Time:  nowTime,
	}
	if err := wallet.UpdateCoinMart(coinMarket); err != nil {
		return "UpdateCoinMart Failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminCoinPriceStatus 更新数字货币 自动更新状态
func UpdateAdminCoinPriceStatus(coinPriceJson *model.AUpdateStatusCoinPriceJson) (string, int) {
	coinPrice := new(model.CoinPrice)
	if err := db.DB.Model(coinPrice).Where("name = ?", coinPriceJson.CoinType).
		Update("auto_update", coinPriceJson.Status).Error; err != nil {
		logrus.Errorf("UpdateAdminCoinPriceStatusToDB falied Err:%s", err.Error())
		return "error", model.STATUS_FAILED
	}
	logrus.Debugf("UpdateAdminCoinPriceStatusToDB name: %s auto_update: %t\n", coinPriceJson.CoinType, coinPriceJson.Status)
	return "success", model.STATUS_SUCCESS
}

// GetAdminUserSettingsList 获取系统设置记录表
func GetAdminUserSettingsList(settingJson *model.ASettingsJson) (*model.Settings, string, int) {
	var settings model.Settings
	if err := db.DB.Where("category = ? AND name = ?", settingJson.Category, settingJson.Name).
		First(&settings).Error; err != nil {
		logrus.Errorf("GetAdminSettingsFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	return &settings, "success", model.STATUS_SUCCESS
}

// UpdateAdminSettings 更新系统设置表
func UpdateAdminSettings(ctx iris.Context, settingJson *model.ASettingsUpdateJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	settings := &model.Settings{
		AuthorId: userId,
		Category: settingJson.Category,
		Name:     settingJson.Name,
		Value:    settingJson.Value,
	}
	if err := UpdateAdminSettingsToDB(settings); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 刷新全局缓存
	set.InitSysSeting()          // 系统设置
	set.UpdateUserSafeSeting()   // 更新安全设置
	wallet.UpdateWalletConfirm() // 更新交易设置
	return "success", model.STATUS_SUCCESS
}

// GetUserLoginList 获取用户登录记录
func GetUserAdminLoginList(ctx iris.Context, userOauthJson *model.UserOauthJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return nil, "error", model.STATUS_FAILED
	}

	var oauthList []model.UserOauth
	var oauthRes []model.UserOauthRes
	count := 0
	offset, err := db.GetOffset(userOauthJson.Page, userOauthJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := db.DB.Model(model.UserOauth{}).Order("id desc").Where("user_id =?", userId).
		Offset(offset).Limit(userOauthJson.Limit).Find(&oauthList).Scan(&oauthRes).Error; err != nil {
		logrus.Errorf("GetUserLoginListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(oauthRes)

	response := model.MyMap{
		"data":   oauthRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 保存设置
func UpdateAdminSettingsToDB(setting *model.Settings) error {
	// 如果未找到，则创建
	sets := new(model.Settings)
	if err := db.DB.Where("category = ? AND name = ?", setting.Category, setting.Name).First(sets).RecordNotFound(); err == true {
		if err := db.DB.Create(setting).Error; err != nil {
			logrus.Errorf("UpdateAdminSettingsToDB  create setting failed err: %s", err.Error())
			return errors.New("error")
		}
		return nil
	}
	if err := db.DB.Model(sets).Updates(setting).Error; err != nil {
		logrus.Errorf("UpdateAdminSettingsToDB  failed err: %s", err.Error())
		return errors.New("error")
	}
	return nil
}

// web hook
func WebHookLog(queryData *model.QueryHookJson) (model.MyMap, string, int) {
	var webHook []model.WebHook
	count := 0
	dbs := db.DB
	if queryData.ProjectName != "" {
		dbs = dbs.Where("project_name = ?", queryData.ProjectName)
	}
	offset, err := db.GetOffset(queryData.Page, queryData.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if err := dbs.Model(model.WebHook{}).Order("ID desc").
		Count(&count).Offset(offset).Limit(queryData.Limit).
		Find(&webHook).Error; err != nil {
		logrus.Errorf("WebHook failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	// count = len(helpArticleCategory)
	response := model.MyMap{
		"data":   webHook,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}
