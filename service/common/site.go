package common

import (
	"github.com/zhxx123/gomonitor/model"
)

func SiteInfo() (map[string]interface{}, string, int) {
	// var userCount int
	// var topicCount int
	// var replyCount int
	// if err := db.DB.Model(&model.User{}).Count(&userCount).Error; err != nil {
	// 	return nil, err.Error(), model.ErrorCode.ERROR
	// }
	// if err := db.DB.Model(&model.Article{}).Count(&topicCount).Error; err != nil {
	// 	return nil, err.Error(), model.ErrorCode.ERROR
	// }
	// if err := db.DB.Model(&model.Comment{}).Count(&replyCount).Error; err != nil {
	// 	return nil, err.Error(), model.ErrorCode.ERROR
	// }

	// var keyvalueconfig model.KeyValueConfig
	// siteConfig := make(map[string]interface{})
	// siteConfig["name"] = ""
	// siteConfig["icp"] = ""
	// siteConfig["title"] = ""
	// siteConfig["description"] = ""
	// siteConfig["keywords"] = ""
	// siteConfig["logoURL"] = "/images/logo.png"
	// siteConfig["bdStatsID"] = ""
	// siteConfig["luosimaoSiteKey"] = ""
	// if err := db.DB.Where("key_name = ?", "site_config").Find(&keyvalueconfig).Error; err != nil {
	// 	logrus.Error(err)
	// } else {
	// 	if err := json.Unmarshal([]byte(keyvalueconfig.Value), &siteConfig); err != nil {
	// 		logrus.Error(err)
	// 	}
	// }

	// var baiduAdKeyValue model.KeyValueConfig
	// baiduAdConfig := make(map[string]interface{})
	// baiduAdConfig["banner760x90"] = ""
	// baiduAdConfig["banner2_760x90"] = ""
	// baiduAdConfig["banner3_760x90"] = ""
	// baiduAdConfig["ad250x250"] = ""
	// baiduAdConfig["ad120x90"] = ""
	// baiduAdConfig["ad20_3"] = ""
	// baiduAdConfig["ad20_3A"] = ""
	// baiduAdConfig["allowBaiduAd"] = false

	// if err := db.DB.Where("key_name = ?", "baidu_ad_config").Find(&baiduAdKeyValue).Error; err != nil {
	// 	logrus.Error(err)
	// } else {
	// 	if err := json.Unmarshal([]byte(baiduAdKeyValue.Value), &baiduAdConfig); err != nil {
	// 		logrus.Error(err)
	// 	}
	// }

	// res := map[string]interface{}{
	// 	"siteConfig":    siteConfig,
	// 	"baiduAdConfig": baiduAdConfig,
	// 	"userCount":     userCount,
	// 	"topicCount":    topicCount,
	// 	"replyCount":    replyCount,
	// }
	return nil, "success", model.ErrorCode.SUCCESS
}
