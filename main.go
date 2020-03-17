package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/zhxx123/gomonitor/controller"
	"github.com/zhxx123/gomonitor/middleware"
	"github.com/zhxx123/gomonitor/model"
)

// 初始化 iris app
func InitIris() {
	api := iris.New()
	api.Logger().SetLevel("info")
	api.Use(logger.New())
	// api 调用统计
	api.Use(middleware.APIStatsD())
	// 404 错误
	api.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(controller.ApiResource(model.NotFound, nil, "404 Not Found"))
	})
	// 500 错误
	api.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.WriteString("Oups something went wrong, try again")
	})
	// iris 中断获取
	iris.RegisterOnInterrupt(func() {
	})
	const maxSize = 1 << 20
	// http 访问控制
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	v1 := api.Party("/api", crs).AllowMethods(iris.MethodOptions)
	{
		// 通用接口
		v1.Post("/login", controller.Login)
		v1.Post("/logout", controller.Logout)
		v1.Get("/article/list", controller.ArtilceList)
		v1.Get("/article/info", controller.ArticleInfo)
	}

	addr := ":7000"
	if err := api.Run(iris.Addr(addr)); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	return
}

// 注册 信号
func signalNotify() {
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-c
		fmt.Printf("[got signal: %s], exiting goblog now\n", s)
		shutdownAll()
	}()
}

// 关闭所有
func shutdownAll() {
	fmt.Println("[All] Shutdown")
	os.Exit(0)
}

// 初始化
func init() {
	// 初始化日志
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

// 系统开始，主函数
func main() {
	signalNotify() // 注册监听程序
	InitIris()     // web 启动
}
