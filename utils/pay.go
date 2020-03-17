package utils

import (
	"fmt"
	"time"
)

const (
	OutTradeNoLength = 23
)

// 生成支付订单 ID GenOutTradeNo
func GenOutTradeNo(paych, paytype, tradetype, userid int) string {
	timestmp := time.Now().Format("01021504")
	irand := GetRandNum(999)
	pay := 1 // 为了和商品订单号区分 3+8+1+1+1+1+5+3
	return fmt.Sprintf("MGD%s%d%d%d%d%05d%s", timestmp, pay, paych, paytype, tradetype, userid, irand)
}

/** 校验订单号
 * -1 完全不匹配
 * 0 用户id不匹配
 * 1 支付渠道不匹配
 * 2 支付类型不匹配
 * 3 完全匹配
 */
func CheckOutTradeNo(paych, paytype, userId int, outTradeNo string) int {
	if len(outTradeNo) != OutTradeNoLength {
		return -1
	}
	tmpUserID := fmt.Sprintf("%05d", userId)
	if outTradeNo[15:20] != tmpUserID {
		return 0
	}
	tmpPaych := fmt.Sprintf("%d", paych)
	if outTradeNo[12:13] != tmpPaych {
		return 1
	}
	tmpPayType := fmt.Sprintf("%d", paytype)
	if outTradeNo[12:13] != tmpPayType {
		return 2
	}
	return 3
}

// 生成内部支付订单 ID
func GenTradeNo(paych, paytype, tradetype, userid int) string {
	timestmp := time.Now().Format("200601021504")
	irand := GetRandNum(9999)
	pay := 1 // 为了和商品订单号区分
	return fmt.Sprintf("%s%d%d%d%d%05d%s", timestmp, pay, paych, paytype, tradetype, userid, irand)
}

// 生成购买订单 ID GenOutTradeNo
func GenOrderOutTradeNo(ordertype, tradetype, userid int) string {
	timestmp := time.Now().Format("01021504")
	irand := GetRandNum(999)
	order := 2 // 为了和支付订单号区分
	// MGD + 时间(8位) + 订单标识(1:充值订单 2:购买订单) + 订单类型(1:云算力 2:矿机) + 交易类型(1:收单 2:退款) + 用户id(5位) + 随机数(3位) = 22
	return fmt.Sprintf("MGD%s%d%d%d%05d%s", timestmp, order, ordertype, tradetype, userid, irand)
}

// 生成内部购买订单号
func GenOrderTradeNo(ordertype, tradetype, userid int) string {
	timestmp := time.Now().Format("200601021504")
	irand := GetRandNum(9999)
	order := 2 // 为了和支付订单号区分
	// 时间(12) + 订单标识(1:充值订单 2:购买订单) + 订单类型(1:云算力 2:矿机) + 交易类型(1:收单 2:退款) + 用户id(5位) + 随机数(4位) = 24
	return fmt.Sprintf("%s%d%d%d%05d%s", timestmp, order, ordertype, tradetype, userid, irand)
}

// 商城商品ID
func GenGoodsID(goodstype int) string {
	timestmp := time.Now().Format("010215")
	irand := GetRandNumInt(999)
	return fmt.Sprintf("1%s%d%03d", timestmp, goodstype, irand)
}

// 生成支付订单
func GenWorkOrderNo(userid int) string {
	timestmp := time.Now().Format("200601021504")
	irand := GetRandNumInt(999)
	// 时间(12) + 用户id(5位) + 随机数(3位) = 20
	return fmt.Sprintf("%s%05d%03d", timestmp, userid, irand)
}
