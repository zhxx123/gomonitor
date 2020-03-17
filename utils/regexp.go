package utils

import (
	"fmt"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

/**
* VerifyEmail 校验邮箱
 */
func VerifyEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	// pattern2 := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	// regexp.MatchString(pattern2, email)
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

/**
* VerifyEmail 校验邮箱
 */
func VerifyUserName(username string) bool {
	pattern := `^\w{4,20}$`
	res, _ := regexp.MatchString(pattern, username)
	return res
}

/**
* VerifyEmail 校验手机号
 */
func VerifyPhone(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}
func VerifyPassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[^A-Za-z0-9_]{1}`
	var res int
	if a, err := regexp.MatchString(num, password); a && err == nil {
		res++
	}

	if b, err := regexp.MatchString(a_z, password); b && err == nil {
		res++
	}
	if c, err := regexp.MatchString(A_Z, password); c && err == nil {
		res++
	}
	if d, err := regexp.MatchString(symbol, password); d && err == nil {
		res++
	}
	fmt.Println(res)
	if res >= 2 {
		return true
	}
	return false
}

// AvoidXSS 避免XSS
func AvoidXSS(theHTML string) string {
	return bluemonday.UGCPolicy().Sanitize(theHTML)
}
