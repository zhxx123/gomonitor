package hook

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/mail"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

var running = false
var queue []*structTaskQueue

type structTaskQueue struct {
	requestBodyString string
}

// AddNewTask add new task
func AddNewTask(bodyContent string) {
	queue = append(queue, newStructTaskQueue(bodyContent))
	checkoutTaskStatus()
}

func newStructTaskQueue(body string) *structTaskQueue {
	return &structTaskQueue{body}
}

func checkoutTaskStatus() {
	if running {
		return
	}
	if len(queue) > 0 {
		requstBody := queue[0].requestBodyString
		queue = queue[1:]
		go startTask(requstBody)
	}
}
func execShell(cmds string) bool {
	if cmds == "" {
		return false
	}
	cmd := exec.Command("/bin/bash", cmds)
	_, err := cmd.Output()
	if err != nil {
		logrus.Error("github pull start task failed")
		return false
	}
	logrus.Debugf("github pull start task success %s\n", cmds)
	return true
}
func startTask(requstBody string) {
	running = true
	request := make(map[string]interface{})
	err := json.Unmarshal([]byte(requstBody), &request)
	if err == nil {

		nowTime := utils.GetNowTime()
		commits := request["commits"].([]interface{})[0].(map[string]interface{})
		output := commits["message"].(string)
		timestamp := commits["timestamp"].(string)
		author := commits["author"].(map[string]interface{})
		authorName := author["name"].(string)
		authorEmail := author["email"].(string)
		refName := request["ref"].(string)
		refName = strings.Replace(refName, "refs/heads/", "", -1)
		pName := request["repository"].(map[string]interface{})
		projectName := pName["name"].(string)
		afterID := request["after"].(string)
		shell := ""
		if projectName == "nodeadmin" {
			shell = "./shell/deploy_admin.sh"
		} else if projectName == "nodeweb" {
			shell = "./shell/deploy.sh"
		}
		//执行 shell
		taskStatus := false
		if config.WebHookConfig.WebHookShell {
			taskStatus = execShell(shell)
		}
		webhook := &model.WebHook{
			Type:        "github",
			Refer:       refName,
			ProjectName: projectName,
			AuthorName:  authorName,
			AuthorEmail: authorEmail,
			Forced:      request["forced"].(bool),
			BeforeID:    request["before"].(string),
			AfterID:     afterID,
			Status:      taskStatus,
			Output:      output,
			TimeStamp:   timestamp,
			Time:        nowTime,
		}
		if err := db.DB.Create(webhook).Error; err != nil {
			logrus.Error(err)
		}
		if checkSendMail(projectName, taskStatus) {
			taskOutput := "失败"
			if taskStatus == true {
				taskOutput = "成功"
			}
			sendMailAuthor := []string{authorEmail}
			mailAuthor := config.WebHookConfig.MailAuthor
			if utils.VerifyEmail(mailAuthor) {
				sendMailAuthor = append(sendMailAuthor, mailAuthor)
			}
			content := "<h4>" + output + "</h4>" +
				"<p>项目:" + projectName + "/" + refName + "</p>" +
				"<p>最新ID:" + afterID + " 打包状态: " + taskOutput + "</p>" +
				"<p>" + authorName + " 提交于 " + timestamp + "</p>"
			title := "提交状态"
			// 发送邮件
			go func() {
				mail.SendHtmlMail(sendMailAuthor, title, content)
			}()
			logrus.Debug("send github log notify email")
		} else {
			logrus.Debug("not send email")
		}
	}
	running = false
	checkoutTaskStatus()
}
func checkSendMail(projectName string, taskStatus bool) bool {
	sendProject := config.WebHookConfig.MailProject
	if strings.Contains(sendProject, projectName) == false {
		return false
	}
	sendLevel := config.WebHookConfig.MailAutoSend
	if sendLevel == 1 && taskStatus == true {
		return false
	}
	return true
}
