package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/zhxx123/gomonitor/utils"
)

var confMap map[string]interface{}

func initJSON() {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		os.Exit(-1)
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

	if err := json.Unmarshal(bytes, &confMap); err != nil {
		fmt.Println("invalid config: ", err.Error())
		os.Exit(-1)
	}

}

type dBConfig struct {
	Dialect      string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
	Charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
	AutoMigrated bool
}

// DBConfig 数据库相关配置
var DBConfig dBConfig

func initDB() {
	utils.SetStructByJSON(&DBConfig, confMap["database"].(map[string]interface{}))
	if DBConfig.Dialect == "mysql" {
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset)
		// url := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True&loc=Local",
		// 	DBConfig.User, DBConfig.Password, DBConfig.Database, DBConfig.Charset)
		DBConfig.URL = url
	} else { // 测试环境数据库
		url := fmt.Sprintf("%s.db", DBConfig.Database)
		DBConfig.URL = url
	}

}

// RedisConfig redis相关配置
type redisConfig struct {
	Host      string
	Port      int
	Password  string
	URL       string
	MaxIdle   int
	MaxActive int
}

var RedisConfig redisConfig

func initRedis() {
	utils.SetStructByJSON(&RedisConfig, confMap["redis"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%d", RedisConfig.Host, RedisConfig.Port)
	RedisConfig.URL = url
}

// mongodb 配置
type mongoConfig struct {
	URL      string
	Database string
}

var MongoConfig mongoConfig

func initMongo() {
	utils.SetStructByJSON(&MongoConfig, confMap["mongodb"].(map[string]interface{}))
}

// WebHookConfig 相关配置
type webHookConfig struct {
	WebHookShell bool
	MailProject  string
	MailAutoSend int
	MailAuthor   string
}

var WebHookConfig webHookConfig

func initWebHook() {
	utils.SetStructByJSON(&WebHookConfig, confMap["webhook"].(map[string]interface{}))
}

// serverConfig 相关配置
type serverConfig struct {
	Version            string
	Addr               string
	APITLSEnabled      bool
	TLSCertPath        string
	TLSKeyPath         string
	AutoSyncGoods      bool
	LogLevel           string
	LogOutput          string
	LogDir             string
	LogFilePrefix      string
	LogTimeInterval    int
	LogMaxNumber       int
	APIPoweredBy       string
	SiteName           string
	Host               string
	ImgHost            string
	Env                string
	APIPrefix          string
	UploadImgDir       string
	ImgPath            string
	MaxMultipartMemory int
	// Port               int
	DataStatsEnabled bool
	TokenSecret      string
	TokenMaxAge      int
	OrderMaxAge      int
	// PassSalt     string
	// LuosimaoVerifyURL string
	// LuosimaoAPIKey    string
	Github           string
	BaiduPushLink    string
	GithubAutoDeploy bool
	GithubSecretKey  string
}

var ServerConfig serverConfig

func initServer() {
	utils.SetStructByJSON(&ServerConfig, confMap["app"].(map[string]interface{}))
	sep := string(os.PathSeparator)
	execPath, _ := os.Getwd()
	length := utf8.RuneCountInString(execPath)
	lastChar := execPath[length-1:]
	if lastChar != sep {
		execPath = execPath + sep
	}
	if ServerConfig.UploadImgDir == "" {
		pathArr := []string{"static", "upload", "img"}
		uploadImgDir := execPath + strings.Join(pathArr, sep)
		ServerConfig.UploadImgDir = uploadImgDir
	}

	// ymdStr := utils.GetTodayYMD("-")

	// if ServerConfig.LogDir == "" {
	// 	ServerConfig.LogDir = execPath
	// } else {
	// 	length := utf8.RuneCountInString(ServerConfig.LogDir)
	// 	lastChar := ServerConfig.LogDir[length-1:]
	// 	if lastChar != sep {
	// 		ServerConfig.LogDir = ServerConfig.LogDir + sep
	// 	}
	// }
	// ServerConfig.LogFile = ServerConfig.LogDir + ymdStr + ".log"
}

// StatsDConfig statsd相关配置
type statsdConfig struct {
	URL     string
	Prefix  string
	Enabled bool
}

var StatsDConfig statsdConfig

func initStatsd() {
	utils.SetStructByJSON(&StatsDConfig, confMap["statsd"].(map[string]interface{}))
}

// rpc相关配置
type rpcConfig struct {
	MGDHost      string
	MGDUser      string
	MGDPassword  string
	MGDConfirm   int
	BTCHost      string
	BTCUser      string
	BTCPassword  string
	BTCConfirm   int
	ETHHost      string
	ETHConfirm   int
	PoolHost     string
	PoolUser     string
	PoolPassword string
}

var RPCConfig rpcConfig

func initRPC() {
	utils.SetStructByJSON(&RPCConfig, confMap["rpc"].(map[string]interface{}))
}

// api doc 配置
type apiDocConfig struct {
	Name string
	URL  string
	Doc  string
}

var APIDocConfig apiDocConfig

func initAPIDocConfig() {
	utils.SetStructByJSON(&APIDocConfig, confMap["apidoc"].(map[string]interface{}))
}

// pay 支付设置
type payConfig struct {
	Timeout  string
	PollTime int
}

var PayConfig payConfig

func initPayConfig() {
	utils.SetStructByJSON(&PayConfig, confMap["pay"].(map[string]interface{}))
}

// ip归属地相关配置
type ipLocConfig struct {
	TBAPI         string
	TCAPI         string
	TCAPIKey      string
	BDAPI         string
	MyExternalAPI string
	MyZhxxAPI     string
}

var IPLocConfig ipLocConfig

func initIPLocConfig() {
	utils.SetStructByJSON(&IPLocConfig, confMap["iploc"].(map[string]interface{}))
}

// 企业邮箱 相关配置
type emailConfig struct {
	MailUser string
	MailPass string
	MailHost string
	MailPort int
	MailFrom string
}

var EmailConfig emailConfig

func initEmailConfig() {
	utils.SetStructByJSON(&EmailConfig, confMap["email"].(map[string]interface{}))
}

// 初始化
func init() {
	initJSON()
	initDB()
	initRedis()
	initMongo()
	initWebHook()
	initServer()
	initStatsd()
	initRPC()
	initAPIDocConfig()
	initIPLocConfig()
	initEmailConfig()
	initPayConfig()
}
