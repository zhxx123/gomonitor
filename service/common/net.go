package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/zhxx123/gomonitor/config"

	"github.com/zhxx123/gomonitor/model"
)

func TaobaoAPI(ip string) *model.TaobaoIP {
	url := fmt.Sprintf("%s%s", config.IPLocConfig.TBAPI, ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var result model.TaobaoIP
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}

	return &result
}
func TencentAPI(ip string) *model.TencentIP {
	url := fmt.Sprintf("%s%s&key=%s", config.IPLocConfig.TCAPI, ip, config.IPLocConfig.TCAPIKey)
	// fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	// return string(out)
	var result model.TencentIP
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}
	return &result
}
func BaiduAPI(ip string) *model.BaiduIP {
	url := fmt.Sprintf("%s%s", config.IPLocConfig.BDAPI, ip)
	// fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	// return string(out)
	var result model.BaiduIP
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}
	return &result
}
func GetIpCity(ip string) string {
	ipinfo := GetIPInfo(ip)
	if ipinfo == nil {
		return ""
	}
	country := ipinfo.Area
	if country == "" {
		country = ipinfo.Country
	}
	ipaddr := fmt.Sprintf("%s%s", country, ipinfo.City)
	return ipaddr
}
func GetIPInfo(ip string) *model.IPInfo {

	var ipInfo model.IPInfo
	// TaoBao 查询结果
	tbRsp := TaobaoAPI(ip)
	if tbRsp != nil && tbRsp.Code == 0 {
		ipInfo.IP = ip
		ipInfo.Country = tbRsp.Data.Country
		ipInfo.Area = tbRsp.Data.Region
		ipInfo.City = tbRsp.Data.City
		ipInfo.Isp = tbRsp.Data.Isp
		if ipInfo.City == "" { //如果没有城市名，就赋为地区名
			ipInfo.City = ipInfo.Area
		}
		return &ipInfo
	}
	// Tencent 查询到结果
	tcRsp := TencentAPI(ip)
	if tcRsp != nil && tcRsp.Status == 0 {
		ipInfo.IP = ip
		ipInfo.Country = tcRsp.Result.AdInfo.Nation
		ipInfo.Area = tcRsp.Result.AdInfo.Province
		ipInfo.City = tcRsp.Result.AdInfo.City
		ipInfo.Lng = tcRsp.Result.Location.Lng
		ipInfo.Lat = tcRsp.Result.Location.Lat
		if ipInfo.City == "XX" { //如果没有城市名，就赋为地区名
			ipInfo.City = ipInfo.Area
		}
		return &ipInfo
	}

	// bdRsp := BaiduAPI(ip)
	// if bdRsp != nil && bdRsp.Status == "0" && len(bdRsp.Data) > 0 {
	// 	data := bdRsp.Data[0]
	// 	ipInfo.City = data.Location
	// }
	return nil
}

// GetExternalIP 获取公网 ip
func GetExternalIP() (string, error) {
	MYEXTERNAL_IP_API := config.IPLocConfig.MyExternalAPI
	rsp, err := http.Get(MYEXTERNAL_IP_API)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	content, _ := ioutil.ReadAll(rsp.Body)
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(rsp.Body)
	// //s := buf.String()
	return string(content), nil
}

// GetExternalIPFromMy 获取公网 ip
func GetExternalIPFromMy() (string, error) {
	MYZHXX_IP_API := config.IPLocConfig.MyZhxxAPI
	rsp, err := http.Get(MYZHXX_IP_API)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	content, _ := ioutil.ReadAll(rsp.Body)
	return string(content), nil
}

// GetUsingIPFromDNS 获取使用的内网 ip
func GetUsingIPFromDNS() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}
