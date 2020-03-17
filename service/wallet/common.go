package wallet

import (
	"errors"
	"strings"

	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/cache"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

const (
	USDT_PRICE_DEFAULT = "6.80"
	MGD_PRICE_DEFAULT  = "0.01"
	ETH_PRICE_DEFAULT  = "1.00"
)
const (
	USDT_CASH_KEY = "USDT_PRICE"
	MGD_CASH_KEY  = "MGD_PRICE"
	ETH_CASH_KEY  = "ETH_PRICE"
)

// 获取usdt价格
func GetUsdtPrice() string {
	value, found := cache.OC.Get(USDT_CASH_KEY)
	if found {
		return value.(string)
	}
	coinPrice, err := GetCoinPriceFromDB("USDT")
	if err != nil {
		cache.OC.Set(USDT_CASH_KEY, USDT_PRICE_DEFAULT, cache.CacheDefaultExpiration*2)
		return USDT_PRICE_DEFAULT
	}
	cache.OC.Set(USDT_CASH_KEY, coinPrice.Price, cache.CacheDefaultExpiration*2)

	return coinPrice.Price
}

// 获取mgd价格
func GetMgdPrice() string {
	value, found := cache.OC.Get(MGD_CASH_KEY)
	if found {
		return value.(string)
	}
	coinPrice, err := GetCoinPriceFromDB("MGD")
	if err != nil {
		cache.OC.Set(MGD_CASH_KEY, MGD_PRICE_DEFAULT, cache.CacheDefaultExpiration*2)
		return MGD_PRICE_DEFAULT
	}
	cache.OC.Set(MGD_CASH_KEY, coinPrice.Price, cache.CacheDefaultExpiration*2)

	return coinPrice.Price
}

// 获取eth价格
func GetEthPrice() string {
	value, found := cache.OC.Get(ETH_CASH_KEY)
	if found {
		return value.(string)
	}
	coinPrice, err := GetCoinPriceFromDB("ETH")
	if err != nil {
		cache.OC.Set(ETH_CASH_KEY, MGD_PRICE_DEFAULT, cache.CacheDefaultExpiration*2)
		return ETH_PRICE_DEFAULT
	}
	cache.OC.Set(ETH_CASH_KEY, coinPrice.Price, cache.CacheDefaultExpiration*2)

	return coinPrice.Price
}

// 获取指定币种价格结构体
func GetCoinPriceWithCoinName(coinNameParam string) (*model.CoinPriceRes, string, int) {
	coinName := strings.ToUpper(coinNameParam)
	if coinName == "" {
		return nil, "error id", model.STATUS_FAILED
	}
	var coinPrice = "0"
	switch coinName {
	case "MGD":
		coinPrice = GetMgdPrice()
	case "ETH":
		coinPrice = GetEthPrice()
	case "USDT":
		coinPrice = GetUsdtPrice()
	case "CNY":
		coinPrice = "1"
	}
	if coinPrice != "0" {
		response := &model.CoinPriceRes{
			Name:  coinName,
			Price: coinPrice,
			Time:  utils.GetNowTime(),
		}
		return response, "success", model.STATUS_SUCCESS
	}
	res, err := GetCoinPriceFromDB(coinName)
	if err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}
	response := &model.CoinPriceRes{
		Name:  res.Name,
		Price: res.Price,
		Time:  res.Time,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取所有币种价格
func GetCoinPrices() (model.MyMap, string, int) {
	var coinPriceList []model.CoinPrice
	var coinPriceListRes []model.CoinPriceJsonRes
	count := 0
	coinPrice := new(model.CoinPrice)
	if err := db.DB.Model(coinPrice).Count(&count).
		Find(&coinPriceList).Scan(&coinPriceListRes).Error; err != nil {
		logrus.Errorf("GetACoinPriceListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(coinPriceListRes)
	logrus.Debugf("GetACoinPriceListFromDB count: %d\n", count)
	response := model.MyMap{
		"data":   coinPriceListRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取指定币种价格
func GetCoinPriceFromDB(coinName string) (*model.CoinPrice, error) {
	var coinPrice model.CoinPrice
	if err := db.DB.Where("name = ? ", coinName).Last(&coinPrice).Error; err == nil {
		logrus.Debugf("GetCoinPriceFromUser Name %s CoinPrice %s DisCount %s Time %s", coinPrice.Name, coinPrice.Price, coinPrice.Discount, utils.TimeStFormat(coinPrice.Time))
		return &coinPrice, nil
	}

	coinMarket := new(model.CoinMarket)
	if err := db.DB.Where("name = ? ", coinName).Last(coinMarket).Error; err != nil {
		logrus.Errorf("GetCoinPriceFromDB CoinMarket service.DB.Create err %s coinName: %s", err.Error(), coinName)
		return nil, err
	}
	logrus.Debugf("GetCoinPriceFromDB CoinMarket %s %s price: %s %s", coinMarket.Name, coinMarket.Url, coinMarket.Price, utils.TimeStFormat(coinMarket.Time))
	coinPrice.Name = coinMarket.Name
	coinPrice.Price = coinMarket.Price
	coinPrice.Time = coinMarket.Time

	return &coinPrice, nil
}

// 爬虫更新价格记录
func UpdateCoinMart(res *model.CoinMarket) error {
	if err := db.DB.Create(res).Error; err != nil {
		logrus.Errorf("UpdateCoinMart service.DB.Create err %s", err.Error())
		return err
	}
	if err := UpdateCoinPrice(res); err != nil {
		logrus.Errorf("UpdateCoinMart UpdateCoinPrice err %s", err.Error())
		return err
	}
	logrus.Infof("UpdateCoinMart %s %s price: %s %s", res.Name, res.Url, res.Price, utils.TimeStFormat(res.Time))
	return nil
}

// 更新价格
func UpdateCoinPrice(coinMarket *model.CoinMarket) error {
	coinPrice := &model.CoinPrice{
		Name:       coinMarket.Name,
		Price:      coinMarket.Price,
		Discount:   "1",
		AutoUpdate: true,
		Time:       coinMarket.Time,
	}
	if err := db.DB.Where("name = ?", coinMarket.Name).First(&coinPrice).RecordNotFound(); err == true {
		if err := db.DB.Create(coinPrice).Error; err != nil {
			logrus.Errorf("UpdateCoinPrice service.DB.Create err %s", err.Error())
			return err
		}
		return nil
	}
	if coinPrice.AutoUpdate == false {
		return nil
	}
	if err := db.DB.Model(coinPrice).Where("name = ?", coinMarket.Name).Updates(map[string]interface{}{"price": coinMarket.Price, "time": coinMarket.Time}).Error; err != nil {
		logrus.Errorf("UpdatePbStatus  update failed err: %s", err.Error())
		return err
	}
	logrus.Infof("UpdateCoinPrice Name: %s Price: %s Discount: %s AutoUpdate: %t Time: %s", coinPrice.Name, coinPrice.Price, coinPrice.Discount, coinPrice.AutoUpdate, utils.TimeStFormat(coinPrice.Time))
	return nil
}

// 更新用户余额
func UpdateUserAccounts(userId int, coinType, txid, addAmount string) (int, error) {
	// 查询资产更新记录表
	walletRecord := new(model.WalletRecord)
	if err := db.DB.Where("tx_id = ? AND added = ? ", txid, true).First(&walletRecord).Error; err == nil {
		return -1, errors.New("txid already added")
	}
	if addAmount == "0" {
		return -1, errors.New("amount is zero")
	}
	logrus.Infof("updateUserAccounts userId: %d, coinType: %s, coinAmount: %s", userId, coinType, addAmount)
	// 事务操作
	tx := db.DB.Begin()
	// 获取当前资产
	userAccount := new(model.UserAccounts)
	if err := tx.Where("user_id = ? AND coin_type = ? ", userId, coinType).First(&userAccount).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("updateUserAccounts service.DB get user accountAmount err %s", err.Error())
		return 0, err
	}
	// 增加余额
	accountAmount := utils.PayStringAdd(userAccount.CoinAmount, addAmount, 2) // (当前余额,增加的余额,保留小数位数)

	// 新建一条充值记录
	if err := tx.Model(&userAccount).Where("user_id = ? AND coin_type = ?", userId, coinType).Update("coin_amount", accountAmount).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("CreateDBPayRefundTx service.DB err %s", err.Error())
		return 0, err
	}
	// 添加资产流水记录
	payChannel, payType := GetCoinTypeFromStr(coinType)
	outTradeNo := utils.GenOutTradeNo(payChannel, payType, 0, userId)
	createAt := utils.GetNowTime()
	userAssetflow := &model.UserAssetflow{
		UserId:      userId,
		OutTradeNo:  outTradeNo,
		TradeType:   model.TRADE_TYPE_INPUT,
		CreateAt:    createAt,
		CoinType:    coinType,
		Amount:      addAmount,
		TotalAmount: accountAmount,
		Description: "区块检测入账",
	}
	if err := tx.Create(userAssetflow).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateWalletRecord service.DB.Create err %s, coinType: %s, userid: %d", err.Error(), userAssetflow.CoinType, userAssetflow.UserId)
		return 0, err
	}

	// 更新 钱包更新记录表
	walletRecord = &model.WalletRecord{
		TxId:     txid,
		CoinType: coinType,
		Added:    true,
	}
	if err := tx.Create(walletRecord).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("UpdateWalletRecord service.DB.Create err %s, coinType: %s, txid: %s", err.Error(), coinType, txid)
		return 0, err
	}
	tx.Commit()

	logrus.Infof("updateUserAccounts success userId: %d, coinType: %s, addAmount: %s, coinAmount: %s", userId, coinType, addAmount, userAccount.CoinAmount)
	return 0, nil
}

// 更新钱包记录表
func UpdateWalletSyncToDB(walletSync *model.WalletSync) error {
	wtSync := new(model.WalletSync)
	if err := db.DB.Model(wtSync).Where("coin_type = ? ", walletSync.CoinType).Updates(walletSync).Error; err != nil {
		logrus.Errorf("UpdateWalletSyncToDB failed err: %s", err.Error())
		return err
	}
	return nil
}

// ETH插入多组数据
func EthMutilyInsert(wtList []*model.EthTx) error {
	sqlStr := "INSERT INTO eth_txes (block_height,tx_hash,tx_index,coin_from,coin_to,tx_type,nonce,value,gas,gas_price,server_ip,time) VALUES "
	valueArgs := []interface{}{}
	const rowSQL = "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var inserts []string
	for _, elem := range wtList {
		inserts = append(inserts, rowSQL)
		valueArgs = append(valueArgs, elem.BlockHeight, elem.TxHash, elem.TxIndex, elem.CoinFrom, elem.CoinTo, elem.TxType, elem.Nonce, elem.Value, elem.Gas, elem.GasPrice, elem.ServerIP, elem.Time)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	err := db.DB.Exec(sqlStr, valueArgs...).Error
	return err
}

// MGD插入多组数据
func MGDMutilyInsert(wtList []*model.MgdTx) error {
	sqlStr := "INSERT INTO mgd_txes (tx_index,tx_id,to_address,tx_type,to_amount,fee,time,server_ip,detail) VALUES "
	valueArgs := []interface{}{}
	const rowSQL = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var inserts []string
	// coinType := "mgd"
	for _, elem := range wtList {
		inserts = append(inserts, rowSQL)
		valueArgs = append(valueArgs, elem.TxIndex, elem.TxId, elem.ToAddress, elem.TxType, elem.ToAmount, elem.Fee, elem.Time, elem.ServerIP, elem.Detail)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	err := db.DB.Exec(sqlStr, valueArgs...).Error
	return err
}

// BTC 插入多组数据
func BTCMutilyInsert(wtList []*model.BtcTx) error {
	sqlStr := "INSERT INTO btc_txes (tx_index,tx_id,to_address,tx_type,to_amount,fee,time,server_ip,detail) VALUES "
	valueArgs := []interface{}{}
	const rowSQL = "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var inserts []string
	// coinType := "mgd"
	for _, elem := range wtList {
		inserts = append(inserts, rowSQL)
		valueArgs = append(valueArgs, elem.TxIndex, elem.TxId, elem.ToAddress, elem.TxType, elem.ToAmount, elem.Fee, elem.Time, elem.ServerIP, elem.Detail)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	err := db.DB.Exec(sqlStr, valueArgs...).Error
	return err
}

// update wallet basic
func UpdateWalletBasicInfoToDB(walletInfo *model.WalletBasic) error {
	badicInfo := new(model.WalletBasic)
	if err := db.DB.Model(badicInfo).Where("name = ?", walletInfo.Name).Updates(walletInfo).Error; err != nil {
		logrus.Errorf("UpdateWalletBasicInfoToDB updates failed err: %s", err.Error())
	}
	logrus.Infof("UpdateWalletBasicInfoToDB name: %s", walletInfo.Name)
	return nil
}

// 检测钱包基础信息表
func CheckExistAndCreateWalletBasic(walletInfo *model.WalletBasic) error {
	badicInfo := new(model.WalletBasic)
	if err := db.DB.Where("name = ?", walletInfo.Name).First(badicInfo).RecordNotFound(); err != true {
		// 已经存在
		return nil
	}
	// 创建 BasicInfo
	if err := db.DB.Create(walletInfo).Error; err != nil {
		logrus.Errorf("CheckExistAndCreateWalletBasic WalletBasic failed err: %s", err.Error())
		return err
	}
	return nil
}

// add  wallet simple info
func AddWalletSimpleInfoToDB(walletInfo *model.WalletSimple) error {
	if err := db.DB.Create(walletInfo).Error; err != nil {
		logrus.Errorf("AddWalletSimpleInfoToDB create WalletSimple failed err: %s", err.Error())
		return err
	}
	logrus.Infof("AddWalletSimpleInfoToDB name: %s", walletInfo.Name)
	return nil
}

// 获取最后一条 交易txid
func GetLastWalletTxFromDB(coinType string) (string, int, error) {
	if coinType == "MGD" {
		var walletTx model.MgdTx
		if err := db.DB.Last(&walletTx).Error; err != nil {
			logrus.Errorf("GetLastWalletTx  Err: %s", err.Error())
			return "", 0, err
		}
		return walletTx.TxId, walletTx.TxIndex, nil
	}
	var btcTx model.BtcTx
	if err := db.DB.Last(&btcTx).Error; err != nil {
		logrus.Errorf("GetLastWalletTx  Err: %s", err.Error())
		return "", 0, err
	}
	return btcTx.TxId, btcTx.TxIndex, nil

}

// 获取当前钱包所有 用户地址
func GetAllUserCoinAddr(coinType string) (map[string]int, error) {
	// 初始化 map
	coinAddr := make(map[string]int)
	var userAccount []model.UserAccounts
	if err := db.DB.Select("user_id,coin_addr").Where("coin_type = ?", coinType).Find(&userAccount).Error; err != nil {
		logrus.Errorf("getAllUserCoinAddr get all user addr failed err: %s", err.Error())
		return nil, err
	}
	for _, value := range userAccount {
		coinAddr[value.CoinAddr] = value.UserId
		logrus.Debugf("getAllUserCoinAddr addr: %s %d", value.CoinAddr, value.UserId)
	}
	return coinAddr, nil
}

// 从钱包更新记录表获取最后更新记录
func GetLastWalletSyncRecordFromDB(coinType string) (int64, error) {
	var walletSync model.WalletSync
	if err := db.DB.Where("coin_type = ?", coinType).First(&walletSync).RecordNotFound(); err == true {
		// 创建 WalletSync
		walletSync.CoinType = coinType
		if err := db.DB.Create(&walletSync).Error; err != nil {
			logrus.Errorf("GetLastWalletSyncRecordFromDB WalletSync failed err: %s", err.Error())
			return 0, err
		}
		return 0, errors.New("not found")
	}
	if err := db.DB.Last(&walletSync).Error; err != nil {
		logrus.Errorf("GetLastWalletSyncRecordFromDB  Err: %s", err.Error())
		return 0, err
	}
	return walletSync.LastBlock, nil
}

// 获取支付方式的渠道类型
func GetCoinTypeFromStr(coinType string) (int, int) {
	switch coinType {
	case "CNY":
		return 0, model.COIN_PAYTYPE_CNY
	case "MGD":
		return model.COIN_CHANNEL_DIGITAL, model.COIN_PAYTYPE_MGD
	case "ETH":
		return model.COIN_CHANNEL_DIGITAL, model.COIN_PAYTYPE_ETH
	}
	return 0, 0
}
