package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// KeyValueConfig key, value配置
type KeyValueConfig struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	KeyName   string     `json:"key"`
	Value     string     `json:"value"`
}

type KeyValueData struct {
	KeyName string `json:"key" validate:"required,min=1"`
	Value   string `json:"value" validate:"required,min=1"`
}

var (
	CHECK_VERIFY_CODE = false // 是否校验邮箱或者手机验证码
	PAUSE_TRADE       = false
)

// 系统设置表
type Settings struct {
	gorm.Model
	AuthorId int    `gorm:"not null; default ''; type:int(10)" json:"author_id"`    //更改者id
	Category string `gorm:"not null; default ''; type:varchar(32)" json:"category"` // 类型
	Name     string `gorm:"not null; default ''; type:varchar(32)" json:"name"`     // 项名
	Value    string `gorm:"not null; default ''; type:text" json:"value"`           // json字符串值
}

type ASettingsUpdateJson struct {
	AuthorId int    `json:"authorId"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type ASettingsJson struct {
	Category string `json:"category"`
	Name     string `json:"name"`
}
