package alipay

// TradeClose https://docs.open.alipay.com/api_1/alipay.trade.close/
func (this *AliWebClient) TradeClose(param TradeClose) (results *TradeCloseRsp, err error) {
	err = this.doRequest("POST", param, &results)
	return results, err
}

// TradeCancel https://docs.open.alipay.com/api_1/alipay.trade.cancel/
func (this *AliWebClient) TradeCancel(param TradeCancel) (results *TradeCancelRsp, err error) {
	err = this.doRequest("POST", param, &results)
	return results, err
}

// TradeRefund https://docs.open.alipay.com/api_1/alipay.trade.refund/
func (this *AliWebClient) TradeRefund(param TradeRefund) (results *TradeRefundRsp, err error) {
	err = this.doRequest("POST", param, &results)
	return results, err
}

// TradePreCreate https://docs.open.alipay.com/api_1/alipay.trade.precreate/
func (this *AliWebClient) TradePreCreate(param TradePreCreate) (results *TradePreCreateRsp, err error) {
	err = this.doRequest("POST", param, &results)
	return results, err
}

// TradeQuery https://docs.open.alipay.com/api_1/alipay.trade.query/
func (this *AliWebClient) TradeQuery(param TradeQuery) (results *TradeQueryRsp, err error) {
	err = this.doRequest("POST", param, &results)
	return results, err
}
