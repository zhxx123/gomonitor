package wallet

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/wallet/ethrpc"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

var ETH_WEI = big.NewInt(10000000000)

var ETH_WEI_ZERO = big.NewInt(0)

type EthSync struct {
	rpc         *ethrpc.EtherRPC
	url         string
	coinAddr    map[string]int
	serverIP    string // 节点服务器钱包
	coinType    string // 货币类型，ETH
	blockNumber int64  // 最后更新的区块高度
	confirm     int64  // 最低确认块数
	update      bool   //当前是否正在更新
}

// var ethsync *EthSync

// init to sync wallet
func InitETHWallet() *EthSync {
	rpcUrl := config.RPCConfig.ETHHost
	confirmStr := config.RPCConfig.ETHConfirm
	confirm := int64(confirmStr)
	if confirm <= 0 {
		confirm = 12
	}
	coinType := "ETH"
	coinAddr, err := GetAllUserCoinAddr(coinType)
	if err != nil {
		return nil
	}
	logrus.Infof("InitETHWallet GetAllUserCoinAddr cointype %s %d", coinType, len(coinAddr))
	client := NewEthRpcClient(rpcUrl)
	ethsync := NewEthSync(client, rpcUrl, coinAddr, confirm)
	return ethsync
}

func NewEthRpcClient(url string) *ethrpc.EtherRPC {
	logrus.Debug("New EthRpcClient:", url)
	client := ethrpc.NewEtherRPC(url, ethrpc.WithETHERDebug(false))
	return client
}

func NewEthSync(client *ethrpc.EtherRPC, url string, coinAddr map[string]int, confirm int64) *EthSync {
	rpc := &EthSync{
		rpc:         client,
		coinAddr:    coinAddr,
		url:         url,
		coinType:    "ETH",
		blockNumber: 0,
		confirm:     confirm,
	}
	return rpc
}
func (eth *EthSync) setWalletConfirm(confirm int64) {
	eth.confirm = confirm
}
func (eth *EthSync) getWalletConfirm() int64 {
	return eth.confirm
}
func (eth *EthSync) setWalletAddress(userId int, coinaddr string) bool {
	if _, ok := eth.coinAddr[coinaddr]; ok != true {
		eth.coinAddr[coinaddr] = userId
		return true
	}
	return false
}

// etherGetBlockByNumber 获取当前最新区块
func (eth *EthSync) etherGetBlockByNumber(blockHeight int64, withTransactions bool) (*ethrpc.Block, error) {
	if eth.rpc != nil {
		return eth.rpc.EthGetBlockByNumber(blockHeight, withTransactions)
	}
	return nil, errors.New("eth.client getbalance error")
}

// etherEthBlockNumber 获取当前最新区块
func (eth *EthSync) etherEthBlockNumber() (int64, error) {
	if eth.rpc != nil {
		return eth.rpc.EthBlockNumber()
	}
	return 0, errors.New("eth.client eth_blockNumber error")
}

func (eth *EthSync) syncETHWallet() bool {
	lastBlockHeight, err := GetLastWalletSyncRecordFromDB(eth.coinType)

	if err != nil {
		logrus.Debugf("lastBlockHeight %d err: %s", lastBlockHeight, err.Error())
		lastBlockHeight, err = eth.etherEthBlockNumber() // 直接获取的最后一个区块，不用+1
	} else {
		lastBlockHeight = lastBlockHeight + 1 //从已知更新的最后一个区块+1开始
	}
	eth.blockNumber = lastBlockHeight
	ip := strings.Split(eth.url, "/")
	serverip := "default"
	if len(ip) == 4 {
		serverip = ip[2]
	}
	logrus.Infof("syncETHWallet LastBlockHeight: %d ip: %s", lastBlockHeight, serverip)
	eth.serverIP = serverip
	return true
}

// start to update transactions once
func (eth *EthSync) UpdateTransaction() bool {
	// 获取最新区块高度
	lastBloclHeight, err := eth.etherEthBlockNumber()
	if err != nil {
		logrus.Errorf("UpdateTransaction etherEthBlockNumber startBlockHeight error: %s", err)
		return false
	}
	if eth.update == true {
		logrus.Debugf("UpdateTransaction: already updated LastBlockHeight %d", lastBloclHeight)
		return false
	}
	eth.update = true
	startBlockHeight := eth.blockNumber
	for {
		logrus.Debugf("UpdateTransaction: startBlockHeight %d LastBlockHeight %d confirm: %d", startBlockHeight, lastBloclHeight, eth.confirm)
		if lastBloclHeight <= startBlockHeight+eth.confirm {
			// 当前未满足条件
			break
		}
		// 获取区块数据
		blocks, err := eth.etherGetBlockByNumber(startBlockHeight, true)
		if err != nil {
			logrus.Errorf("etherGetBlockByNumber error: %s", err)
			break
		}
		blockNumber := blocks.Number
		if blockNumber != startBlockHeight {
			logrus.Errorf("etherGetBlockByNumber error blocks, height: %d expect height: %d", blockNumber, startBlockHeight)
			continue
		}

		transList := blocks.Transactions
		// listtrans size
		// transLength := len(transList)
		var wtList []*model.EthTx
		for _, result := range transList {
			// logrus.Debugf("index: %d %+v\n", index, result)
			transacion, err := eth.checkTransactions(result, blocks.Timestamp)
			if err != nil { //表示数据库更新出错
				eth.update = false
				return false
			}
			if transacion == nil {
				continue
			}

			wtList = append(wtList, transacion)
		}
		// 更新最新同步区块个数 wallet_syncs
		times := utils.GetNowTime()
		walletSync := &model.WalletSync{
			CoinType:  eth.coinType,
			LastBlock: startBlockHeight,
			UpdateAt:  times,
		}
		if err := UpdateWalletSyncToDB(walletSync); err != nil {
			logrus.Errorf("UpdateTransaction: UpdateWalletSyncToDB succefully lastBlock %d err: %s", startBlockHeight, err.Error())
		}
		startBlockHeight = startBlockHeight + 1
		eth.blockNumber = startBlockHeight
		if len(wtList) > 0 {
			err := EthMutilyInsert(wtList)
			if err != nil {
				logrus.Errorf("UpdateTransaction: EthMutilyInsert failed! lastblockheight: %d err: %s", startBlockHeight-1, err.Error())
			}
			logrus.Infof("UpdateTransaction: EthMutilyInsert succefully lastblockheight %d", startBlockHeight-1)
		}
		time.Sleep(500 * time.Millisecond)
	}
	eth.update = false
	// defer CloseRPCClient()
	return false
}

// generate one sql
func (eth *EthSync) checkTransactions(detail ethrpc.Transaction, timestamp int64) (*model.EthTx, error) {

	toaddress := detail.To
	fromaddress := detail.From
	if len(toaddress) == 0 || len(fromaddress) == 0 {
		return nil, errors.New("empty address")
	}
	txType := ""
	toamount := "0"
	value_to, ok_to := eth.coinAddr[toaddress]
	if ok_to == true {
		txType = "reveive"
	}
	_, ok_from := eth.coinAddr[fromaddress]
	if ok_from == true {
		if txType != "" {
			txType = "self"
		} else {
			txType = "send"
		}
	}
	if txType != "" {
		detail.Value.Div(&detail.Value, ETH_WEI)
		if detail.Value.Cmp(ETH_WEI_ZERO) > 0 {
			bigStr := detail.Value.String()
			toamount = utils.BigIntDivToString(bigStr, 8) //移动八位小数
		}
		if txType == "recevice" || txType == "self" { //只有是收款时候才更新用户账户
			if status, err := UpdateUserAccounts(value_to, eth.coinType, detail.Hash, toamount); err != nil {
				logrus.Errorf("checkTransactions updateUserAccounts err: %s userid: %d cointype: %s addr: %s txtype: %s amount: %s", err.Error(), value_to, eth.coinType, toaddress, txType, toamount)
				if status != -1 { // 表示数据库出错，直接退出当次同步
					return nil, errors.New("update error")
				}
			}
			logrus.Infof("checkTransactions updateUserAccounts receive toaddress %s userid: %d cointype: %s txtype: %s amount: %s updateted", toaddress, value_to, eth.coinType, txType, toamount)
		}
		gasPriceStr := detail.GasPrice.String()
		gasPriceStr = utils.BigIntDivToString(gasPriceStr, 9)
		ethTx := &model.EthTx{
			BlockHeight: int64(*detail.BlockNumber),
			TxHash:      detail.Hash,
			TxIndex:     *detail.TransactionIndex,
			CoinFrom:    detail.From,
			CoinTo:      detail.To,
			TxType:      txType,
			Nonce:       detail.Nonce,
			Value:       toamount,
			Gas:         detail.Gas,
			GasPrice:    gasPriceStr,
			ServerIP:    eth.serverIP,
			Time:        timestamp,
		}
		return ethTx, nil
	}
	return nil, nil
}

// ETH钱包基本信息同步
func (eth *EthSync) syncETHWalletInfo() bool {
	// res, err := eth.getinfo()
	// if err != nil {
	// 	logrus.Errorf("getinfo error %s", err.Error())
	// 	return false
	// }
	// version := util.ParseWalletVersion(res.Version)

	// netModel := "Main"
	// if res.Testnet == true {
	// 	netModel = "Testnet"
	// }
	// time := utils.GetNowTime()
	// walletBasic := &model.WalletBasic{
	// 	Name:     "ETH",
	// 	Version:  version,
	// 	NetModel: netModel,
	// 	UpdateAt: time,
	// }
	// if err := UpdateWalletBasicInfoToDB(walletBasic); err != nil {
	// 	logrus.Errorf("syncETHWalletInfo UpdateWalletBasicInfoToDB err %s", err)
	// }
	return false
}

// ETH钱包状态信息同步
func (eth *EthSync) syncETHMiningInfo() bool {
	// res, err := eth.getmininginfo()
	// if err != nil {
	// 	logrus.Errorf("getmininginfo error: %s", err.Error())
	// 	return false
	// }
	// balance, err := eth.getbalance()
	// if err != nil {
	// 	logrus.Errorf("getbalance error: %s", err.Error())
	// }
	// difficulty := util.ParseFloatToStrWithPoint(res.Difficulty, 2)
	// networkhash := util.ParseFloatToStrWithPoint(res.NetworkHashps, 0)
	// times := utils.GetNowTime()
	// walletSimple := &model.WalletSimple{
	// 	Name:        "ETH",
	// 	Balance:     balance,
	// 	BlockHeight: res.Blocks,
	// 	Difficulty:  difficulty,
	// 	NetworkHash: networkhash,
	// 	UpdateAt:    times,
	// }
	// if err := AddWalletSimpleInfoToDB(walletSimple); err != nil {
	// 	logrus.Errorf("syncETHMiningInfo AddWalletSimpleInfoToDB err %s", err)
	// }
	return false
}
