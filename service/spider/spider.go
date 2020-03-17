package spider

import (
	"fmt"
	"strings"
	"time"

	"github.com/zhxx123/gomonitor/service/task"
	"github.com/zhxx123/gomonitor/service/wallet"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

var coinMarketSpider *CoinMarket
var aicoinSpiderBTC *AICoin
var aicoinSpiderETH *AICoin
var aicoinSpiderUSDT *AICoin

func InitSpider() {
	// aicoin
	btcName := "BTC"
	aicoinSpiderBTC = InitAICoinSpider(btcName)

	ethName := "ETH"
	aicoinSpiderETH = InitAICoinSpider(ethName)

	usdtName := "USDT"
	aicoinSpiderUSDT = InitAICoinSpider(usdtName)

	// coinmarketcap
	mgdName := "massgrid"
	coinMarketSpider = InitCoinMarketCapSpider(mgdName)
	// coinMarketSpider.Update()
	task.AddTimer(time.Minute*5, SpiderBTC)
	logrus.Debug("InitSpider AddTimer SpiderBTC, time duation 5 mins, repated")

	task.AddTimer(time.Minute*7, SpiderETH)
	logrus.Debug("InitSpider AddTimer SpiderETH, time duation 7 mins, repated")

	task.AddTimer(time.Minute*8, SpiderMGD)
	logrus.Debug("InitSpider AddTimer SpiderMGD,  time duation 9 mins, repated")

	task.AddTimer(time.Minute*11, SpiderUSDT)
	logrus.Debug("InitSpider AddTimer SpiderUSDT,  time duation 11 mins, repated")
	fmt.Println("init spider")
}
func SpiderBTC(param map[string]interface{}) bool {
	return aicoinSpiderBTC.Update()
}
func SpiderETH(param map[string]interface{}) bool {
	return aicoinSpiderETH.Update()
}

func SpiderUSDT(param map[string]interface{}) bool {
	return aicoinSpiderUSDT.Update()
}

func SpiderMGD(param map[string]interface{}) bool {
	return coinMarketSpider.Update()
}

type AICoin struct {
	aiCoinSpider *AICoinSpider
}
type CoinMarket struct {
	coinMarketSpider *CoinMarketSpider
}

// init aicoin
func InitAICoinSpider(coinName string) *AICoin {
	hduUrl := "https://www.aicoin.cn/currencies/all/cny/1/desc?lang=en&search="
	hduSpider := NewAICoinSpider(hduUrl, coinName)
	hduoj := &AICoin{
		aiCoinSpider: hduSpider,
	}
	return hduoj
}

// init massgrid coinmarketcap
func InitCoinMarketCapSpider(coinName string) *CoinMarket {
	hduUrl := "https://coinmarketcap.com/zh/currencies/massgrid/ratings/"
	hduSpider := NewCoinMarketSpider(hduUrl, coinName)
	hduoj := &CoinMarket{
		coinMarketSpider: hduSpider,
	}
	return hduoj
}
func (aiCoin *AICoin) Update() bool {
	// 首先从数据库获取最后一条更新记录
	res, err := aiCoin.aiCoinSpider.AICoin()
	if err != nil || res == nil {
		return false
	}
	res.Time = utils.GetNowTime()
	if err := wallet.UpdateCoinMart(res); err != nil {
		return false
	}
	return false
}

func (coinMarket *CoinMarket) Update() bool {
	// 首先从数据库获取最后一条更新记录
	res, err := coinMarket.coinMarketSpider.CoinMarket()
	if err != nil || res == nil {
		return false
	}
	time := utils.GetNowTime()
	res.Time = time
	price := res.Price

	if strings.Contains(price, "USD") {
		price = strings.Replace(price, "USD", "", -1)
		usdtPrice := wallet.GetUsdtPrice()
		res.Price = utils.FloatMutiplyStr(price, usdtPrice)
		// fmt.Println(price, usdtPrice, res.Price)
	}
	// fmt.Printf("%+v\n", res)
	if err := wallet.UpdateCoinMart(res); err != nil {
		return false
	}
	return false
}
