package model

import "github.com/jinzhu/gorm"

// type BasicInfoDataRes struct {
// 	Data   *[]SystemBasic `json:"data"`
// 	Length int            `json:"length"`
// }
// type SimpleInfoDataRes struct {
// 	Data   *[]SystemSimple `json:"data"`
// 	Length int             `json:"length"`
// }

type LoginUsers struct {
	LoginName string `gorm:"not null; default ''; type:varchar(20)" json:"login_name"`
	TTY       string `gorm:"not null; default ''; type:varchar(20)" json:"tty"`
	LoginTime string `gorm:"not null; default ''; type:varchar(20)" json:"login_time"`
	LoginAddr string `gorm:"not null; default ''; type:varchar(20)" json:"login_addr"`
	Detail    string `gorm:"not null; default ''; type:varchar(20)" json:"detail"`
}
type Basics struct {
	MyUserId
	NickName    string `gorm:"not null; default ''; type:varchar(64)" json:"nick_name"`
	NetRemoteIp string `gorm:"not null; default ''; type:varchar(20)" json:"net_remote_ip"`
	CpuCount    int    `gorm:"not null; default 0; type:int(10)" json:"cpu_count"`
	CpuName     string `gorm:"not null; default ''; type:varchar(64)" json:"cpu_name"`
	OsArch      string `gorm:"not null; default ''; type:varchar(20)" json:"os_arch"`
	OsByteOrder string `gorm:"not null; default ''; type:varchar(20)" json:"os_byte_order"`
	OsSystem    string `gorm:"not null; default ''; type:varchar(20)" json:"os_system"`
	NetMac      string `gorm:"not null; default ''; type:varchar(20)" json:"net_mac"`
	DiskTotal   string `gorm:"not null; default ''; type:varchar(20)" json:"disk_total"`
	MemTotal    string `gorm:"not null; default ''; type:varchar(20)" json:"mem_total"`
}
type SystemBasic struct {
	gorm.Model
	MyUserId
	NickName    string `gorm:"not null; default ''; type:varchar(64)" json:"nick_name"`
	NetRemoteIp string `gorm:"not null; default ''; type:varchar(20)" json:"net_remote_ip"`
	CpuCount    int    `gorm:"not null; default 0; type:int(10)" json:"cpu_count"`
	CpuName     string `gorm:"not null; default ''; type:varchar(64)" json:"cpu_name"`
	OsArch      string `gorm:"not null; default ''; type:varchar(20)" json:"os_arch"`
	OsByteOrder string `gorm:"not null; default ''; type:varchar(20)" json:"os_byte_order"`
	OsSystem    string `gorm:"not null; default ''; type:varchar(20)" json:"os_system"`
	NetMac      string `gorm:"not null; default ''; type:varchar(20)" json:"net_mac"`
	DiskTotal   string `gorm:"not null; default ''; type:varchar(20)" json:"disk_total"`
	MemTotal    string `gorm:"not null; default ''; type:varchar(20)" json:"mem_total"`
	UpdateAt    int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"`
}
type Simples struct {
	MyUserId
	NetRemoteIp string        `json:"net_remote_ip"`
	NetLocalIp  string        `json:"net_local_ip"`
	MemUsed     string        `json:"mem_used"`
	MemUsage    string        `json:"mem_usage"`
	SysUptime   string        `json:"sys_uptime"`
	DiskUsed    string        `json:"disk_used"`
	DiskUsage   string        `json:"disk_usage"`
	NetByteSent int64         `json:"net_byte_sent"`
	NetByteRecv int64         `json:"net_byte_recv"`
	CpuAverage  *[]string     `json:"cpu_average"`
	LoginCount  int           `json:"login_count"`
	LoginUser   *[]LoginUsers `json:"login_user"`
}
type SystemSimple struct {
	gorm.Model
	MyUserId
	NetRemoteIp string `gorm:"not null; default '' ; type:varchar(20)" json:"net_remote_ip"`
	NetLocalIp  string `gorm:"not null; default '' ; type:varchar(20)" json:"net_local_ip"`
	MemUsed     string `gorm:"not null; default ''; type:varchar(20)" json:"mem_used"`
	MemUsage    string `gorm:"not null; default ''; type:varchar(12)" json:"mem_usage"`
	SysUptime   string `gorm:"not null; default ''; type:varchar(64)" json:"sys_uptime"`
	DiskUsed    string `gorm:"not null; default ''; type:varchar(20)" json:"disk_used"`
	DiskUsage   string `gorm:"not null; default ''; type:varchar(12)" json:"disk_usage"`
	NetByteSent int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"net_byte_sent"`
	NetByteRecv int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"net_byte_recv"`
	CpuAverage  string `gorm:"not null; default ''; type:varchar(128)" json:"cpu_average"`
	LoginCount  int    `gorm:"not null; default 0; type:int(10)" json:"login_count"`
	LoginUser   string `gorm:"not null; default ''; type:varchar(1024)" json:"login_user"`
	UpdateAt    int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"`
}
type SysInfoJson struct {
	UID   string `json:"uid"`
	Page  int    `json:"page" validate:"number,min=0,max=100"`
	Limit int    `json:"limit" validate:"number,min=0,max=100"`
}
