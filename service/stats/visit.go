package stats

import (
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/zhxx123/gomonitor/service/db"

	"github.com/globalsign/mgo/bson"
	"github.com/zhxx123/gomonitor/model"
)

// PV 增加一次页面访问
func PV(ctx iris.Context, clientInfo *model.ClientInfo) (string, int) {
	var err error
	var userVisit model.UserVisit
	userVisit.ID = bson.NewObjectId()
	userVisit.Platform = clientInfo.Platform
	userVisit.ClientID = clientInfo.ClientID
	userVisit.OSName = clientInfo.OSName
	userVisit.OSVersion = clientInfo.OSVersion
	userVisit.Language = clientInfo.Language
	userVisit.Country = clientInfo.Country
	userVisit.DeviceModel = clientInfo.DeviceModel
	userVisit.DeviceWidth, err = strconv.Atoi(clientInfo.DeviceWidth)
	if err != nil {
		return "无效的deviceWidth", model.ErrorCode.ERROR
	}
	userVisit.DeviceHeight, err = strconv.Atoi(clientInfo.DeviceHeight)
	if err != nil {
		return "无效的deviceHeight", model.ErrorCode.ERROR
	}
	userVisit.IP = ctx.RemoteAddr()
	userVisit.Date = time.Now()
	userVisit.Referrer = clientInfo.Referrer
	userVisit.URL = clientInfo.URL
	userVisit.BrowserName = clientInfo.BrowserName
	userVisit.BrowserVersion = clientInfo.BrowserVersion
	if userVisit.ClientID == "" {
		return "clientId不能为空", model.ErrorCode.ERROR
	}

	if err := db.MongoDB.C("userVisit").Insert(&userVisit); err != nil {
		logrus.Error(err)
		return "error", model.ErrorCode.ERROR
	}
	return "success", model.ErrorCode.SUCCESS
}
