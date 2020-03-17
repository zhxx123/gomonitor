package gopay

import "github.com/zhxx123/gomonitor/service/pay/gopay/alipay"

// transferAliPayPreCreate
// alipay.TradePreCreateRsp ==> ChargePreCreateRsp
func transferAliPayPreCreate(rsp *alipay.TradePreCreateRsp, chargeRsp *ChargePreCreateRsp) error {
	chargeRsp.Code = rsp.Content.Code
	chargeRsp.Msg = rsp.Content.Msg
	chargeRsp.SubCode = rsp.Content.SubCode
	chargeRsp.SubMsg = rsp.Content.SubMsg

	chargeRsp.OutTradeNo = rsp.Content.OutTradeNo
	chargeRsp.QRCode = rsp.Content.QRCode
	return nil
}

// transferAliPayPreCreate
// alipay.TradePreCreateRsp ==> ChargePreCreateRsp
func transferAliPayQuery(rsp *alipay.TradeQueryRsp, chargeRsp *ChargeQueryRsp) error {
	chargeRsp.Code = rsp.Content.Code
	chargeRsp.Msg = rsp.Content.Msg
	chargeRsp.SubCode = rsp.Content.SubCode
	chargeRsp.SubMsg = rsp.Content.SubMsg

	chargeRsp.TradeNo = rsp.Content.TradeNo
	chargeRsp.OutTradeNo = rsp.Content.OutTradeNo
	chargeRsp.BuyerLogonId = rsp.Content.BuyerLogonId
	chargeRsp.BuyerUserId = rsp.Content.BuyerUserId
	chargeRsp.TradeStatus = rsp.Content.TradeStatus
	chargeRsp.TotalAmount = rsp.Content.TotalAmount
	chargeRsp.SendPayDate = rsp.Content.SendPayDate
	return nil
}

// transferAliPayPreCreate
// alipay.TradePreCreateRsp ==> ChargePreCreateRsp
func transferAliPayCancel(rsp *alipay.TradeCancelRsp, chargeRsp *ChargeCancelRsp) error {
	chargeRsp.Code = rsp.Content.Code
	chargeRsp.Msg = rsp.Content.Msg
	chargeRsp.SubCode = rsp.Content.SubCode
	chargeRsp.SubMsg = rsp.Content.SubMsg

	chargeRsp.TradeNo = rsp.Content.TradeNo
	chargeRsp.OutTradeNo = rsp.Content.OutTradeNo
	chargeRsp.RetryFlag = rsp.Content.RetryFlag
	chargeRsp.Action = rsp.Content.Action
	return nil
}

// transferAliPayPreCreate
// alipay.TradePreCreateRsp ==> ChargePreCreateRsp
func transferAliPayRefund(rsp *alipay.TradeRefundRsp, chargeRsp *ChargeRefundRsp) error {
	chargeRsp.Code = rsp.AliPayTradeRefund.Code
	chargeRsp.Msg = rsp.AliPayTradeRefund.Msg
	chargeRsp.SubCode = rsp.AliPayTradeRefund.SubCode
	chargeRsp.SubMsg = rsp.AliPayTradeRefund.SubMsg

	chargeRsp.TradeNo = rsp.AliPayTradeRefund.TradeNo
	chargeRsp.OutTradeNo = rsp.AliPayTradeRefund.OutTradeNo
	chargeRsp.BuyerLogonId = rsp.AliPayTradeRefund.BuyerLogonId
	chargeRsp.BuyerUserId = rsp.AliPayTradeRefund.BuyerUserId
	chargeRsp.FundChange = rsp.AliPayTradeRefund.FundChange
	chargeRsp.RefundFee = rsp.AliPayTradeRefund.RefundFee
	chargeRsp.GmtRefundPay = rsp.AliPayTradeRefund.GmtRefundPay
	chargeRsp.StoreName = rsp.AliPayTradeRefund.StoreName

	return nil
}

// transferAliPayPreCreate
// alipay.TradePreCreateRsp ==> ChargePreCreateRsp
func transferAliPayClose(rsp *alipay.TradeCloseRsp, chargeRsp *ChargeCloseRsp) error {
	chargeRsp.Code = rsp.AliPayTradeClose.Code
	chargeRsp.Msg = rsp.AliPayTradeClose.Msg
	chargeRsp.SubCode = rsp.AliPayTradeClose.SubCode
	chargeRsp.SubMsg = rsp.AliPayTradeClose.SubMsg

	chargeRsp.TradeNo = rsp.AliPayTradeClose.TradeNo
	chargeRsp.OutTradeNo = rsp.AliPayTradeClose.OutTradeNo
	return nil
}
