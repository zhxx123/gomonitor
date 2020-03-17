package log

import (
	"fmt"
	"os"

	"github.com/zhxx123/gomonitor/config"
	"github.com/sirupsen/logrus"
)

// var fd *os.File
// var Logrus = log.New()

func InitLog() {
	// 格式
	output := config.ServerConfig.LogOutput
	logLevel := config.ServerConfig.LogLevel
	// 输出到控制台
	if output == "console" {

		level, ok := logrus.ParseLevel(logLevel)
		if ok == nil {
			logrus.SetLevel(level)
		} else {
			logrus.SetLevel(logrus.WarnLevel)
		}
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05", //时间格式化
			FullTimestamp:   true,
		})
		logrus.SetOutput(os.Stdout)
	} else {
		// 屏蔽默认的 logrus 实例输入
		NewLfsHook(logLevel)
		// logrus.AddHook() // 添加 log hook

	}
	fmt.Println("init Log", output, logLevel)
}

// 日志文件 退出函数
func CloseLog() {
	fmt.Println("close logs")
	// if fd != nil {
	// 	logrus.Println("close files logs")
	// 	fd.Sync()
	// 	fd.Close()
	// } else {
	// 	fmt.Println("[Logs] No opening logs")
	// }
}
