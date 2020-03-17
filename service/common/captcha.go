package common

import (
	"errors"
	"time"

	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/cache"
	"github.com/zhxx123/gomonitor/service/mail"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

var store = base64Captcha.NewMemoryStore(base64Captcha.GCLimitNumber, 5*time.Minute)

// 获取图片验证码
func CreateCaptchaCode(verfiyCode *model.CaptchaCodeJson) (*model.CaptchaCodeRes, string, int) {
	var id, data string
	switch verfiyCode.Type {
	case 1:
		id, data = GenerateCaptcha(verfiyCode.UID, "digtial")
		break
	// case 2:
	// 	data = GenerateCaptcha(verfiyCode.UID, "audio")
	// 	break
	// case 3:
	// 	data = GenerateCaptcha(verfiyCode.UID, "string")
	default:
		return nil, "error type", model.STATUS_FAILED
	}
	res := &model.CaptchaCodeRes{
		Data: data,
		UID:  id,
	}
	// 将当前验证码加入缓存,并且设置有效期
	// OC.Set(res.UID, code, CacheDefaultExpiration)
	return res, "success", model.STATUS_SUCCESS
}

// 生成图片验证码
// base64Captcha create http handler
func GenerateCaptcha(captchaIds, captchaType string) (string, string) {

	var driver base64Captcha.Driver
	switch captchaType {
	case "audio":
		driver = &base64Captcha.DriverAudio{
			Length:   6,
			Language: "zh",
		}
		break
	case "string":
		driver = &base64Captcha.DriverString{
			Height:          35,
			Width:           85,
			NoiseCount:      20,
			ShowLineOptions: 2,
			Length:          4,
		}
		break
	case "chinese":
		driver = &base64Captcha.DriverChinese{
			Height:          35,
			Width:           85,
			NoiseCount:      20,
			ShowLineOptions: 2,
			Length:          4,
		}
		break
	case "math":
		driver = &base64Captcha.DriverMath{
			Height:          35,
			Width:           85,
			NoiseCount:      20,
			ShowLineOptions: 2,
		}
		break
	default:
		driver = &base64Captcha.DriverDigit{
			Height:   35,
			Width:    85,
			Length:   4,
			MaxSkew:  0.6,
			DotCount: 60,
		}
	}
	c := base64Captcha.NewCaptcha(driver, store)
	captchaId, base64blob, _ := c.Generate()
	logrus.Infof("GenerateCaptcha inputId %s type %s genId %s", captchaIds, captchaType, captchaId)
	return captchaId, base64blob
}

// 校验验证图片验证码
func CheckCaptchaCodes(code *model.CodeJson) (string, int) {
	if err := CaptchaVerify(code.UID, code.Code); err != nil {
		return "error captcha code", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// base64Captcha verify http handler
func CaptchaVerify(captchaIds, verifyValue string) error {
	if captchaIds == "" {
		return errors.New("error id")
	}
	if verifyValue == "" {
		return errors.New("error code")
	}
	if verifyResult := store.Verify(captchaIds, verifyValue, true); verifyResult != true {
		return errors.New("错误的验证码")
	}
	return nil
}

// 获取并发送邮箱验证码
func GetEmailVerifyCodeAndSend(address string, option int) (string, error) {
	code := utils.GetRandSixNum()
	siteName := config.ServerConfig.SiteName

	var content string
	var title string
	if option == 1 { // 用户注册
		content = "<p>我们收到您在 " + siteName + "注册请求, 验证码: " + code + " (5 分钟内有效).</p>" +
			"<p>如果您没有在 " + siteName + " 填写过注册信息, 说明有人滥用了您的邮箱, 请删除此邮件, 我们对给您造成的打扰感到抱歉.</p>" +
			"<p>(这是一封自动产生的邮件，请勿回复。)</p>"
		title = "账号注册"
	}
	if option == 2 { // 用户找回密码
		content = "<p>您正在找回密码, 验证码: " + code + " (5 分钟内有效).</p>" +
			"<p>感谢你对" + siteName + "的支持，希望你在" + siteName + "的体验有益且愉快。</p>" +
			"<p>(这是一封自动产生的邮件，请勿回复。)</p>"
		title = "找回密码"
	}
	content += "<p><img src=\"https://massgrid.com/assets/img/wechatkefu-qrcode.png\" style=\"height: 42px;\"/></p>"
	//fmt.Println(content)
	mailTo := []string{address}
	go func() {
		mail.SendHtmlMail(mailTo, title, content)
	}()
	return code, nil
}

//获取并发送手机验证码
func GetPhoneVerifyCodeAndSend(address string) (string, error) {
	return "", nil
}

// 校验验证码
func CheckVerifyCode(uid, address, verifyCode string, igVerifyCode, isClear bool) error {
	if igVerifyCode == true || model.CHECK_VERIFY_CODE == false {
		logrus.Info("email or phone verify code not check")
		return nil
	}
	if uid == "" {
		return errors.New("错误验证码")
	}
	if address == "" {
		return errors.New("错误验证码")
	}
	value, found := cache.OC.Get(address)
	if found {
		if isClear {
			cache.OC.Delete(address) // 清除缓存
		}
		res, ok := (value).(model.VerifyCodeRes)
		if ok != true {
			return errors.New("错误的验证码")
		}
		if res.UID == uid && res.VerifyCode == verifyCode {
			return nil
		}
		return errors.New("验证码有误")

	}
	return errors.New("验证码有误")
}
