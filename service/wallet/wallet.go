package wallet

import (
	"errors"
	"os"
	"time"

	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/set"
	"github.com/zhxx123/gomonitor/service/task"
	"github.com/sirupsen/logrus"
)

var mgdsync *MgdSync
var ethsync *EthSync

func InitWallet() {
	//初始化数据库
	InitWalletBasicDB()
	task.AddTimer(time.Second*60, SyncMGDInfo)
	logrus.Debug("InitWallet AddCallback SyncMGDInfo, update wallet, time duation 10s, repated")
	task.AddTimer(time.Second*120, SyncMGDMiningInfo)
	logrus.Debug("InitWallet AddCallback SyncMGDMiningInfo, update wallet, time duation 10s, repated")

	// mgd 同步策略
	mgdsync = InitMGDWallet()
	logrus.Debug("InitWallet AddCallback SyncMGD, init wallet, time duation 5s, only once")
	task.AddCallback(time.Second*5, SyncMGD)
	logrus.Debug("InitWallet AddCallback SyncUpdateTrans, update wallet, time duation 10s, repated")
	task.AddTimer(time.Second*10, SyncUpdateTrans)

	// eth 同步策略
	ethsync = InitETHWallet()
	logrus.Debug("InitWallet AddCallback SyncETH, init wallet, time duation 5s, only once")
	task.AddCallback(time.Second*25, SyncETH)
	logrus.Debug("InitWallet AddCallback SyncUpdateEthTrans, update wallet, time duation 10s, repated")
	task.AddTimer(time.Second*30, SyncUpdateEthTrans)
}

// 初始化数据库
func InitWalletBasicDB() {
	walletBasic := &model.WalletBasic{
		Name: "MGD",
	}
	if err := CheckExistAndCreateWalletBasic(walletBasic); err != nil {
		logrus.Errorf("Wallet Sync MGD Basic err: %s", err.Error())
		os.Exit(-1)
	}
	walletBasicETH := &model.WalletBasic{
		Name: "ETH",
	}
	if err := CheckExistAndCreateWalletBasic(walletBasicETH); err != nil {
		logrus.Errorf("Wallet Sync ETH Basic err: %s", err.Error())
		os.Exit(-1)
	}

	walletBasicBTC := &model.WalletBasic{
		Name: "BTC",
	}
	if err := CheckExistAndCreateWalletBasic(walletBasicBTC); err != nil {
		logrus.Errorf("Wallet Sync BTC Basic err: %s", err.Error())
		os.Exit(-1)
	}
	logrus.Debug("InitWallet InitWalletBasicDB MGD, ETH, BTC BasicDB success")
}

// 初始化同步MGD
func SyncMGD(param map[string]interface{}) bool {
	return mgdsync.syncMGDWallet()
}

// 同步 MGD 交易
func SyncUpdateTrans(param map[string]interface{}) bool {
	return mgdsync.UpdateTransaction()
}

// 同步MGD 钱包基本信息
func SyncMGDInfo(param map[string]interface{}) bool {
	return mgdsync.syncMGDWalletInfo()
}

// 同步MGD 钱包状态信息
func SyncMGDMiningInfo(param map[string]interface{}) bool {
	return mgdsync.syncMGDMiningInfo()
}

// eth
// 初始化同步ETH
func SyncETH(param map[string]interface{}) bool {
	return ethsync.syncETHWallet()
}

// 同步 ETH 交易
func SyncUpdateEthTrans(param map[string]interface{}) bool {
	return ethsync.UpdateTransaction()
}

// 更新同步参数，区块确认数
func UpdateWalletConfirm() bool {
	mgdConfirm, _ := set.GetConfMapKeyValue("coin_confirm_mgd", 3.0)
	mgdConf := int(mgdConfirm.(float64))
	ethConfirm, _ := set.GetConfMapKeyValue("coin_confirm_eth", 12.0)
	ethConf := int(ethConfirm.(float64))
	if mgdConf <= 0 && mgdConf >= 20 {
		return false
	}
	if ethConf <= 0 && ethConf >= 60 {
		return false
	}
	if mgdsync != nil {
		mgdconfirm := mgdsync.getWalletConfirm()
		if mgdConf != mgdconfirm {
			mgdsync.setWalletConfirm(mgdConf)
		}
	}
	if ethsync != nil {
		ethconfirm := ethsync.getWalletConfirm()
		if int64(ethConf) != ethconfirm {
			ethsync.setWalletConfirm(int64(ethConf))
		}
	}
	logrus.Infof("UpdateWalletConfirm mgd: %d eth: %d", mgdConf, ethConf)
	return true
}

func UpdateWalletAddress(coinType, coinAddr string, userId int) error {
	if coinType == "MGD" {
		if mgdsync != nil {
			if status := mgdsync.setWalletAddress(userId, coinAddr); status != true {
				return errors.New("data error")
			}
			return nil
		}
		return errors.New("info null pointer mgdsync")
	}
	if coinType == "ETH" {
		if ethsync != nil {
			if status := ethsync.setWalletAddress(userId, coinAddr); status != true {
				return errors.New("data error")
			}
			return nil
		}
		return errors.New("info null pointer ethsync")
	}
	return errors.New("unknow cointype")
}
