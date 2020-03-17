package cron

import (
	"fmt"

	"github.com/zhxx123/gomonitor/model"

	"github.com/zhxx123/gomonitor/config"
	"github.com/robfig/cron"
)

// var cronMap = map[string]func(){}
var c *cron.Cron

// func init() {
// 	// if config.ServerConfig.Env != model.DevelopmentMode {
// 	// 	cronMap["0 0 3 * * *"] = yesterdayCron
// 	// }
// 	// if config.ServerConfig.AutoSyncGoods = true {
// 	// 	cronMap
// 	// }
// }

// New 构造cron
// func New() *cron.Cron {
// 	c := cron.New()
// 	for spec, cmd := range cronMap {
// 		c.AddFunc(spec, cmd)
// 	}
// 	return c
// }

func InitCron() {
	c = cron.New()
	if config.ServerConfig.Env != model.DevelopmentMode {
		c.AddFunc("0 0 3 * * *", yesterdayCron)
	}
	if config.ServerConfig.AutoSyncGoods == true { // 每隔5分钟同步矿场机器状态
		c.AddFunc("0 */5 * * *", syncMinerInfoCron)
	}
	c.Start()
	fmt.Println("init cron")
}
