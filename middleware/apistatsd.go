package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kataras/iris/v12/context"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/service/stats"
	"github.com/sirupsen/logrus"
)

// 获取 请求路径
func getReqPath(ctx context.Context) string {
	pathArr := strings.Split(ctx.Request().URL.Path, "/")
	for i := len(pathArr) - 1; i >= 0; i-- {
		if pathArr[i] == "" {
			pathArr = append(pathArr[:i], pathArr[i+1:]...)
		}
	}
	for i, path := range pathArr {
		if matched, err := regexp.MatchString("^[0-9]+$", path); matched && err == nil {
			pathArr[i] = "id"
		}
	}
	pathArr = append([]string{strings.ToLower(ctx.Request().Method)}, pathArr...)
	return strings.Join(pathArr, "_")
}

// APIStatsD 统计api请求
func APIStatsD() context.Handler {
	return func(ctx context.Context) {

		t := time.Now()
		ctx.Next()

		if config.StatsDConfig.URL == "" {
			return
		}
		if config.StatsDConfig.Enabled != true {
			return
		}
		duration := time.Since(t)
		durationMS := int64(duration / 1e6) // 转成毫秒

		reqPath := getReqPath(ctx)
		reqPathAndMethod := fmt.Sprintf("api,req_path=%s", reqPath)
		if err := (*stats.StatsdClient).Timing(reqPathAndMethod, durationMS, 1); err != nil {
			logrus.Error(err)
		}
		// logrus.Debugf("reqPath: %s %d\n", reqPath, durationMS)
		status := ctx.ResponseWriter().StatusCode()
		if status != http.StatusGatewayTimeout && durationMS > 5000 {
			timeoutReqPath := fmt.Sprintf("timeout,req_path=%s", reqPath)
			if err := (*stats.StatsdClient).Inc(timeoutReqPath, 1, 1); err != nil {
				logrus.Error(err)
			}
		}
	}
}
