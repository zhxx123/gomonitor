package alipay

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/zhxx123/gomonitor/utils"
)

var (
	kSuccess = []byte("success")
)

func NewRequest(method, url string, params url.Values) (*http.Request, error) {
	var m = strings.ToUpper(method)
	var body io.Reader
	if m == "GET" || m == "HEAD" {
		if len(params) > 0 {
			if strings.Contains(url, "?") {
				url = url + "&" + params.Encode()
			} else {
				url = url + "?" + params.Encode()
			}
		}
	} else {
		body = strings.NewReader(params.Encode())
	}
	return http.NewRequest(m, url, body)
}

func (this *AliWebClient) NotifyVerify(partnerId, notifyId string) bool {
	var param = url.Values{}
	param.Add("service", "notify_verify")
	param.Add("partner", partnerId)
	param.Add("notify_id", notifyId)
	req, err := NewRequest("GET", this.notifyVerifyDomain, param)
	if err != nil {
		return false
	}

	rep, err := this.Client.Do(req)
	if err != nil {
		return false
	}
	defer rep.Body.Close()

	data, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return false
	}
	if string(data) == "true" {
		return true
	}
	return false
}

func (this *AliWebClient) GetTradeNotification(req *http.Request) (*TradeNotification, error) {
	if req == nil {
		return nil, errors.New("request 参数不能为空")
	}
	req.ParseForm()
	var m = make(map[string]string)
	for k, v := range req.Form {
		// k不会有多个值的情况
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
	}
	var aliPayNotify TradeNotification
	err := utils.MapStrToStruct(m, &aliPayNotify)
	if err != nil {
		return nil, err
	}

	ok, err := this.VerifySign(req.Form)
	if ok == false {
		return nil, err
	}
	return &aliPayNotify, err
}

//  AckNotification 确认
func (this *AliWebClient) AckNotification(w http.ResponseWriter) {
	AckNotification(w)
}

// AckNotification 验证成功,用于向支付宝后台系统发送 成功消息
func AckNotification(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(kSuccess)
}
