package wallet

import (
	"encoding/json"
	"errors"

	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/wallet/mgdrpc"
	"github.com/zhxx123/gomonitor/utils"

	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	MGD_UPDATE_INTERVAL = 5
	MGD_TX_INDEX_OFFSET = 2
)

type MGDMasterNode struct {
	rpcHost     string
	rpcUsername string
	rcpPassword string
}

func newMGDMNode(rpchost, rpcusername, rpcpassword string) *MGDMasterNode {
	return &MGDMasterNode{rpcHost: rpchost, rpcUsername: rpcusername, rcpPassword: rpcpassword}
}

type MgdSync struct {
	rpc      *mgdrpc.MgdRPC
	mnNode   *MGDMasterNode
	coinAddr map[string]int
	serverIP string // 节点服务器钱包
	coinType string // 货币类型，mgd
	lastTxid string // 最后一次更新的 txid
	txIndex  int    // lastTxid 的 index
	dbIndex  int    // 数据库记录的最后数值
	confirm  int    // 最低确认块数
}

// var mgdsync *MgdSync

// init to sync wallet
func InitMGDWallet() *MgdSync {
	rpcHost := config.RPCConfig.MGDHost
	rpcUser := config.RPCConfig.MGDUser
	rpcPassword := config.RPCConfig.MGDPassword
	confirm := config.RPCConfig.MGDConfirm
	if confirm <= 0 {
		confirm = 3
	}
	mnNode := newMGDMNode(rpcHost, rpcUser, rpcPassword)
	coinType := "mgd"
	coinAddr, err := GetAllUserCoinAddr(coinType)
	if err != nil {
		return nil
	}
	client := NewMgdRpcClient(mnNode)
	mgdsync := NewMgdSync(client, mnNode, coinAddr, confirm)
	return mgdsync
}

func NewMgdRpcClient(mnnode *MGDMasterNode) *mgdrpc.MgdRPC {
	logrus.Debug("New MgdRpcClient:", mnnode)
	client := mgdrpc.NewMgdRPC(mnnode.rpcHost, mnnode.rpcUsername, mnnode.rcpPassword, mgdrpc.WithMGDDebug(false))
	return client
}

func NewMgdSync(client *mgdrpc.MgdRPC, mnNode *MGDMasterNode, coinAddr map[string]int, confirm int) *MgdSync {
	rpc := &MgdSync{
		rpc:      client,
		coinAddr: coinAddr,
		mnNode:   mnNode,
		coinType: "mgd",
		lastTxid: "",
		txIndex:  0,
		dbIndex:  0,
		confirm:  confirm,
	}
	return rpc
}
func (mgd *MgdSync) setWalletConfirm(confirm int) {
	mgd.confirm = confirm
}
func (mgd *MgdSync) getWalletConfirm() int {
	return mgd.confirm
}
func (mgd *MgdSync) setWalletAddress(userId int, coinaddr string) bool {
	if _, ok := mgd.coinAddr[coinaddr]; ok != true {
		mgd.coinAddr[coinaddr] = userId
		return true
	}
	return false
}
func (mgd *MgdSync) syncMGDWallet() bool {
	lasttxid, dbindex, _ := GetLastWalletTxFromDB(mgd.coinType)
	logrus.Infof("syncWallet getLastWalletTx last txid: %s, dbindex: %d", lasttxid, dbindex)
	ftxid, index := mgd.getLastTxIndex(lasttxid)
	logrus.Infof("syncWallet: last txid: %s  index: %d", ftxid, index)
	ip := strings.Split(mgd.mnNode.rpcHost, ":")
	mgd.serverIP = ip[1][2:]
	mgd.lastTxid = lasttxid
	mgd.txIndex = index
	mgd.dbIndex = dbindex
	return true
	// go autoUpdateTrans(client, index, dbindex, lasttxid, serverIP)
}
func (mgd *MgdSync) listTransactions(count, from int) ([]mgdrpc.ListTransactions, error) {
	if mgd.rpc != nil {
		return mgd.rpc.GetListTransactions(count, from, mgd.confirm)
	}
	return nil, errors.New("mgd.client listTransactions error")
}
func (mgd *MgdSync) getTransaction(txHash string) (*mgdrpc.Transaction, error) {
	if mgd.rpc != nil {
		return mgd.rpc.GetTransaction(txHash)
	}
	return nil, errors.New("mgd.client getTransaction error")
}

// wallet info
func (mgd *MgdSync) getinfo() (*mgdrpc.WalletInfo, error) {
	if mgd.rpc != nil {
		return mgd.rpc.GetInfo()
	}
	return nil, errors.New("mgd.client getinfo error")
}

// wallet mining info
func (mgd *MgdSync) getmininginfo() (*mgdrpc.MiningInfo, error) {
	if mgd.rpc != nil {
		return mgd.rpc.GetMiningInfo()
	}
	return nil, errors.New("mgd.client getmininginfo error")
}

// wallet balance info
func (mgd *MgdSync) getbalance() (string, error) {
	if mgd.rpc != nil {
		return mgd.rpc.GetBalance()
	}
	return "-1", errors.New("mgd.client getbalance error")
}

// getLastTxIndex
func (mgd *MgdSync) getLastTxIndex(lastTxid string) (string, int) {
	from, count := 0, 1000
	for {
		listTrans, err := mgd.listTransactions(count, from)
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
func (mgd *MgdSync) UpdateTransaction() bool {
	totalCount := mgd.dbIndex
	updateInterval := MGD_UPDATE_INTERVAL
	txIndexOffset := MGD_TX_INDEX_OFFSET
	failedOffset := 0
	from := mgd.txIndex
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
		logrus.Debugf("UpdateTransaction: from %d count %d lasttxid %s", from, count, mgd.lastTxid)
		transList, err := mgd.listTransactions(count, from)
		if err != nil {
			logrus.Errorf("listTransactions error: %s", err)
			return false
		}
		// listtrans size
		transLength := len(transList)
		isFindTxid := false
		var wtList []*model.MgdTx
		for index, result := range transList {
			if mgd.lastTxid == result.TxID {
				isFindTxid = true
				// from = index
				continue
			}
			if isFindTxid || mgd.lastTxid == "" {
				if index < transLength-1 && transList[index+1].TxID == result.TxID {
					continue
				}
				totalCount = totalCount + 1
				logrus.Debugf("tx index: %d id: %s", totalCount, result.TxID)
				transacion, err := mgd.getTransaction(result.TxID)
				walletTxList, err := mgd.getTransactionSql(transacion, mgd.serverIP, totalCount)
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
		if isFindTxid == false && mgd.lastTxid != "" {
			failedOffset += updateInterval
		}
		// logrus.Debugf("from: %d  count: %d  failedoffset: %d", from, count, failedOffset)
		// updte last txid
		if transLength > 0 {
			// update db
			if len(wtList) > 0 {
				err := MGDMutilyInsert(wtList)
				if err == nil {
					// 数据库更新成功，同步信息 重新赋值
					mgd.txIndex = from
					mgd.dbIndex = totalCount
					mgd.lastTxid = transList[transLength-1].TxID
					logrus.Infof("UpdateTransaction: mutilyInsert succefully wtlist len: %d  new lasttxid: %s", len(wtList), mgd.lastTxid)
				} else {
					logrus.Errorf("UpdateTransaction: mutilyInsert failed! last txid: %s err %s", mgd.lastTxid, err)
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
func (mgd *MgdSync) getTransactionSql(detail *mgdrpc.Transaction, serverIP string, realIndex int) (*model.MgdTx, error) {
	// 设置精确度为 0.00000001
	var a utils.Accuracy = func() float64 { return 0.00000001 }

	toamount, txtype, toaddress, fee := "0", "", "", "0"
	coinType := "MGD"
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
				if value, ok := mgd.coinAddr[toaddress]; ok == true {
					// userid, "mgd", 接收数量
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
	walletTx := &model.MgdTx{
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

// MGD钱包基本信息同步
func (mgd *MgdSync) syncMGDWalletInfo() bool {
	res, err := mgd.getinfo()
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
		Name:     "MGD",
		Version:  version,
		NetModel: netModel,
		UpdateAt: time,
	}
	if err := UpdateWalletBasicInfoToDB(walletBasic); err != nil {
		logrus.Errorf("syncMGDWalletInfo UpdateWalletBasicInfoToDB err %s", err)
	}
	return false
}

// MGD钱包状态信息同步
func (mgd *MgdSync) syncMGDMiningInfo() bool {
	res, err := mgd.getmininginfo()
	if err != nil {
		logrus.Errorf("getmininginfo error: %s", err.Error())
		return false
	}
	balance, err := mgd.getbalance()
	if err != nil {
		logrus.Errorf("getbalance error: %s", err.Error())
	}
	difficulty := utils.ParseFloatToStrWithPoint(res.Difficulty, 2)
	networkhash := utils.ParseFloatToStrWithPoint(res.NetworkHashps, 0)
	time := utils.GetNowTime()
	walletSimple := &model.WalletSimple{
		Name:        "MGD",
		Balance:     balance,
		BlockHeight: res.Blocks,
		Difficulty:  difficulty,
		NetworkHash: networkhash,
		UpdateAt:    time,
	}
	if err := AddWalletSimpleInfoToDB(walletSimple); err != nil {
		logrus.Errorf("syncMGDMiningInfo AddWalletSimpleInfoToDB err %s", err)
	}
	return false
}
