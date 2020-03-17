package cron

import (
	"github.com/zhxx123/gomonitor/service/order"
	"github.com/sirupsen/logrus"
)

func syncMinerInfoCron() {
	// 获取所有机器详情，然后写入数据库
	res, err := order.GetFarmsInfo()
	if err != nil {
		logrus.Errorf("syncMinerInfoCron %s", err.Error())
		return
	}
	// 查找机器是否存在
	if len(res) == 0 {
		logrus.Debugf("syncMinerInfoCron no machine")
		return
	}
	for _, value := range res {
		order.UpdateFarmsInfoToDB(value)
		logrus.Infof("SyncMinerInfoCron update farmserver %s %s %d", value.FarmID, value.MinerType, value.AvailableCount)
	}
	// 更新数据库,不主动更新数据库

}
