package hook

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

func GithubHook(ctx iris.Context, content string) (string, int) {
	if config.ServerConfig.GithubAutoDeploy != true {
		return "not runing", model.STATUS_SUCCESS
	}
	pushMethod := ctx.GetHeader("x-github-event")
	if pushMethod != "push" {
		logrus.Error("git push-event error")
		return "Unmatch x-github-event", model.STATUS_FAILED
	}
	signature := ctx.GetHeader("X-Hub-Signature")
	secretKey := config.ServerConfig.GithubSecretKey
	if utils.VerifySignature(signature, content, secretKey) != true {
		logrus.Error("git verifysignature failed")
		return "error signature", model.STATUS_FAILED
	}
	// 添加任务,执行脚本
	AddNewTask(content)
	logrus.Debug("add git auto deploy job")

	return "success", model.STATUS_SUCCESS
}

func GithubHookLog(queryData *model.QueryHookJson) (model.MyMap, string, int) {
	var webHook []model.WebHook
	count := 0
	offset, err := db.GetOffset(queryData.Page, queryData.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	nodeweb := "nodeweb"
	if err := db.DB.Model(model.WebHook{}).Where("project_name = ?", nodeweb).Order("ID desc").
		Count(&count).Offset(offset).Limit(queryData.Limit).
		Find(&webHook).Error; err != nil {
		logrus.Errorf("WebHook failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	// count = len(helpArticleCategory)
	response := model.MyMap{
		"data":   webHook,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}
