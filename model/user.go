package model

import (
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

const (
	// 验证码链接有效期
	ActiveDuration = 5 * 60
	ResetDuration  = 10 * 60
)

const (
	// UserRoleNormal 普通用户
	UserRoleNormal = 1

	// UserRoleEditor 网站编辑
	UserRoleEditor = 2

	// UserRoleAdmin 管理员
	UserRoleAdmin = 3

	// UserRoleSuperAdmin 超级管理员
	UserRoleSuperAdmin = 4

	// UserRoleCrawler 爬虫，网站编辑或管理员登陆后台后，操作爬虫去抓取文章
	// 这时，生成的文章，其作者是爬虫账号。没有直接使用爬虫账号去登陆的情况.
	UserRoleCrawler = 5
)

const (
	UserStatusNormal = 0
)

const (
	// MaxPassLen 密码的最大长度
	MaxPassLen = 20

	// MinPassLen 密码的最小长度
	MinPassLen = 6
)

const (
	WORKORDER_PROCESS_STATUS   = 1201 // process 处理中
	WORKORDER_COMPLETED_STATUS = 1202 // completed 已完成
)
const (
	TRADE_TYPE_INPUT    = iota + 1 //用户资产增加，也就是平台资产增加
	TRADE_TYPE_USER_OUT            //用户资产减少
	TRADE_TYPE_SYS_OUT             // 平台资产减少
)

// 用户注册数据库表
type User struct {
	gorm.Model

	Username    string `gorm:"not null default '';type:varchar(30)"`
	UserID      int    `gorm:"unique;not null default 0; type:int(10) "`
	Email       string `gorm:"unique;not null default '' ;type:varchar(20)"`
	Password    string `gorm:"not null default '';type:varchar(255)"`
	Phone       string `gorm:"not null default '';type:varchar(20)"`
	RoleID      int    `gorm:"not null default 0; type:int(10)"`
	RegisterAt  int64  `gorm:"not null default 0 BIGINT(20)"`
	RegisterIp  string `gorm:"not null default '' ;type:varchar(50)"`
	LastLoginAt int64  `gorm:"not null; default 0; BIGINT(20)"`
	UserStatus  int64  `gorm:"not null default 0 ;type:BIGINT(20)"`
}

// CheckPassword 验证密码是否正确
func (user User) CheckPassword(inPass string) bool {
	if inPass == "" || user.Password == "" {
		return false
	}
	return bcrypt.Match(inPass, user.Password)
}

// EncryptPassword 给密码加密
func (user User) EncryptPassword(password string) string {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)
	return hash
}

// 用户登录认证,token 管理
type UserOauth struct {
	gorm.Model
	UserId    int    `gorm:"not null; default 0; comment:'UserId'; int(10)" json:"user_id"`
	RoleId    int    `gorm:"not null; default 0; comment:'RoleId'; int(10)" json:"role_id"`
	LoginAt   int64  `gorm:"not null; default 0; comment:'LgoinAt'; BIGINT(20)" json:"login_at"`
	Token     string `gorm:"not null; default ''; comment:'Token'; type:varchar(255)" json:"token"`
	Secret    string `gorm:"not null; default ''; comment:'Secret'; type:varchar(255)" json:"secret"`
	LoginType string `gorm:"not null; default ''; comment:'LoginType'; type:varchar(50)" json:"login_type"` // 登录设备类型
	LoginIp   string `gorm:"not null; default ''; comment:'LoginIp'; type:varchar(50)" json:"login_ip"`
	LoginCity string `gorm:"not null; default ''; comment:'LoginCity'; type:varchar(50)" json:"login_city"`
	ExpressIn int64  `gorm:"not null; default 0; comment:'ExpressIn'; BIGINT(20)" json:"express_in"`
	Revoked   bool   `gorm:"not null; default 0; comment:'Revoked'; tinyint(1)" json:"revoked"`
}

// UserAccounts 用户资产余额表
type UserAccounts struct {
	gorm.Model
	UserId        int    `gorm:"not null; default 0; type:int(10)" json:"user_id"`
	CoinType      string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`
	CoinAddr      string `gorm:"not null; default ''; type:varchar(80)" json:"coin_addr"`
	CoinAmount    string `gorm:"not null; default '0'; type:varchar(30)" json:"coin_amount"`
	VirtualAmount string `gorm:"not null; default '0'; type:varchar(30)" json:"virtual_amount"`
}

// 用户消息表
type UserMessage struct {
	gorm.Model
	UserId    int    `gorm:"not null; default 0; type:int(10)" json:"user_id"`           // 消息作者id
	ArticleId int    `gorm:"unique;not null; default 0; type:int(10)" json:"article_id"` // 文章唯一标识符号
	Category  int    `gorm:"not null; default 0; type:int(2) " json:"category"`          //消息类别，充值，消费，其他
	Title     string `gorm:"not null; default ''; type:varchar(128)" json:"title"`       // 文章标题
	Summary   string `gorm:"not null; default ''; type:varchar(128)" json:"summary"`     // 文章摘要
	Content   string `gorm:"not null; default ''; type:text" json:"content"`
	PushedAt  int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"pushed_at"`   // 发布时间
	Announcer string `gorm:"not null; default ''; type:varchar(20)" json:"announcer"` // 发布者标签 xxx团队
	AuthorId  int    `gorm:"not null; default 0; type:int(10)" json:"author_id"`      // 当前消息作者ID
	Readed    bool   `gorm:"not null; default 0; type:tinyint(1)" json:"read_ed"`     // 当前消息是否已经阅读
	ReadAt    int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"read_at"`
	Status    bool   `gorm:"not null; default 0; type:tinyint(1)" json:"status"` // 当前消息是否已经发布
}

// 用户登录设置表
type UserLoginSets struct {
	gorm.Model
	UserId     int    `gorm:"not null; default 0; type:int(10) "`     // 用户id
	LoginCount int    `gorm:"not null; default 0; type:int(10) "`     // 登录失败次数
	LoginAt    int64  `gorm:"not null; default 0; type:BIGINT(20)"`   // 首次登录失败时间
	IPList     string `gorm:"not null; default 0; type:varchar(255)"` // 常用登录ip地址列表，冒号分隔符号存储;
	Verify     bool   `gorm:"not null; default 0; type:tinyint(1)"`   // 是否需要开启验证码验证
}

// 用户提交工单记录表
type WorkOrder struct {
	gorm.Model
	UserId      int    `gorm:"not null; default 0; type:int(10)"`              // 用户id
	WorkId      string `gorm:"unique;not null; default '0'; type:varchar(30)"` // 工单标识
	IssueType   int    `gorm:"not null; default 0; type:int(10)"`              // 问题类型
	Description string `gorm:"not null; default ''; type:varchar(255)"`        // 问题描述
	Email       string `gorm:"not null; default ''; type:varchar(50)"`         // 联系方式
	ImgUri      string `gorm:"not null; default ''; type:varchar(128)"`        // 工单附件图片名称，分号分割
	CreateAt    int64  `gorm:"not null; default 0; type:BIGINT(20)"`           // 工单创建时间
	Status      int    `gorm:"not null; default 0; type:int(10)"`              // 工单状态
	OperatorId  int    `gorm:"not null; default 0; type:int(10) "`             // 处理者 id
}

// 用户挖矿账户表
type MinerAccounts struct {
	gorm.Model
	UserId      int    `gorm:"not null; default 0; type:int(10)" json:"user_id"`            // 用户id
	CoinType    string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`     // 币种类型
	MinerPool   string `gorm:"not null; default ''; type:varchar(80)" json:"miner_worker"`  //矿池链接
	CoinAddr    string `gorm:"not null; default ''; type:varchar(256)" json:"coin_addr"`    // 地址,或者 矿工名称
	AddrType    int    `gorm:"not null; default 0; type:int(10)" json:"addr_type"`          // 地址类型 1: 属于用户 2: 属于系统
	CoinAmount  string `gorm:"not null; default '0'; type:varchar(30)" json:"coin_amount"`  // 地址当前剩余收益(针对与属于系统的地址)
	TotalAmount string `gorm:"not null; default '0'; type:varchar(30)" json:"total_amount"` // 地址累计收益
	Status      int    `gorm:"not null; default 0; type:int(10)" json:"status"`             // 当前状态是否启用 1启用 其他不启用
}

// 用户登录提交表单
type UserLogin struct {
	LoginType   string `json:"loginType" validate:"required"`
	SigninInput string `json:"signinInput" validate:"required"`
	Password    string `json:"password" validate:"required,gte=6,lte=20"`
	VerifyCode  string `json:"verifyCode"`
	UID         string `json:"uid"` // 验证码申请时候的 uid
}

// 用户登录提交表单
type UserUpdatePwd struct {
	Type         string `json:"type" validate:"required"`
	AddressInput string `json:"addressInput" validate:"required"`
	Password     string `json:"password" validate:"required,gte=6,lte=20"`
	VerifyCode   string `json:"verifyCode"`
	UID          string `json:"uid"` // 验证码申请时候的 uid
}

//用户注册提交表单 struct
type UserJson struct {
	// Username   string `json:"username" validate:"gte=0,lte=20"`
	Email      string `json:"email" validate:"required"`
	Phone      string `json:"phone"`
	Password   string `json:"password" validate:"required,gte=6,lte=50"`
	VerifyCode string `json:"verifyCode"`
	UID        string `json:"uid"` // 验证码申请时候的 uid
}

// 更新用户名(昵称)
type UserUpdateJson struct {
	Username string `json:"username"`
}

// 更新密码
type PasswordUpdateData struct {
	Password string `json:"password" validate:"required,min=6,max=20"`
	NewPwd   string `json:"newPwd" validate:"required,min=6,max=20"`
}

// 用户账户返回表
type UserAccountRes struct {
	CoinType      string `json:"coin_type"`
	CoinAddr      string `json:"coin_addr"`
	CoinAmount    string `json:"coin_amount"`
	VirtualAmount string `json:"virtual_amount"`
}

// 用户获取图片验证码
type CaptchaCodeJson struct {
	UID string `json:"uid" validate:"required,gte=30,lte=36"`
	// 1:数字验证码, 2:音频验证码 3:公式验证码
	Type int `json:"type" validate:"required,gte=0,lte=2"`
}
type CaptchaCodeRes struct {
	UID  string `json:"uid"`
	Data string `json:"data"`
}

// 校验图片验证码
type CodeJson struct {
	UID  string `json:"uid" validate:"required,gte=20,lte=36"`
	Code string `json:"code" validate:"required"`
}

// 用户获取手机或者邮箱验证码
type VerifyCodeJson struct {
	UID     string `json:"uid" validate:"required,gte=32,lte=32"`
	Address string `json:"address" validate:"required"`            // 验证码接收方地址
	Type    int    `json:"type" validate:"required,gte=0,lte=2"`   // 1:email, 2:phone
	Option  int    `json:"option" validate:"required,gte=1,lte=2"` //1: 用户注册 2: 找回密码
}

// 用户验证码返回
type VerifyCodeRes struct {
	UID        string `json:"uid"`
	VerifyCode string `josn:"verify_code"`
}

// 用户消息列表
type UserMessageJson struct {
	// UserId int `json:"user_id"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
type UserMessageRes struct {
	ArticleId int    `json:"article_id"` // 文章唯一标识符号
	Category  int    `json:"category"`   // 文章类型
	Title     string `json:"title"`      // 文章标题
	PushedAt  int64  `json:"pushed_at"`  // 发布时间
	Readed    bool   `json:"readed"`     // 消息是否已经阅读
}

// 用户消息详情
type MessageDetailJson struct {
	ArticleId int `json:"articleId"` // 文章唯一标识符号
}
type MessageDetailRes struct {
	ArticleId int    `json:"article_id"` // 文章唯一标识符号
	PushedAt  int64  `json:"pushed_at"`  // 发布时间
	Announcer string `json:"announcer"`  // 发布者标签 xxx团队
	Title     string `json:"title"`      // 文章标题
	Summary   string `json:"summary"`    // 文章摘要
	Content   string `json:"content"`    // 文章内容
	Readed    bool   `json:"readed"`     // 消息是否已经阅读
}
type WorkOrderIDJson struct {
	UID string `json:"uid" validate:"required,gte=6,lte=36"`
}

// 用户提交工单
type WorkOrderJson struct {
	UID         int    `json:"uid" validate:"required"`         // uuid 订单ID
	WorkId      string `json:"workId" validate:"required"`      // "xxxxx" 工单标识
	IssueType   int    `json:"issueType"  validate:"required"`  // 1 问题类型,1:"订单",2:"安全"
	Description string `json:"decription"  validate:"required"` // "xxxx" 问题描述
	Email       string `json:"email"  validate:"required"`      // "aa@qq.com" 联系方式
	ImgUri      string `json:"imgUri"`                          // "url" 工单附件图片名称，分号分割
	// OperatorId  int    `json:"operateor_id"` // 处理者 id
}

// 获取用户登录记录
type UserOauthJson struct {
	UserId int `json:"userId"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
}

// 用户登录记录简要
type UserOauthRes struct {
	LoginAt   int64  `json:"login_at"`
	LoginIp   string `json:"login_ip"`
	LoginCity string `json:"login_city"`
	LoginType string `json:"login_type"` // 登录设备类型
	// ExpressIn int64  `json:"express_in"`
}
type UserOauthDataRes struct {
	Data   *[]UserOauthRes `json:"data"`
	Length int             `json:"length"`
}

type AssetflowJson struct {
	StartAt  int64  `json:"startAt"`  //开始时间
	EndAt    int64  `json:"endAt"`    //结束时间
	CoinType string `json:"coinType"` // 类型
	Page     int    `json:"page" validate:"required,number,min=1,max=100"`
	Limit    int    `json:"limit" validate:"required,number,min=1,max=20"`
}
type AssetflowRes struct {
	OutTradeNo  string `json:"out_trade_no"` // 订单号
	TradeType   int    `json:"trade_type"`   // 交易类型
	CreateAt    int64  `json:"create_at"`    // 创建时间
	CoinType    string `json:"coin_type"`    // 充值资产类型，0: cny; 1: mgd; 2: eth
	Amount      string `json:"amount"`       // 账户金额
	TotalAmount string `json:"total_amount"` // 当前账户金额
	Description string `json:"description"`  // 交易字符串描述
}

// ********************************************************
// admin
// ********************************************************
// 用户列表获取
type AUserJson struct {
	Email    string `json:"email"`
	RoleType int    `json:"roleType"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}
type AUserInfoRes struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	UserId   int    `json:"user_id"`
	Email    string `json:"email"`
	// Password   string `json:"password"`
	Phone      string `json:"phone"`
	Roles      []int  `json:"roles"`
	RegisterAt int64  `json:"register_at"`
	RegisterIp string `json:"register_ip"`
	UserStatus int64  `json:"user_status"`
}
type AUserRes struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	UserId     int    `json:"user_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	RoleID     int    `json:"role_id"`
	RegisterAt int64  `json:"register_at"`
	RegisterIp string `json:"register_ip"`
	UserStatus int64  `json:"user_status"`
}
type AUserDataRes struct {
	Data   *[]AUserRes `json:"data"`
	Length int         `json:"length"`
}

// 获取用户消息
type AUserMessageJson struct {
	UserId   int  `json:"userId"`
	Status   bool `json:"status"`
	IsPushed bool `json:"isPushed"`
	Page     int  `json:"page"`
	Limit    int  `json:"limit"`
}
type AUserMessageDataRes struct {
	Data   *[]UserMessage `json:"data"`
	Length int            `json:"length"`
}

// 添加 用户消息
type AEditUserMessageJson struct {
	ArticleId int    `json:"articleId"` // 文章id，只在修改时候有用
	UserId    int    `json:"userId"`    // 消息作者id
	Category  int    `json:"category"`  // 消息类别
	Title     string `json:"title"`     // 文章标题
	Summary   string `json:"summary"`   // 文章摘要
	Content   string `json:"content"`   // 文章内容
	PushedAt  int64  `json:"pushedAt"`  // 发布时间
	Announcer string `json:"announcer"` // 发布者标签 xxx团队
	AuthorId  int    `json:"authorId"`  // 当前消息作者ID
}

// 更新消息
type AUpdateUserMessageJson struct {
	UserId    int  `json:"userId"`
	ArticleId int  `json:"articleId"`
	Status    bool `json:"status"`
}

type AUserAccountsJson struct {
	UserId   int    `json:"userId"`
	CoinType string `json:"coinType"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type AUserAccountsDataRes struct {
	Data   *[]UserAccounts `json:"data"`
	Length int             `json:"length"`
}

// 更新用户信息
type AUserUpdateStatusJson struct {
	UserId     int    `json:"userId"`
	Email      string `json:"email"`
	RoleID     int    `json:"roleId"`
	UserStatus int64  `json:"userStatus"`
}

// 用户登录日志
type ALoginLogsJson struct {
	UserId    int  `json:"userId"`
	Revoked   bool `json:"revoked"`
	IsRevoked bool `json:"isRevoked"`
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
}

type AUserOauthRes struct {
	ID      int   `json:"id"`
	UserId  int   `json:"user_id"`
	RoleId  int   `json:"role_id"`
	LoginAt int64 `json:"login_at"`
	// Token     string `json:"-"` // "-" 表示忽略此字段
	// Secret    string `json:"secret"`
	LoginType string `json:"login_type"`
	LoginIp   string `json:"login_ip"`
	LoginCity string `json:"login_city"`
	ExpressIn int64  `json:"express_in"`
	Revoked   bool   `json:"revoked"`
}

type ALoginLogsResData struct {
	Data   *[]AUserOauthRes `json:"data"`
	Length int              `json:"length"`
}
