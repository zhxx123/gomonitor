package mail

import (
	"github.com/zhxx123/gomonitor/config"
	"gopkg.in/gomail.v2"
)

func SendGoHtmlMail(mailTo []string, subject, body string) error {
	host := config.EmailConfig.MailHost
	port := config.EmailConfig.MailPort
	user := config.EmailConfig.MailUser
	password := config.EmailConfig.MailPass
	emailFrom := config.EmailConfig.MailFrom

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)  //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Bcc", user)        //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)
	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
}

func SendGoTextMail(mailTo []string, subject, body string) error {
	host := config.EmailConfig.MailHost
	port := config.EmailConfig.MailPort
	user := config.EmailConfig.MailUser
	password := config.EmailConfig.MailPass
	emailFrom := config.EmailConfig.MailFrom

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)  //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Bcc", user)        //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(host, port, user, password)
	return d.DialAndSend(m)
}
