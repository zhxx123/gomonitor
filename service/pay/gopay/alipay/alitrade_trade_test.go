package alipay

import (
	"testing"
)

// 统一收单接口,预创建交易
func TestAliPay_TradePreCreate(t *testing.T) {
	t.Log("========== TradePreCreate ==========")
	var p = TradePreCreate{}
	payCh, payType, tradeType, userId := 1, 1, 1, 3
	p.OutTradeNo = util.GenOutTradeNo(payCh, payType, tradeType, userId)
	p.Subject = "测试订单04"
	p.TotalAmount = "109"
	var goodItem = &GoodsDetailItem{GoodsId: "no_10010", GoodsName: "6卡p102", Price: "1.7"}
	p.GoodsDetail = append(p.GoodsDetail, goodItem)
	p.TimeOutExpress = "5m"
	p.QrCodeTimeOutExpress = "5m" //交易超时关闭时间, 5分钟之后开始二维码失效,无法付款

	rsp, err := client.TradePreCreate(p)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Content.Code != K_SUCCESS_CODE {
		t.Fatal(rsp.Content.Msg, rsp.Content.SubMsg)
	}
	t.Log(util.StructToStr(rsp))
	t.Logf("param: %+v", p)
	t.Logf("OutTradeNo: %s %s", p.OutTradeNo, util.NowTimeFormat())
	t.Log(rsp.Content.QRCode)
}

//统一收单接口, 查询交易
func TestAliPay_TradeQuery(t *testing.T) {
	t.Log("========== TradeQuery ==========")
	var p = TradeQuery{}
	p.OutTradeNo = "MGD12191351111000003576"
	rsp, err := client.TradeQuery(p)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Content.Code != K_SUCCESS_CODE {
		t.Fatal(rsp.Content.Msg, rsp.Content.SubMsg)
	}

	t.Log(util.StructToStr(rsp))
	t.Log(rsp.Content)
}

// 统一收单接口, 取消订单
func TestAliPay_TradeCancel(t *testing.T) {
	t.Log("========== TradeCancel ==========")
	var p = TradeCancel{}
	p.OutTradeNo = "MGD07291404110003851"
	rsp, err := client.TradeCancel(p)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Content.Code != K_SUCCESS_CODE {
		t.Fatal(rsp.Content.Msg, rsp.Content.SubMsg)
	}
	t.Log(util.StructToStr(rsp))
	t.Log(rsp.Content)
}

// 统一收单接口,退款
func TestAliPay_TradeRefund(t *testing.T) {
	t.Log("========== TradeRefund ==========")
	var p = TradeRefund{}
	p.RefundAmount = "0.10"
	p.OutTradeNo = "MGD07291424110003784"
	payCh, payType, tradeType, userId := 1, 1, 1, 3
	p.OutRequestNo = util.GenOutTradeNo(payCh, payType, tradeType, userId)
	rsp, err := client.TradeRefund(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(util.StructToStr(rsp))
	t.Logf("%v", rsp.AliPayTradeRefund)
}

// 统一收单接口, 关闭订单
func TestAliPay_TradeClose(t *testing.T) {
	t.Log("========== TradeClose ==========")
	var p = TradeClose{}
	p.OutTradeNo = "MGD07291424110003784"
	rsp, err := client.TradeClose(p)
	if err != nil {
		t.Fatal(err)
	}
	if rsp.IsSuccess() == false {
		t.Fatal(rsp.AliPayTradeClose.Msg, rsp.AliPayTradeClose.SubMsg)
	}
	t.Log(util.StructToStr(rsp))
	t.Log(rsp.AliPayTradeClose)
}
