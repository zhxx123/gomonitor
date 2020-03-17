package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/zhxx123/gomonitor/config"
	"github.com/sirupsen/logrus"
)

var emailPool *email.Pool
var fromEmailUsername string
var ccEmailUsername string
var emailTemplate = `
<div style="background-color:white;border-top:2px solid #12ADDB;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
	<h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;">
	<span style="color: #12ADDB;font-weight:bold;">
		{{.Title}}
	</span>
	</h2>
	<div style="padding:0 12px 0 12px; margin-top:18px;">
	{{if .Content}}
		<p>
			{{.Content}}
		</p>
	{{end}}
	<div style="background-color: #f5f5f5;padding: 10px 15px;margin:0 auto; text-align:center;word-wrap:break-word;">
		<img src="https://massgrid.com/assets/img/wechatkefu-qrcode.png" alt="公众号二维码">
		<p>欢迎关注公众号</p>
	</div>
	{{if .Url}}
	<p>
		<a style="text-decoration:none; color:#12addb" href="{{.Url}}" target="_blank" rel="noopener">点击访问官网</a>
	</p>
	{{end}}
	</div>
</div>
`

// InitEmail 初始化邮件系统
func InitEmail() {
	host := config.EmailConfig.MailHost
	port := config.EmailConfig.MailPort
	user := config.EmailConfig.MailUser
	password := config.EmailConfig.MailPass
	emailFrom := config.EmailConfig.MailFrom
	addr := fmt.Sprintf("%s:%d", host, port)
	emailPool, _ = email.NewPool(addr, 2, smtp.PlainAuth("", user, password, host))
	fromEmailUsername = emailFrom
	ccEmailUsername = user
	fmt.Printf("init email %s %s\n", addr, user)
}

// SendText 发送文本文件
func SendTextEmail(to []string, subject, content string) error {
	e := email.NewEmail()
	e.From = fromEmailUsername
	e.To = to
	e.Subject = subject
	e.Text = []byte(content)
	logrus.Debugf("send  text email: %s %s %s\n", e.From, e.To, e.Subject)
	return emailPool.Send(e, 10*time.Second)
}
func SendHtmlEmail(to []string, subject, html string) error {
	e := email.NewEmail()
	e.From = fromEmailUsername
	e.To = to
	e.Bcc = []string{ccEmailUsername}
	e.Subject = subject
	e.HTML = []byte(html)
	logrus.Debugf("send html email: %s %s %s\n", e.From, e.To, e.Subject)
	return emailPool.Send(e, 10*time.Second)
}

func SendTemplateEmail(to []string, subject, content, title, quoteContent, url string) error {
	tpl, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, map[string]interface{}{
		"Title":   title,
		"Content": content,
		"Url":     url,
	})
	if err != nil {
		return err
	}

	html := b.String()
	return SendHtmlEmail(to, subject, html)
}
