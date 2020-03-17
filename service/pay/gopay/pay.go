package gopay

import (
	"errors"
	"fmt"
	"sync"

	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

var one sync.Once

func InitPay() {
	one.Do(func() { //只执行一次
		InitAliPay()
		//InitWePay()
	})
}

// DoPay 用户预支付下单支付接口
func DoTradePreCreate(charge *ChargePreCreate) (*ChargePreCreateRsp, error) {
	var err error
	err = CheckCharge(charge)
	if err != nil {
		return nil, err
	}
	switch charge.PayChannel {
	case ALI_PAY:
		return doAliPayPreCreate(charge)
	// case WECHAT_PAY:
	// return doWechatPay(charge)
	default:
		err = errors.New("不支持的支付渠道")
	}
	return nil, err
}

// DoTradeQuery 支付查询接口
func DoTradeQuery(charge *ChargeQuery) (*ChargeQueryRsp, error) {
	var err error

	switch charge.PayChannel {
	case ALI_PAY:
		return doAliPayQuery(charge)
	// case WECHAT_PAY:
	// return doWechatPay(charge)
	default:
		err = errors.New("不支持的支付渠道")
	}
	return nil, err
}

// DoTradeCancel 支付取消接口
func DoTradeCancel(charge *ChargeCancel) (*ChargeCancelRsp, error) {
	var err error

	switch charge.PayChannel {
	case ALI_PAY:
		return doAliPayCancel(charge)
	// case WECHAT_PAY:
	// return doWechatPay(charge)
	default:
		err = errors.New("不支持的取消渠道")
	}
	return nil, err
}

// DoTradeRefund 支付退款接口
func DoTradeRefund(charge *ChargeRefund) (*ChargeRefundRsp, error) {
	var err error

	switch charge.PayChannel {
	case ALI_PAY:
		return doAliPayRefund(charge)
	// case WECHAT_PAY:
	// return doWechatPay(charge)
	default:
		err = errors.New("不支持的退款渠道")
	}
	return nil, err
}

// DoTradeClose 支付退款接口
func DoTradeClose(charge *ChargeClose) (*ChargeCloseRsp, error) {
	var err error

	switch charge.PayChannel {
	case ALI_PAY:
		return doAliPayClose(charge)
	// case WECHAT_PAY:
	// return doWechatPay(charge)
	default:
		err = errors.New("不支持的退款渠道")
	}
	return nil, err
}

// doAliPay 支付宝预支付接口
func doAliPayPreCreate(charge *ChargePreCreate) (*ChargePreCreateRsp, error) {
	var chargeRsp ChargePreCreateRsp
	var err error
	if charge.PayType == ALI_WEB { //支付宝网页支付
		rsp, err := tradeAliPreCreate(charge)
		if err != nil {
			logrus.Errorf("doAliPay client.TradePreCreate err %s", err.Error())
			return nil, err
		}
		if rsp.IsSuccess() == false {
			logrus.Errorf("doAliPay rsp.Content.Code err %s %s %s", rsp.Content.Code, rsp.Content.Msg, rsp.Content.SubMsg)
		}
		// 转换为 chargeRsp
		err = transferAliPayPreCreate(rsp, &chargeRsp)

	} else if charge.PayType == ALI_APP { //支付宝APP 支付

	} else {
		err = errors.New("支付宝,不支持的付款方式")
	}
	return &chargeRsp, err
}

// 用户查询接口
func doAliPayQuery(charge *ChargeQuery) (*ChargeQueryRsp, error) {
	var chargeRsp ChargeQueryRsp
	var err error
	if charge.PayType == ALI_WEB { //支付宝网页支付
		rsp, err := tradeAliQuery(charge)
		if err != nil {
			logrus.Errorf("doAliPayQuery client.tradeAliQuery err %s", err.Error())
			return nil, err
		}
		if rsp.IsSuccess() == false {
			logrus.Errorf("doAliPayQuery rsp.Content.Code err %s %s %s", rsp.Content.Code, rsp.Content.Msg, rsp.Content.SubMsg)
		}
		// 根据错误判断,是否需要重试,比如订单号重复
		err = transferAliPayQuery(rsp, &chargeRsp)
	} else if charge.PayType == ALI_APP { //支付宝APP 支付

	} else {
		err = errors.New("支付宝,不支持的查询方式")
	}
	return &chargeRsp, err
}

// 取消支付接口
func doAliPayCancel(charge *ChargeCancel) (*ChargeCancelRsp, error) {
	var chargeRsp ChargeCancelRsp
	var err error
	if charge.PayType == ALI_WEB { //支付宝网页支付
		rsp, err := tradeAliCancel(charge)
		if err != nil {
			logrus.Errorf("doAliPayCancel client.tradeAliCancel err %s", err)
			return nil, err
		}
		if rsp.IsSuccess() == false {
			logrus.Errorf("doAliPayCancel rsp.Content.Code err %s %s %s", rsp.Content.Code, rsp.Content.Msg, rsp.Content.SubMsg)
		}
		// 根据错误判断,是否需要重试,比如订单号重复
		err = transferAliPayCancel(rsp, &chargeRsp)
	} else if charge.PayType == ALI_APP { //支付宝APP 支付

	} else {
		err = errors.New("支付宝,不支持的取消方式")
	}

	return &chargeRsp, err
}

// 退款接口
func doAliPayRefund(charge *ChargeRefund) (*ChargeRefundRsp, error) {
	var chargeRsp ChargeRefundRsp
	var err error
	if charge.PayType == ALI_WEB { //支付宝网页支付
		rsp, err := tradeAliRefund(charge)
		if err != nil {
			logrus.Errorf("doAliPayRefund client.doAliPayRefund err %s", err)
			return nil, err
		}
		if rsp.IsSuccess() == false {
			logrus.Errorf("doAliPayRefund rsp.Content.Code err %s %s %s", rsp.AliPayTradeRefund.Code, rsp.AliPayTradeRefund.Msg, rsp.AliPayTradeRefund.SubMsg)
		}
		// 根据错误判断,是否需要重试,比如订单号重复
		err = transferAliPayRefund(rsp, &chargeRsp)
	} else if charge.PayType == ALI_APP { //支付宝APP 支付

	} else {
		err = errors.New("支付宝,不支持的取消方式")
	}

	return &chargeRsp, err
}

// 关闭接口
func doAliPayClose(charge *ChargeClose) (*ChargeCloseRsp, error) {
	var chargeRsp ChargeCloseRsp
	var err error
	if charge.PayType == ALI_WEB { //支付宝网页支付
		rsp, err := tradeAliPayClose(charge)
		if err != nil {
			logrus.Errorf("doAliPayClose client.tradeAliCancel err %s", err)
			return nil, err
		}
		if rsp.IsSuccess() == false {
			logrus.Errorf("doAliPayClose rsp.Content.Code err %s %s %s", rsp.AliPayTradeClose.Code, rsp.AliPayTradeClose.Msg, rsp.AliPayTradeClose.SubMsg)
		}
		// 根据错误判断,是否需要重试,比如订单号重复
		err = transferAliPayClose(rsp, &chargeRsp)
	} else if charge.PayType == ALI_APP { //支付宝APP 支付

	} else {
		err = errors.New("支付宝,不支持的类型")
	}

	return &chargeRsp, err
}

// WECHAT PAY 支付
//doWechatPay 微信支付
// func doWechatPay(charge *ChargePreCreate) (map[string]string, error) {
// 	var resPonse map[string]string
// 	var err error
// 	if charge.PayType == WECHAT_WEB { // 微信网页支付

// 	} else if charge.PayType == WECHAT_APP { //微信APP 支付

// 	} else if charge.PayType == WECHAT_APP { //微信小程序 支付

// 	} else {
// 		err = errors.New("微信,不支持的付款方式")
// 	}
// 	return resPonse, err
// }

// 初步验证支付内容
func CheckCharge(charge *ChargePreCreate) error {
	// payMethod 支付宝或者微信支付
	if charge.PayChannel <= 0 {
		errs := fmt.Sprintf("支付渠道错误: %d", charge.PayChannel)
		return errors.New(errs)
	}
	if charge.PayType < 0 {
		errs := fmt.Sprintf("支付类型错误: %d", charge.PayType)
		return errors.New(errs)
	}
	moneyTotal, err := utils.ParseStrToFloat(charge.TotalAmount)
	if err != nil {
		errs := fmt.Sprintf("支付金额,%s", err)
		return errors.New(errs)
	}
	var a utils.Accuracy = func() float64 { return 0.01 }
	if a.Smaller(moneyTotal, 0.01) {
		errs := fmt.Sprintf("总金额不能小于 0.01: 当前金额: %f", moneyTotal)
		return errors.New(errs)
	}
	return nil
}
