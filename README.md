### golang web 监控

基于 golang iris web 框架构建,构建 web api 监控系统，客户端使用 Statsd(go 语言版)

监控系统基于 Grafana、Telegraf、InfluxDB

相关文章地址

#### 客户端接口模拟请求代码
```
package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func GetRandNumInt(maxnum int) int {
	rand.Seed(time.Now().UnixNano())
	ikind := rand.Intn(maxnum)
	return ikind
}
func GetMethod() {
	api := []string{"/article/list", "/article/info"}
	//生成要访问的url
	url := "http://127.0.0.1:7000/api"

	for {
		time.Sleep(500 * time.Millisecond)

		randomIndex := GetRandNumInt(2)
		reqPath := fmt.Sprintf("%s%s", url, api[randomIndex])
		resp, err := http.Get(reqPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(body, err)
		}
		fmt.Printf("%s %d\n", reqPath, resp.StatusCode)

	}
}
func PostMethod() {
	api := []string{"/login", "/logout"}
	//生成要访问的url
	url := "http://127.0.0.1:7000/api"

	for {
		time.Sleep(500 * time.Millisecond)

		randomIndex := GetRandNumInt(2)
		reqPath := fmt.Sprintf("%s%s", url, api[randomIndex])
		resp, err := http.Post(reqPath, "", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(body, err)
		}
		fmt.Printf("%s %d\n", reqPath, resp.StatusCode)

	}
}
func main() {
	go GetMethod()
	PostMethod()
}

```