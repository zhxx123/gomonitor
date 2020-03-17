package model

const (
	SUCCESS  = 0 // 成功
	ERROR    = 1 // 通用失败
	NotFound = 2
)

const StatsdURL = "10.0.0.5:8125" // 填写你的 statsd 服务端接收地址
const StatsdPreFix = ""           // statsd 客户端前缀，可留空
