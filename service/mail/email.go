package mail

import "github.com/sirupsen/logrus"

func SendHtmlMail(to []string, subject, body string) {
	err := SendHtmlEmail(to, subject, body)
	if err == nil {
		logrus.Debug("email SendHtmlMail success st")
		return
	}
	logrus.Errorf("email SendHtmlEmail failed %s\n", err.Error())
	err = SendGoHtmlMail(to, subject, body)
	if err != nil {
		logrus.Errorf("email SendGoHtmlMail failed %s\n", err.Error())
		return
	}
	logrus.Debug("email SendHtmlMail success ed")
	return
}

func SendTextMail(to []string, subject, body string) {
	err := SendTextEmail(to, subject, body)
	if err == nil {
		logrus.Debug("email SendTextMail success st")
		return
	}
	logrus.Errorf("email SendTextEmail failed %s\n", err.Error())
	err = SendGoTextMail(to, subject, body)
	if err != nil {
		logrus.Errorf("email SendGoTextMail failed %s\n", err.Error())
		return
	}
	logrus.Debug("email SendTextMail success ed")
	return
}
