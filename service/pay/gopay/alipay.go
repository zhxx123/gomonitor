package gopay

import (
	"fmt"
	"io/ioutil"
	_ "net/http"
	"os"

	"github.com/zhxx123/gomonitor/service/pay/gopay/alipay"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

var (
	appID    = "2016101000656342"
	sellerID = "2088102178996355" //非必须选项
	// RSA2(SHA256)
	privateKey   = getPrivKey()
	aliPublicKey = getPubKey() //此处为 支付宝的 公钥
)

const (
	keyPath = ""
)

func getPrivKey() string {
	privKey, err := ioutil.ReadFile(keyPath + "./key/rsa2048_private.txt")
	if err != nil {
		logrus.Errorf("get privateKey failed: %s", err.Error())
	}
	return string(privKey[:])
}
func getPubKey() string {
	pubKey, err := ioutil.ReadFile(keyPath + "./key/rsa2048_public.txt")
	if err != nil {
		logrus.Errorf("get publicKey failed: %s", err.Error())
	}
	return string(pubKey[:])
}

var client *alipay.AliWebClient

func InitAliPay() {
	var err error
	client, err = alipay.New(appID, sellerID, aliPublicKey, privateKey, false)
	if err != nil {
		fmt.Printf("init alipay failed: %s", err.Error())
		os.Exit(-1)
	} else {
		fmt.Println("init alipay success")
	}
	// logrus.Info(client)
}

// 支付宝支付 相关接口
// tradePreCreate 支付宝预创建接口订单
func tradeAliPreCreate(charge *ChargePreCreate) (*alipay.TradePreCreateRsp, error) {
	var p = alipay.TradePreCreate{}
	if charge.OutTradeNo == "" {
		p.OutTradeNo = utils.GenOutTradeNo(charge.PayChannel, charge.PayType, charge.TradeType, charge.UserId)
	} else {
		p.OutTradeNo = charge.OutTradeNo
	}
	p.Subject = charge.Subject
	p.TotalAmount = charge.TotalAmount
	goodItem := &alipay.GoodsDetailItem{GoodsId: charge.GoodsId, GoodsName: "CNY充值", Price: "0"}
	p.GoodsDetail = append(p.GoodsDetail, goodItem)
	p.NotifyURL = charge.CallbackURL
	p.TimeOutExpress = charge.QrCodeTimeOutExpress
	p.QrCodeTimeOutExpress = charge.QrCodeTimeOutExpress //交易超时关闭时间, 5分钟之后开始二维码失效,无法付款
	if client == nil {
		InitAliPay()
	}

	return client.TradePreCreate(p)
}

//tradeAliQuery 支付宝查询订单
func tradeAliQuery(charge *ChargeQuery) (*alipay.TradeQueryRsp, error) {
	var p = alipay.TradeQuery{}
	p.OutTradeNo = charge.OutTradeNo
	return client.TradeQuery(p)
}

//tradeAliReFund 支付宝退款订单
func tradeAliRefund(charge *ChargeRefund) (*alipay.TradeRefundRsp, error) {
	var p = alipay.TradeRefund{}
	p.OutTradeNo = charge.OutTradeNo //商户订单号
	p.OutRequestNo = utils.GenOutTradeNo(charge.PayChannel, charge.PayType, charge.TradeType, charge.UserId)
	p.RefundAmount = charge.RefundAmount
	p.RefundReason = charge.RefundReason
	if client == nil {
		InitAliPay()
	}
	return client.TradeRefund(p)
}

//tradeAliCancel 支付宝取消订单
func tradeAliCancel(charge *ChargeCancel) (*alipay.TradeCancelRsp, error) {
	var p = alipay.TradeCancel{}
	p.OutTradeNo = charge.OutTradeNo
	p.TradeNo = charge.TradeNo
	if client == nil {
		InitAliPay()
	}
	return client.TradeCancel(p)
}

//tradeAliPayClose 支付宝关闭订单
func tradeAliPayClose(charge *ChargeClose) (*alipay.TradeCloseRsp, error) {
	var p = alipay.TradeClose{}
	p.OutTradeNo = charge.OutTradeNo
	p.TradeNo = charge.TradeNo
	p.OperatorId = charge.OperatorId
	if client == nil {
		InitAliPay()
	}
	return client.TradeClose(p)
}
