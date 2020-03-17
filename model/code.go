package model

type errorCode struct {
	SUCCESS      int
	ERROR        int
	NotFound     int
	LoginError   int
	LoginTimeout int
	InActive     int
}

// ErrorCode 错误码
var ErrorCode = errorCode{
	SUCCESS:      0,
	ERROR:        1,
	NotFound:     404,
	LoginError:   1000, //用户名或密码错误
	LoginTimeout: 1001, //登录超时
	InActive:     1002, //未激活账号
}

const (
	STATUS_SUCCESS        = 2000 // 成功
	STATUS_FAILED         = 2001 // 通用失败
	STATUS_EXPIRED        = 2002 // 登录过期
	STATUS_AUTH_ERROR     = 2003 // 没有权限
	STATUS_PARAM_ERROR    = 2004 // 参数错误
	STATUS_FREQUENT_LIMIT = 2005 // 访问过快
)
