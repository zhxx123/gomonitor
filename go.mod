module github.com/zhxx123/gomonitor

go 1.13

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/betacraft/yaag v1.0.0
	github.com/cactus/go-statsd-client/statsd v0.0.0-20191106001114-12b4e2b38748
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/fasthttp-contrib/websocket v0.0.0-20160511215533-1f3b11f56072 // indirect
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/garyburd/redigo v1.6.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/iris-contrib/middleware/cors v0.0.0-20191219204441-78279b78a367
	github.com/iris-contrib/middleware/jwt v0.0.0-20191219204441-78279b78a367
	github.com/jameskeane/bcrypt v0.0.0-20120420032655-c3cd44c1e20f
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/jordan-wright/email v0.0.0-20200307200233-de844847de93
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/iris v11.1.1+incompatible
	github.com/kataras/iris/v12 v12.1.8
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/mojocn/base64Captcha v1.3.0
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robfig/cron v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/tebeka/strftime v0.1.3 // indirect
	github.com/valyala/fasthttp v1.9.0 // indirect
	github.com/yudai/pp v2.0.1+incompatible // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

replace github.com/mojocn/base64Captcha v1.3.0 => github.com/mgdzhxx/base64Captcha v1.4.0
