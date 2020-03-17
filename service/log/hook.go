package log

import (
	"fmt"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/zhxx123/gomonitor/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewLfsHook(logLevel string) {
	// log := logrus.New()
	logPath := config.ServerConfig.LogDir
	logTimeInterval := config.ServerConfig.LogTimeInterval
	// logMaxNumber := config.ServerConfig.LogMaxNumber
	logFilePrefix := config.ServerConfig.LogFilePrefix

	baseLogPath := path.Join(logPath, logFilePrefix)

	writer, err := rotatelogs.New(
		baseLogPath+"_%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath), // 生成软链，指向最新日志文件

		// rotatelogs.WithMaxAge(logTimeInterval),        // 文件最大保存时间
		// rotatelogs.WithRotationCount(lgoMaxNumber),  // 最多存365个文件

		rotatelogs.WithRotationTime(time.Hour*time.Duration(logTimeInterval)), // 日志切割时间间隔
	)
	if err != nil {
		fmt.Printf("config local file system logger error. %+v", errors.WithStack(err))
	}
	level, ok := logrus.ParseLevel(logLevel)
	if ok == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	// lfsHook := lfshook.NewHook(lfshook.WriterMap{
	// 	logrus.DebugLevel: writer,
	// 	logrus.InfoLevel:  writer,
	// 	logrus.WarnLevel:  writer,
	// 	logrus.ErrorLevel: writer,
	// 	logrus.FatalLevel: writer,
	// 	logrus.PanicLevel: writer,
	// }, &logrus.TextFormatter{
	// 	DisableColors:   true,
	// 	TimestampFormat: "2006-01-02 15:04:05", //时间格式化
	// 	FullTimestamp:   true})
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式化
		FullTimestamp:   true,
	})
	logrus.SetOutput(writer)
	return
}
