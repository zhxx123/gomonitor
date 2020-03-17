package user

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/utils"
)

// 登录
func Login() (string, int) {
	// 随机模拟超时
	if utils.GetRandNumInt(100) > 80 {
		times := utils.GetRandNumInt(10000)
		logrus.Debug("sleep: ", times)
		time.Sleep(time.Duration(times) * time.Millisecond)
	}
	return "success", model.SUCCESS
}

func Logout() (string, int) {
	// 随机模拟超时
	if utils.GetRandNumInt(100) > 80 {
		times := utils.GetRandNumInt(10000)
		logrus.Debug("sleep: ", times)
		time.Sleep(time.Duration(times) * time.Millisecond)
	}
	return "success", model.SUCCESS
}

func ArtilceList() (string, int) {
	// 随机模拟超时
	if utils.GetRandNumInt(100) > 70 {
		times := utils.GetRandNumInt(8000)
		logrus.Debug("sleep: ", times)
		time.Sleep(time.Duration(times) * time.Millisecond)
	}
	return "success", model.SUCCESS
}

func ArticleInfo() (string, int) {
	// 随机模拟超时
	if utils.GetRandNumInt(100) > 90 {
		times := utils.GetRandNumInt(8000)
		logrus.Debug("sleep: ", times)
		time.Sleep(time.Duration(times) * time.Millisecond)
	}
	return "success", model.SUCCESS
}
