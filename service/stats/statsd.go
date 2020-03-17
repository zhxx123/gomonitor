package stats

import (
	"fmt"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/zhxx123/gomonitor/model"
)

// StatsdClient statsd 客户端
var StatsdClient *statsd.Statter

func init() {
	if model.StatsdURL == "" {
		return
	}
	client, err := statsd.NewBufferedClient(model.StatsdURL, model.StatsdPreFix, 300*time.Millisecond, 512)
	if err != nil {
		fmt.Println(err.Error())
	}
	StatsdClient = &client
}
