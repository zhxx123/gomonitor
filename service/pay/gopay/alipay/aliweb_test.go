package alipay

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	appID    = "2016101000656342"
	sellerID = "2088102178996355" //非必须选项
	// RSA2(SHA256)
	privateKey   = getPrivKey()
	aliPublicKey = getPubKey() //此处为 支付宝的 公钥
)

const (
	keyPath_test = "/home/zhxx/go/src/github.com/massgrid/"
)

func getPrivKey() string {
	privKey, err := ioutil.ReadFile(keyPath_test + "goserver/service/gopay/key/rsa2048_private.txt")
	if err != nil {
		fmt.Println("privateKey:", err)
	}
	// fmt.Println("pubkey:\n",string(privKey[:]))
	return string(privKey[:])
}
func getPubKey() string {
	pubKey, err := ioutil.ReadFile(keyPath_test + "goserver/service/gopay/key/rsa2048_public.txt")
	if err != nil {
		fmt.Println("publicKey:", err)
	}
	// fmt.Println("pubkey:\n",string(pubKey[:]))
	return string(pubKey[:])
}

var client *AliWebClient

func init() {
	var err error
	client, err = New(appID, sellerID, aliPublicKey, privateKey, false)
	fmt.Println(client)
	if err != nil {
		fmt.Println("初始化支付宝失败, 错误信息为:", err)
		os.Exit(-1)
	} else {
		fmt.Println("初始化成功")
	}

	http.HandleFunc("/alipay", func(rep http.ResponseWriter, req *http.Request) {
		var noti, _ = client.GetTradeNotification(req)
		if noti != nil {
			fmt.Println("支付成功")
		} else {
			fmt.Println("支付失败", noti)
		}
		client.AckNotification(rep) // 确认收到通知消息
	})
	//http.ListenAndServe(":80", nil)
}
