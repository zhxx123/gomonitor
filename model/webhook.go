package model

import "github.com/jinzhu/gorm"

type WebHook struct {
	gorm.Model
	Type        string `gorm:"not null; default ''; type:varchar(32)" json:"type"`         // 类型,例如: github
	ProjectName string `gorm:"not null; default ''; type:varchar(32)" json:"project_name"` // 项目名称 nodeweb
	Refer       string `gorm:"not null; default ''; type:varchar(64)" json:"refer"`        // 类型,例如: master
	AuthorName  string `gorm:"not null; default ''; type:varchar(32)" json:"author_name"`  // 作者用户名
	AuthorEmail string `gorm:"not null; default ''; type:varchar(32)" json:"author_email"` // 作者邮箱
	Forced      bool   `gorm:"not null; default 0; type:tinyint(1)" json:"forced"`           // 是否强制提交
	BeforeID    string `gorm:"not null; default ''; type:varchar(64)" json:"before_id"`    // 上一次提交id
	AfterID     string `gorm:"not null; default ''; type:varchar(64)" json:"after_id"`     // 上一次提交id
	Status      bool   `gorm:"not null; default 0; type:tinyint(1)" json:"status"`           // 本地是否执行成功
	Output      string `gorm:"not null; default ''; type:varchar(512)" json:"output"`      // 本地是否执行成功
	TimeStamp   string `gorm:"not null; default ''; type:varchar(32)" json:"time_stamp"`   // 本地是否执行成功
	Time        int64  `gorm:"not null; default 0; type:varchar(64)" json:"time"`          //执行时间
}

// json
type QueryHookJson struct {
	ProjectName string `json:"projectName"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}
