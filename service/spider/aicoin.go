package spider

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zhxx123/gomonitor/model"
	"github.com/sirupsen/logrus"
)

type AICoinSpider struct {
	Url    string
	Name   string
	client *http.Client
}

func NewAICoinSpider(url, name string) *AICoinSpider {
	aiCoinSpider := &AICoinSpider{
		Url:    url,
		Name:   name,
		client: &http.Client{Timeout: time.Second * 10},
	}
	return aiCoinSpider
}
func (this *AICoinSpider) AICoin() (*model.CoinMarket, error) {
	// Request the HTML page.
	newUrl := fmt.Sprintf("%s%s", this.Url, this.Name)

	logrus.Debug("aicoin url:", newUrl)
	reqest, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		logrus.Errorf("AICoin http.NewRequest err: %s", err.Error())
		return nil, err
	}
	//增加header选项
	reqest.Header.Set("User-Agent", "PostmanRuntime/7.18.0")
	reqest.Header.Set("Host", "www.aicoin.cn")
	reqest.Header.Set("Accept", "text/html")
	reqest.Header.Set("Accept-Encoding", "gzip, deflate")
	reqest.Header.Set("Connection", "keep-alive")

	//处理返回结果
	resp, err := this.client.Do(reqest)
	if err != nil {
		logrus.Errorf("AICoinOJ client.Do error: %s", err.Error())
		return nil, err
	}
	// resp, err := this.client.Get(newUrl)

	var res *model.CoinMarket
	// if err != nil {
	// 	logrus.Errorf("AICoinOJ http.Get error: %s", err)
	// 	return res, err
	// }
	defer resp.Body.Close()
	gzipbody, err := gzip.NewReader(resp.Body)
	if err != nil {
		logrus.Errorf("AICoinOJ gzip.NewReader error: %s", err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(gzipbody)
	if err != nil {
		logrus.Errorf("AICoinOJ goquery.NewDocumentFromReader error: %s", err)
		return res, err
	}
	// Find the review items
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		res = this.Process(s)
	})
	return res, nil
}
func (this *AICoinSpider) Process(table *goquery.Selection) *model.CoinMarket {
	var result *model.CoinMarket
	table.Find("tr").Each(func(indexOfTr int, tr *goquery.Selection) {
		// fmt.Printf("Review %d\n %s\n", indexOfTr, tr.Text())
		if indexOfTr == 1 { // 不包含表头
			if trRes := this.ProcessTr(tr); trRes != nil {
				result = trRes
			}
		}
	})
	return result
}
func (this *AICoinSpider) ProcessTr(tr *goquery.Selection) *model.CoinMarket {
	var tdRes []string
	tr.Find("td").Each(func(_ int, td *goquery.Selection) {
		line := td.Text()
		tdRes = append(tdRes, line)
	})
	if len(tdRes) >= 3 && strings.Contains(tdRes[1], this.Name) {
		hduojRes := &model.CoinMarket{
			Name:  this.Name,
			Url:   "www.aicoin.cn",
			Price: strings.Replace(tdRes[3], ",", "", -1),
		}
		return hduojRes
	}
	return nil
}
