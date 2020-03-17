package wallet

import (
	"encoding/json"
	"errors"

	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/wallet/btcrpc"
	"github.com/zhxx123/gomonitor/utils"

	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	UPDATE_INTERVAL = 5
	TX_INDEX_OFFSET = 2
)

type MasterNode struct {
	rpcHost     string
	rpcUsername string
	rcpPassword string
}

func newMNode(rpchost, rpcusername, rpcpassword string) *MasterNode {
	return &MasterNode{rpcHost: rpchost, rpcUsername: rpcusername, rcpPassword: rpcpassword}
}

type BtcSync struct {
	rpc      *btcrpc.BtcRPC
	mnNode   *MasterNode
	coinAddr map[string]int
	serverIP string // 节点服务器钱包
	coinType string // 货币类型，btc
	lastTxid string // 最后一次更新的 txid
	txIndex  int    // lastTxid 的 index
	dbIndex  int    // 数据库记录的最后数值
	confirm  int    // 最低确认块数
}

// var btcsync *BtcSync

// init to sync wallet
func InitBTCWallet() *BtcSync {
	rpcHost := config.RPCConfig.BTCHost
	rpcUser := config.RPCConfig.BTCUser
	rpcPassword := config.RPCConfig.BTCPassword
	confirm := config.RPCConfig.BTCConfirm
	if confirm <= 0 {
		confirm = 3
	}
	mnNode := newMNode(rpcHost, rpcUser, rpcPassword)
	coinType := "btc"
	coinAddr, err := GetAllUserCoinAddr(coinType)
	if err != nil {
		return nil
	}
	client := NewBtcRpcClient(mnNode)
	btcsync := NewBtcSync(client, mnNode, coinAddr, confirm)
	return btcsync
}

func NewBtcRpcClient(mnnode *MasterNode) *btcrpc.BtcRPC {
	logrus.Debug("New BtcRpcClient:", mnnode)
	client := btcrpc.NewBtcRPC(mnnode.rpcHost, mnnode.rpcUsername, mnnode.rcpPassword, btcrpc.WithBTCDebug(false))
	return client
}

func NewBtcSync(client *btcrpc.BtcRPC, mnNode *MasterNode, coinAddr map[string]int, confirm int) *BtcSync {
	rpc := &BtcSync{
		rpc:      client,
		coinAddr: coinAddr,
		mnNode:   mnNode,
		coinType: "btc",
		lastTxid: "",
		txIndex:  0,
		dbIndex:  0,
		confirm:  confirm,
	}
	return rpc
}
func (btc *BtcSync) setBTCWalletConfirm(confirm int) {
	btc.confirm = confirm
}
func (btc *BtcSync) getBTCWalletConfirm() int {
	return btc.confirm
}
func (btc *BtcSync) setWalletAddress(userId int, coinaddr string) bool {
	if _, ok := btc.coinAddr[coinaddr]; ok != true {
		btc.coinAddr[coinaddr] = userId
		return true
	}
	return false
}
func (btc *BtcSync) syncBTCWallet() bool {
	lasttxid, dbindex, _ := GetLastWalletTxFromDB(btc.coinType)
	logrus.Infof("syncWallet getLastWalletTx last txid: %s, dbindex: %d", lasttxid, dbindex)
	ftxid, index := btc.getLastTxIndex(lasttxid)
	logrus.Infof("syncWallet: last txid: %s  index: %d", ftxid, index)
	ip := strings.Split(btc.mnNode.rpcHost, ":")
	btc.serverIP = ip[1][2:]
	btc.lastTxid = lasttxid
	btc.txIndex = index
	btc.dbIndex = dbindex
	return true
	// go autoUpdateTrans(client, index, dbindex, lasttxid, serverIP)
}
func (btc *BtcSync) listTransactions(count, from int) ([]btcrpc.ListTransactions, error) {
	if btc.rpc != nil {
		return btc.rpc.GetListTransactions(count, from, btc.confirm)
	}
	return nil, errors.New("btc.client listTransactions error")
}
func (btc *BtcSync) getTransaction(txHash string) (*btcrpc.Transaction, error) {
	if btc.rpc != nil {
		return btc.rpc.GetTransaction(txHash)
	}
	return nil, errors.New("btc.client getTransaction error")
}

// wallet info
func (btc *BtcSync) getinfo() (*btcrpc.WalletInfo, error) {
	if btc.rpc != nil {
		return btc.rpc.GetInfo()
	}
	return nil, errors.New("btc.client getinfo error")
}

// wallet mining info
func (btc *BtcSync) getmininginfo() (*btcrpc.MiningInfo, error) {
	if btc.rpc != nil {
		return btc.rpc.GetMiningInfo()
	}
	return nil, errors.New("btc.client getmininginfo error")
}

// wallet balance info
func (btc *BtcSync) getbalance() (string, error) {
	if btc.rpc != nil {
		return btc.rpc.GetBalance()
	}
	return "-1", errors.New("btc.client getbalance error")
}

// getLastTxIndex
func (btc *BtcSync) getLastTxIndex(lastTxid string) (string, int) {
	from, count := 0, 1000
	for {
		listTrans, err := btc.listTransactions(count, from)
		if err != nil {
			logrus.Errorf("listTransactions error: %s", err)
			break
		}
		lenTrans := len(listTrans)

		if lenTrans <= 0 {
			break
		}
		if lastTxid == "" {
			if lenTrans < count {
				lastTxid = listTrans[0].TxID
				txidIndex := from + lenTrans
				return lastTxid, txidIndex
			}
		} else {
			for i, result := range listTrans {
				if lastTxid == result.TxID {
					// find the last sample txid
					if i < lenTrans-1 && listTrans[i+1].TxID == lastTxid {
						continue
					}
					txidIndex := from + lenTrans - i
					return lastTxid, txidIndex
				}
			}
		}
		from = from + count
		time.Sleep(500 * time.Millisecond)
		logrus.Debugf("last txid %s from %d count %d transLength: %d", lastTxid, from, count, lenTrans)
	}
	return "", from
}

// start to update transactions once
func (btc *BtcSync) UpdateTransaction() bool {
	totalCount := btc.dbIndex
	updateInterval := UPDATE_INTERVAL
	txIndexOffset := TX_INDEX_OFFSET
	failedOffset := 0
	from := btc.txIndex
	if from > 100 {
		updateInterval = 100
		txIndexOffset = 10
	}
	for {
		// count: Number of acquisitions at a time
		// from: Number of start index
		// Update 10 more data each time
		count := updateInterval + failedOffset
		if from < updateInterval && failedOffset == 0 {
			count = from + txIndexOffset
		}
		from = from - count + txIndexOffset
		fcomplate := false
		if from <= 0 {
			from = 0
		}
		logrus.Debugf("UpdateTransaction: from %d count %d lasttxid %s", from, count, btc.lastTxid)
		transList, err := btc.listTransactions(count, from)
		if err != nil {
			logrus.Errorf("listTransactions error: %s", err)
			return false
		}
		// listtrans size
		transLength := len(transList)
		isFindTxid := false
		var wtList []*model.BtcTx
		for index, result := range transList {
			if btc.lastTxid == result.TxID {
				isFindTxid = true
				// from = index
				continue
			}
			if isFindTxid || btc.lastTxid == "" {
				if index < transLength-1 && transList[index+1].TxID == result.TxID {
					continue
				}
				totalCount = totalCount + 1
				logrus.Debugf("tx index: %d id: %s", totalCount, result.TxID)
				transacion, err := btc.getTransaction(result.TxID)
				walletTxList, err := btc.getTransactionSql(transacion, btc.serverIP, totalCount)
				if err != nil {
					logrus.Errorf("UpdateTransaction getTransactionSql err %s", err)
					continue
				}
				wtList = append(wtList, walletTxList)
			}
		}
		if isFindTxid && len(wtList) == 0 {
			// 本次同步完成
			fcomplate = true
		}
		if isFindTxid == false && btc.lastTxid != "" {
			failedOffset += updateInterval
		}
		// logrus.Debugf("from: %d  count: %d  failedoffset: %d", from, count, failedOffset)
		// updte last txid
		if transLength > 0 {
			// update db
			if len(wtList) > 0 {
				err := BTCMutilyInsert(wtList)
				if err == nil {
					// 数据库更新成功，同步信息 重新赋值
					btc.txIndex = from
					btc.dbIndex = totalCount
					btc.lastTxid = transList[transLength-1].TxID
					logrus.Infof("UpdateTransaction: mutilyInsert succefully wtlist len: %d  new lasttxid: %s", len(wtList), btc.lastTxid)
				} else {
					logrus.Errorf("UpdateTransaction: mutilyInsert failed! last txid: %s err %s", btc.lastTxid, err)
				}
			}
		}
		if fcomplate {
			//初始化更新间隔 为 100 ,更新完成 更改间隔10
			return false
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
	// defer CloseRPCClient()
}

// generate one sql
func (btc *BtcSync) getTransactionSql(detail *btcrpc.Transaction, serverIP string, realIndex int) (*model.BtcTx, error) {
	// 设置精确度为 0.00000001
	var a utils.Accuracy = func() float64 { return 0.00000001 }

	toamount, txtype, toaddress, fee := "0", "", "", "0"
	coinType := "BTC"
	fee = utils.ParseFloatToStr(detail.Fee)
	if a.Equal(detail.Amount, 0) {
		txtype = "toMe"
	} else {
		isToMany := true
		for _, dt := range detail.Details {
			if a.Equal(detail.Amount, dt.Amount) {
				isToMany = false
				// break
			}
			txtype = dt.Category
			toaddress = dt.Address
			toamount = utils.ParseFloatToStr(dt.Amount)
			//更新用户账户资产
			// 判断交易类型,更新数据库
			//如果是收账
			if txtype == "receive" {
				// 如果是当前数据库地址
				if value, ok := btc.coinAddr[toaddress]; ok == true {
					// userid, "btc", 接收数量
					if status, err := UpdateUserAccounts(value, coinType, detail.TxID, toamount); err != nil {
						logrus.Errorf("getTransactionSql updateUserAccounts err: %s userid: %d cointype: %s addr: %s amount: %s", err.Error(), value, coinType, toaddress, toamount)
						if status != -1 { // 表示数据库出错，直接退出当次同步
							return nil, errors.New("update error")
						}
					}
				} else {
					logrus.Infof("getTransactionSql updateUserAccounts not found address %s userid: %d cointype: %s amount: %s", toaddress, value, coinType, toamount)
				}
			}
		}
		if isToMany {
			txtype = "toMany"
			toamount = utils.ParseFloatToStr(detail.Amount)
		}
	}
	detail.Hex = ""
	jsonDetail, err := json.Marshal(detail)
	if err != nil {
		logrus.Errorf("getTransactionSql json.Marshal err %s", err)
		return nil, err
	}
	walletTx := &model.BtcTx{
		TxIndex:   realIndex,
		TxId:      detail.TxID,
		ToAddress: toaddress,
		TxType:    txtype,
		ToAmount:  toamount,
		Fee:       fee,
		Time:      detail.TimeReceived,
		ServerIP:  serverIP,
		Detail:    string(jsonDetail),
	}
	return walletTx, nil
}

// BTC钱包基本信息同步
func (btc *BtcSync) syncBTCWalletInfo() bool {
	res, err := btc.getinfo()
	if err != nil {
		logrus.Errorf("getinfo error %s", err.Error())
		return false
	}
	version := utils.ParseWalletVersion(res.Version)

	netModel := "Main"
	if res.Testnet == true {
		netModel = "Testnet"
	}
	time := utils.GetNowTime()
	walletBasic := &model.WalletBasic{
		Name:     "BTC",
		Version:  version,
		NetModel: netModel,
		UpdateAt: time,
	}
	if err := UpdateWalletBasicInfoToDB(walletBasic); err != nil {
		logrus.Errorf("syncBTCWalletInfo UpdateWalletBasicInfoToDB err %s", err)
	}
	return false
}

// BTC钱包状态信息同步
func (btc *BtcSync) syncBTCMiningInfo() bool {
	res, err := btc.getmininginfo()
	if err != nil {
		logrus.Errorf("getmininginfo error: %s", err.Error())
		return false
	}
	balance, err := btc.getbalance()
	if err != nil {
		logrus.Errorf("getbalance error: %s", err.Error())
	}
	difficulty := utils.ParseFloatToStrWithPoint(res.Difficulty, 2)
	networkhash := utils.ParseFloatToStrWithPoint(res.NetworkHashps, 0)
	time := utils.GetNowTime()
	walletSimple := &model.WalletSimple{
		Name:        "BTC",
		Balance:     balance,
		BlockHeight: res.Blocks,
		Difficulty:  difficulty,
		NetworkHash: networkhash,
		UpdateAt:    time,
	}
	if err := AddWalletSimpleInfoToDB(walletSimple); err != nil {
		logrus.Errorf("syncBTCMiningInfo AddWalletSimpleInfoToDB err %s", err)
	}
	return false
}
