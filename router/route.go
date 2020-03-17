package router

import (
	"time"

	"github.com/betacraft/yaag/yaag"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/controller"
	"github.com/zhxx123/gomonitor/middleware"
)

// BUG(zhxx): 我是bug说明

// 路由注册
func RegisterRoute(api *iris.Application) {
	const maxSize = 1 << 20
	// http 访问控制
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})
	//api 文档配置
	appName := config.APIDocConfig.Name
	appDoc := config.APIDocConfig.Doc
	appUrl := config.APIDocConfig.URL
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: appName,
		DocPath:  appDoc + "/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": appUrl,
			"Staging":    "linux",
		},
	})
	limiter := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})

	// This is an example on how to limit only GET and POST requests.
	limiter.SetMethods([]string{"GET", "POST"})

	// 接口
	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{
		// v1.Use(middleware.RefreshTokenCookie)
		v1.Use(middleware.LimitHandler(limiter))
		// v1.Use(middleware.NewYaag()) // <- IMPORTANT, register the middleware.
		apiPrefix := config.ServerConfig.APIPrefix
		// 通用接口
		v1.PartyFunc(apiPrefix, func(p router.Party) {
			p.Post("/register", controller.Signup)     //注册用户
			p.Post("/login", controller.Signin)        // 登录
			p.Post("/reset", controller.ResetPassword) //重置密码
			// 获取图片验证码
			p.Post("/captchacode", controller.GetCaptchaCode)
			// 验证图片验证码
			p.Post("/verifycode", controller.CheckCaptchaCode)
			// 获取手机或者邮箱验证码
			p.Post("/emailcode", controller.GetEmailPhoneCode)

			// 产品
			p.PartyFunc("/product", func(h router.Party) {
				h.Get("/list", controller.GetAllProduct)
				// h.Get("/detail/{id:string}", controller.GetProDetail)
			})
			// 帮助中心
			p.PartyFunc("/help", func(h router.Party) {
				h.Get("/type", controller.GetHelpTypeList)
				h.Get("/list", controller.GetHelpList)
				h.Get("/detail", controller.GetArticleDetail)

			})
			// 公告
			p.PartyFunc("/notice", func(n router.Party) {
				n.Get("/list", controller.GetNoticeList)
				n.Get("/detail", controller.GetArticleDetail)
			})

			// 访问信息提交
			p.Get("/stats/visit", controller.AddPV)
			// github webhook 自动更新
			p.Post("/notify/gitpush", controller.GithubHook)
			// github logs
			p.Get("/github", controller.GithubHookLog)
		})

		// users
		v1.PartyFunc("/user", func(u router.Party) {
			u.Use(middleware.JwtHandler().Serve, middleware.AuthToken)
			// users info
			u.PartyFunc("/users", func(users router.Party) {
				users.Get("/asset", controller.GetUserAsset)
				users.Get("/assetflow", controller.GetAssetFlow)
				users.Get("/profile", controller.GetUserProfile)
				users.Put("/username/update", controller.UpdateUserName)
				users.Put("/password/update", controller.UpdatePassword)
				users.Post("/logout", controller.Signout)
			})
			// user message
			u.PartyFunc("/message", func(m router.Party) {
				m.Get("/list", controller.GetMessageList)
				m.Get("/detail", controller.GetMessageDetail)
				m.Put("/update", controller.UpdateMessage)
			})

			// user order
			u.PartyFunc("/order", func(order router.Party) {
				order.Post("/id", controller.GetOrderID)
				order.Post("/precreate", controller.OrderCreate)
				order.Post("/pay", controller.OrderCreatePay)
				order.Get("/query", controller.OrderQuery)
				order.Post("/cancel", controller.OrderCancel)
				order.Post("/update", controller.OrderUpdate)
				// 租用订单 列表
				order.Get("/list", controller.GetUserOrderList)
			})
			// user pay
			u.PartyFunc("/pay", func(pays router.Party) {
				pays.Post("/id", controller.GetPayID)
				pays.Post("/precreate", controller.PayCreate)
				pays.Get("/query", controller.PayQuery)
				pays.Post("/cancel", controller.PayCancel)
				pays.Post("/update", controller.PayUpdate)
				// 支付订单列表
				// pays.Post("/list/{id:string}", controller.GetPayList)
			})
			// upload
			u.PartyFunc("/upload", func(upload router.Party) {
				upload.Post("/image", iris.LimitRequestBodySize(maxSize+1<<19), controller.UploadImage)
				upload.Post("/id", controller.GetWorkOrderID)
				upload.Post("/workorder", controller.WorkOrder)
			})

			// others
			u.PartyFunc("/login", func(l router.Party) {
				l.Get("/logs", controller.GetLoginList) //获取登录记录
			})
			// coins
			u.PartyFunc("/coin", func(coin router.Party) {
				coin.Get("/price", controller.GetCoinPrice)         //获取货币价格
				coin.Get("/pricelist", controller.GetCoinPriceList) //获取货币价格
			})

		})

		// admin
		v1.Post("/admin/login", controller.Signin)
		v1.PartyFunc("/admin", func(admin router.Party) {
			admin.Use(middleware.JwtHandler().Serve, middleware.AuthAdminToken)
			// system
			admin.PartyFunc("/sys", func(sys router.Party) {
				sys.Get("/basic", controller.GetAdminSysBasicListInfo)
				sys.Get("/simple", controller.GetAdminSysSimpleInfo)
			})
			// wallet
			admin.PartyFunc("/wallet", func(c router.Party) {
				c.Get("/basic", controller.GetAdminWalletBasicInfo)
				c.Get("/simple", controller.GetAdminWalletSimpleInfo)
				c.Get("/address", controller.GetAdminWalletAddressRecordList)
				c.Post("/upaddr", controller.UpdateAdminWalletAddress)
				c.Put("/coinaddr", controller.UpdateAdminUserCoinAddress)
			})
			// coinprice
			admin.PartyFunc("/coinprice", func(c router.Party) {
				c.Get("/list", controller.GetAdminCoinPrice)
				c.Get("/market", controller.GetAdminCoinMarketList)
				c.Put("/update", controller.UpdateAdminCoinPrice)
				c.Put("/updatestatus", controller.UpdateAdminCoinPriceAutoUp)
			})
			// users
			admin.PartyFunc("/users", func(users router.Party) {
				// 获取所有用户列表
				users.Get("/list", controller.GetAdminUserList)
				// 获取当前用户信息
				users.Get("/info", controller.GetAdminUserInfo)
				// 获取指定用户信息
				users.Get("/{id:int}", controller.GetAdminSpecificUserInfo)
				// 更新指定用户信息
				users.Put("/update", controller.UpdateAdminUserStatusRole)
				// 新建用户
				// users.Post("/create", controller.CreateAdminUser)
				// 删除用户
				users.Delete("/{id:int}", controller.DeleteAdminUser)
				// 下线用户
				users.Put("/logout/{id:int}", controller.LogoutAdminUser)
				// 用户登录日志记录
				users.Get("/loginlogs", controller.GetAdminLoginLogs)
				// 管理员用户登出
				users.Put("/logout", controller.UserAdminLogout)
			})

			// message
			admin.PartyFunc("/message", func(m router.Party) {
				m.Get("/list", controller.GetAdminMessageList) // 用户消息列表
				m.Post("/add", controller.AddAdminUserMessage)
				m.Put("/updatestatus", controller.UpdateAdminMessageStatus) // 更新指定用户的消息
				m.Put("/update", controller.UpdateAdminMessage)             // 更新指定用户的消息
			})
			// finances
			admin.PartyFunc("/finances", func(f router.Party) {
				f.Get("/summary", controller.GetAdminSummay)
				f.Get("/report", controller.GetAdminFinanceReportList)
				f.Post("/add", controller.AddAdminFinanceReport)
				f.Put("/updatestatus", controller.UpdateAdminReportStatus)
			})
			// 数字货币
			admin.PartyFunc("/digtial", func(f router.Party) {
				f.Get("/list", controller.GetAdminDigtialList)
				f.Post("/update", controller.UpdateAdminDigtial)
			})
			// 法币
			admin.PartyFunc("/office", func(f router.Party) {
				f.Get("/list", controller.GetAdminOfficeList)
				f.Post("/update", controller.UpdateAdminOffice)
			})
			// 资产
			admin.PartyFunc("/asset", func(f router.Party) {
				f.Get("/list", controller.GetAdminUserAssetList)
				f.Get("/flowlist", controller.GetAdminAssetFlow)
			})
			// 虚拟充值
			admin.PartyFunc("/virtual", func(f router.Party) {
				f.Get("/list", controller.GetAdminVirtualList)
				f.Post("/create", controller.AddAdminAssetVirtual)
			})
			// order 订单
			admin.PartyFunc("/orders", func(o router.Party) {
				o.Get("/list", controller.GetAdminOrderList)
				o.Put("/updatestatus", controller.UpdateAdminOrderStatus)
				o.Get("/query/{id:int}", controller.QueryAdminOrder)
				// 矿场机器订单列表
				o.Get("/miner", controller.GetAdminMinerList)
			})
			// commodits 商品
			admin.PartyFunc("/commodits", func(c router.Party) {
				// 矿机，云算力
				c.Get("/list", controller.GetAdminGoodsList)
				c.Post("/add", controller.CreateAdminGoods)
				c.Post("/update", controller.UpdateAdminGood)
				c.Put("/updatestatus", controller.UpdateAdminGoodStatus)
				c.Get("/query/{id:string}", controller.QueryAdminGood)
			})
			// farm 矿场
			admin.PartyFunc("/farms", func(f router.Party) {
				// 矿机，云算力
				f.Get("/list", controller.GetAdminFarmsList)
			})

			// article notice
			admin.PartyFunc("/article", func(a router.Party) {
				a.Get("/list", controller.GetAdminArticleList)
				a.Get("/detail", controller.QueryAdminArticle)
				a.Put("/update", controller.UpdateAdminArticle)
				a.Post("/add", controller.AddAdminArticle)
				a.Put("/updatestatus", controller.UpdateAdminArticleStatus)
				a.Put("/updatereadcount", controller.UpdateAdminArticleReadCount)
				a.Get("/helpcategory/list", controller.GetAdminHelpType)
				a.Post("/helpcategory/add", controller.AddAdminHelpType)
				a.Put("/helpcategory/update", controller.UpdateAdminHelpType)
				a.Put("/helpcategory/updatestatus", controller.UpdateStatusAdminHelpType)
			})
			// permission
			admin.PartyFunc("/permission", func(p router.Party) {
				p.Get("/list", controller.GetAdminUserMangerList)
			})
			// settings
			admin.PartyFunc("/settings", func(set router.Party) {
				set.Get("/list", controller.GetAdminSettingsList)
				set.Put("/update", controller.UpdateAdminSettings)
			})

			// logs
			admin.PartyFunc("/logs", func(log router.Party) {
				log.Get("/list", controller.GetAdminLogList)
				log.Get("/login", controller.GetAdminUserLoginLogs)
				log.Get("/github", controller.WebHookLog)
			})
		})
	}
}
