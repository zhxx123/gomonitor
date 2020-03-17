package model

import "github.com/jinzhu/gorm"

// 资产报表
type Assets struct {
	gorm.Model
	FlowId         uint   `gorm:"not null; default 0; type:int(10)" json:"flow_id"`              //报表最终号，资产流水ID
	CoinType       string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`       //报表资产类型
	Count          int    `gorm:"not null; default 0; type:int(10)" json:"count"`                //统计区间内交易笔数
	IncreaceAmount string `gorm:"not null; default ''; type:varchar(32)" json:"increace_amount"` //统计区间金额变化
	ReduceAmount   string `gorm:"not null; default ''; type:varchar(32)" json:"reduce_amount"`   //统计区间金额变化
	TotalAmount    string `gorm:"not null; default ''; type:varchar(32)" json:"total_amount"`    //当前系统总额
	StartAt        int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"start_at"`          //报表开始时间
	EndAt          int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"end_at"`            //报表截至时间
	TotalTime      int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"total_time"`        //统计区间时长
	AuthorId       int    `gorm:"not null; default 0; type:int(10)" json:"author_id"`            //当前报表审核者ID
	CreateAt       int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"create_at"`         //创建时间
	Status         bool   `gorm:"not null; default 0; type:tinyint(1)" json:"status"`              //是否已经审核，只有审核过的，才可以继续统计
}

// 系统资产表
type SystemAccount struct {
	gorm.Model
	CoinType   string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`
	CoinAmount string `gorm:"not null; default ''; type:varchar(32)" json:"coin_amount"`
	UpdateAt   int64  `gorm:"not null; default 0; type:BIGINT(20)" json:"update_at"` //创建时间
}

// 获取资产报表
type AReportJson struct {
	CoinType string `json:"coinType"`
	Status   bool   `json:"status"`
	IsStatus bool   `json:"isstatus"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// 更新报表状态
type AUpdateReportStatusJson struct {
	ID int `json:"id"`
	// AuthorId int  `json:"authorId"`
	Status bool `json:"status"`
}

// 新增报表
type AUpdateReportsJson struct {
	// AuthorId int    `json:"authorId"`
	CoinType string `json:"coinType"`
	EndAt    string `json:"endAt"`
}
