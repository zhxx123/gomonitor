package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//生成随机字符串
func GetRandomStrWithn(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()+[]{}/<>;:=.,?"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

//RandomStr 获取一个随机字符串
func GetRandomStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// RandNum 生成随机数字字符
func GetRandNum(maxnum int) string {
	rand.Seed(time.Now().UnixNano())
	ikind := rand.Intn(maxnum)
	if maxnum == 999 {
		return fmt.Sprintf("%03d", ikind)
	}
	return fmt.Sprintf("%04d", ikind)
}
func GetRandSixNum() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func GetRandNumInt(maxnum int) int {
	rand.Seed(time.Now().UnixNano())
	ikind := rand.Intn(maxnum)
	return ikind
}

/**
 * base64 解码
 * @method func
 */
func Base64Decode(str string) string {
	s, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(s)
}

// 解析版本号
func ParseWalletVersion(nVersion int) string {
	res := ""
	if nVersion%100 == 0 {
		res = fmt.Sprintf("%d.%d.%d", nVersion/1000000, (nVersion/10000)%100, (nVersion/100)%100)
	} else {
		res = fmt.Sprintf("%d.%d.%d.%d", nVersion/1000000, (nVersion/10000)%100, (nVersion/100)%100, nVersion%100)
	}
	return res
}
func StrShelter(data, key string) (string, string) {
	if data == "" {
		return "", data
	}
	res := strings.Split(data, key)
	if len(res) == 2 {
		tmpData := res[0]
		l := len(tmpData)
		if l >= 3 {
			return fmt.Sprintf("%s**", tmpData[0:2]), res[1]
		}
		return fmt.Sprintf("%s*", tmpData[0:1]), res[1]
	}
	return "", data
}

// 字符串脱敏处理
func StrWithShelter(data string, types int) string {
	if data == "" {
		return data
	}
	// 如果是邮箱
	switch types {
	case 1:
		pre_email, res := StrShelter(data, "@")
		suf_email, suffix := StrShelter(res, ".")
		if pre_email != "" && res != "" && suffix != "" {
			return fmt.Sprintf("%s@%s.%s", pre_email, suf_email, suffix)
		}
		return data
	case 2:
		l := len(data)
		if l > 4 {
			return fmt.Sprintf("%s****", data[0:l-4])
		}
		return data
	}
	return data
}
func GetMachineType(userAgent string) string {
	ua := Parse(userAgent)
	if ua.Mobile {
		return "移动 Web"
	}
	return "PC Web"
}
