package utils

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
)

// RelativeURLToAbsoluteURL 相对URL转绝对URL
func RelativeURLToAbsoluteURL(curURL string, baseURL string) (string, error) {
	curURLData, err := url.Parse(curURL)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	baseURLData, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	curURLData = baseURLData.ResolveReference(curURLData)
	return curURLData.String(), nil
}

/**
 * SplitHostPort 分割 host:port
 */
func SplitHostPort(hostPort string) (string, string) {
	hostPart, portPart, err := net.SplitHostPort(hostPort)
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			return hostPort, ""
		}
		return "", ""
	}

	return hostPart, portPart
}
func RemoteAddr(ctx iris.Context) string {
	reqHeader := [...]string{"X-Forwarded-For", "X-Real-IP", "Host"}
	for _, value := range reqHeader {
		if addr := ctx.GetHeader(value); addr != "" {
			return addr
		}
	}
	return ctx.RemoteAddr()
}

// inet_ntoa ip地址 int to string
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// inet_aton ip地址 string to int
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

// isPublicIP 是否公网ip
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
