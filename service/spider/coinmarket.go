package spider

import (
	"compress/gzip"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zhxx123/gomonitor/model"
	"github.com/sirupsen/logrus"
)

type CoinMarketSpider struct {
	Url    string
	Name   string
	client *http.Client
}

func NewCoinMarketSpider(url, name string) *CoinMarketSpider {
	CoinMarketSpider := &CoinMarketSpider{
		Url:    url,
		Name:   name,
		client: &http.Client{Timeout: time.Second * 15},
	}
	return CoinMarketSpider
}
func (this *CoinMarketSpider) CoinMarket() (*model.CoinMarket, error) {
	// Request the HTML page.
	// newUrl := fmt.Sprintf("%s%s", this.Url, this.Name)
	newUrl := this.Url
	logrus.Debug("coinmarket url: ", newUrl)
	reqest, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		logrus.Errorf("CoinMarket http.NewRequest err: %s", err.Error())
		return nil, err
	}
	//增加header选项
	reqest.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36")
	reqest.Header.Set("host", "coinmarketcap.com")
	reqest.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	reqest.Header.Set("accept-encoding", "gzip, deflate")
	reqest.Header.Set("accept-language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("origin", "https://coinmarketcap.com")
	reqest.Header.Set("referer", "https://coinmarketcap.com/currencies/massgrid/")
	//处理返回结果
	resp, err := this.client.Do(reqest)
	if err != nil {
		logrus.Errorf("CoinMarketOJ client.Do error: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	gzipbody, err := gzip.NewReader(resp.Body)

	if err != nil {
		logrus.Errorf("CoinMarketOJ gzip.NewReader error: %s", err)
		return nil, err
	}

	// body, err := ioutil.ReadAll(gzipbody)
	// if err == nil {
	// 	fmt.Println(string(body))
	// }

	var res *model.CoinMarket
	doc, err := goquery.NewDocumentFromReader(gzipbody)
	if err != nil {
		logrus.Errorf("CoinMarketOJ goquery.NewDocumentFromReader error: %s", err)
		return res, err
	}
	doc.Find(".cmc-details-panel-about__table").Each(func(i int, s *goquery.Selection) {
		// fmt.Printf("class Review %d\n %s\n", i, s.Text())
		if i == 0 {
			res = this.Process(s)
		}
	})
	return res, nil
}
func (this *CoinMarketSpider) Process(table *goquery.Selection) *model.CoinMarket {
	var trResult string
	table.Find("div").Each(func(indexOfTr int, tr *goquery.Selection) {
		// fmt.Printf("tr Review %d\n %s\n", indexOfTr, tr.Text())
		if indexOfTr == 2 { // 不包含表头
			str := tr.Text()
			str = strings.Replace(str, " ", "", -1)
			// str = strings.Replace(str, "USD", "", -1)
			str = strings.Replace(str, "$", "", -1)
			// 去除换行符
			trResult = strings.Replace(str, "\n", "", -1)
		}
	})
	if len(trResult) > 3 {
		hduojRes := &model.CoinMarket{
			Name:  "MGD",
			Url:   "coinmarketcap.com",
			Price: trResult, //strings.Replace(trResult[0], ",", "", -1),
		}
		return hduojRes
	}
	return nil
}
