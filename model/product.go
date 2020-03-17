package model

import "github.com/jinzhu/gorm"

type Products struct {
	gorm.Model
	AuthorId       int            `gorm:"not null; default 0; type:int(10)" json:"author_id"`              // 发布者id
	GoodsType      int            `gorm:"not null; default 0; type:int(10)" json:"goods_type"`             // 商品类型，1:云算力，2:矿机
	GoodsId        string         `gorm:"unique;not null; default ''; type:varchar(32)" json:"goods_id"`   // 商品id
	GoodsName      string         `gorm:"not null; default ''; type:varchar(32)" json:"goods_name"`        // 商品名称
	OrgPrice       string         `gorm:"not null; default ''; type:varchar(30)" json:"org_price"`         // 原始价格
	CurPrice       string         `gorm:"not null; default ''; type:varchar(30)" json:"cur_price"`         // 当前价格 一年期的价格为基准价，90天就*1.05,30天就*1.1
	Quantity       int            `gorm:"not null; default 0; type:int(10)" json:"quantity"`               // 商品数量
	TotalQuantity  int            `gorm:"not null; default 0; type:int(10)" json:"total_quantity"`         // 商品总数量
	Unit           string         `gorm:"not null; default ''; type:varchar(10)" json:"unit"`              // 商品单位
	Description    string         `gorm:"not null; default ''; type:varchar(30)" json:"description"`       // 商品介绍
	ImageUri       string         `gorm:"not null; default ''; type:varchar(128)" json:"image_uri"`        // 商品图片路径
	PushedAt       int64          `gorm:"not null; default 0; type:BIGINT(20)" json:"pushed_at,omitempty"` // 发布时间
	FarmID         string         `gorm:"not null; default ''; type:varchar(64)" json:"farm_id"`           // 矿场ID标识
	MinerGoodsType string         `gorm:"not null; default ''; type:varchar(64)" json:"miner_goods_type"`  // 矿场商品标识
	Status         bool           `gorm:"not null; default 0; type:tinyint(1)" json:"status"`                // 当前商品是否已经发布
	Details        ProductDetails `json:"details"`
}
type ProductDetails struct {
	gorm.Model
	ProductsID      uint   `gorm:"not null; default 0; type:int(10)" json:"products_id"`           // productID,商品表，相当于外键
	Label           int    `gorm:"not null; default 0; type:int(5)" json:"label"`                  // 商品标签,1:显示 限时，热卖，按钮,2:不显示
	CoinType        string `gorm:"not null; default ''; type:varchar(10)" json:"coin_type"`        // 币种 BTC，ETH，MGD 等等
	RoiSt           int    `gorm:"not null; default 0; type:int(10)" json:"roi_st"`                // 投资最低回报率(理论值)
	RoiEnd          int    `gorm:"not null; default 0; type:int(10)" json:"roi_end"`               // 投资最高回报率
	Mining          int    `gorm:"not null; default 0; type:int(10)" json:"mining"`                // 算力大小
	MiningUnit      string `gorm:"not null; default ''; type:varchar(10)" json:"mining_unit"`      // 算力单位
	DesirableOutput string `gorm:"not null; default ''; type:varchar(32)" json:"desirable_output"` // 期望产出
	Power           string `gorm:"not null; default ''; type:varchar(20)" json:"power"`            // 功耗
	PowerPrice      string `gorm:"not null; default ''; type:varchar(20)" json:"power_price"`      // 电价
	OccupyPrice     string `gorm:"not null; default ''; type:varchar(20)" json:"occupy_price"`     // 占位费用
	MangerPrice     string `gorm:"not null; default ''; type:varchar(20)" json:"manger_price"`     // 管理费用
	LeaseTime       string `gorm:"not null; default ''; type:varchar(20)" json:"lease_time"`       // 租赁时长
	StartTime       string `gorm:"not null; default ''; type:varchar(20)" json:"start_time"`       // 上架日期，1天内，或者5天
}

// farm server info
type FarmServer struct {
	gorm.Model
	FarmID         string           `gorm:"not null; default ''; type:varchar(64)" json:"farm_id"`
	MinerType      string           `gorm:"not null; default ''; type:varchar(32)" json:"miner_type"`
	PriceList      []MinerPriceList `gorm:"ForeignKey:FarmServerID" json:"price_list"`
	AvailableCount uint64           `gorm:"not null; default 0; type:BIGINT(20)"  json:"available_count"`
	CreateAt       int64            `gorm:"not null; default 0; type:BIGINT(20)" json:"create_at"`
}
type MinerPriceList struct {
	FarmServerID uint   `gorm:"not null; default 0; type:BIGINT(20)" json:"farm_server_id"`
	Value        uint64 `gorm:"not null; default 0; type:BIGINT(20)" json:"value"` // 价格
	Time         uint64 `gorm:"not null; default 0; type:BIGINT(20)" json:"price"` // 时长
}

// 商品列表详情查询
type ProductsQueryJson struct {
	GoodsType int `json:"goodsType"` // 商品类型，1:云算力，2:矿机
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}
type ProductsQueryRes struct {
	ID            uint           `json:"id"`
	GoodsType     int            `json:"goods_type" validate:"required,number,min=1,max=2"` // 商品类型，1:云算力，2:矿机
	GoodsId       string         `json:"goods_id" validate:"required"`                      // 商品id
	GoodsName     string         `json:"goods_name" validate:"required"`                    // 商品名称
	OrgPrice      string         `json:"org_price" validate:"required"`                     // 原始价格
	CurPrice      string         `json:"cur_price" validate:"required"`                     // 当前价格
	Quantity      int            `json:"quantity" validate:"required"`                      // 商品数量
	TotalQuantity int            `json:"total_quantity" validate:"required"`                // 商品总数量
	Unit          string         `json:"unit" validate:"required"`                          // 商品单位
	Description   string         `json:"description" validate:"required"`                   // 商品介绍
	ImageUri      string         `json:"image_uri"`                                         // 商品图片路径
	PushedAt      int64          `json:"pushed_at"`                                         // 发布时间
	Status        bool           `json:"status"`                                            // 当前商品是否已经发布
	Details       ProductDetails `json:"details"`                                           // 商品详情
}

type AProductsJson struct {
	GoodsType int `json:"goodsType"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}

// 更新商品
type AUpdateProductsJson struct {
	AuthorId       int                 `json:"authorId"`                                         // 发布者id
	GoodsId        string              `json:"goodsId"`                                          // 商品id
	GoodsType      int                 `json:"goodsType" validate:"required,number,min=1,max=2"` // 商品类型，1:云算力，2:矿机
	GoodsName      string              `json:"goodsName" validate:"required"`                    // 商品名称
	OrgPrice       string              `json:"orgPrice" validate:"required"`                     // 原始价格
	CurPrice       string              `json:"curPrice" validate:"required"`                     // 当前价格
	Quantity       int                 `json:"quantity" validate:"required"`                     // 可用商品数量
	TotalQuantity  int                 `json:"totalQuantity" validate:"required"`                // 商品总数量
	Unit           string              `json:"unit" validate:"required"`                         // 商品单位
	Description    string              `json:"description" validate:"required"`                  // 商品介绍
	ImageUri       string              `json:"imageUri"`                                         // 商品图片路径
	PushedAt       int64               `json:"pushedAt"`                                         // 发布时间
	FarmID         string              `json:"farmId"`                                           // 矿场ID标识
	MinerGoodsType string              `json:"minerGoodsType"`
	Status         bool                `json:"status"` // 当前商品是否已经发布
	Details        ProductsDetailsJson `json:"details"`
}
type ProductsDetailsJson struct {
	ProductsID      uint   `json:"productsId"`      // productID,商品表，相当于外键
	Label           int    `json:"label"`           // 商品标签，是否显示 限时，热卖，按钮
	CoinType        string `json:"coinType"`        // 币种 BTC，ETH，MGD 等等
	RoiSt           int    `json:"roiSt"`           // 投资最低回报率(理论值)
	RoiEnd          int    `json:"roiEnd"`          // 投资最高回报率
	Mining          int    `json:"mining"`          // 算力大小
	MiningUnit      string `json:"miningUnit"`      // 算力单位
	DesirableOutput string `json:"desirableOutput"` // 期望产出
	Power           string `json:"power"`           // 功耗
	PowerPrice      string `json:"powerPrice"`      // 电价
	OccupyPrice     string `json:"occupyPrice"`     // 占位费用
	MangerPrice     string `json:"mangerPrice"`     // 管理费用
	LeaseTime       string `json:"leaseTime"`       // 租赁时长
	StartTime       string `json:"startTime"`       // 上架日期，1天内，或者5天
}

// 更新商品 状态
type AUpdateProductStatusJson struct {
	AuthorId int    `json:"authorId"`                    // 发布者id
	GoodsId  string `json:"goodsId" validate:"required"` // 商品id
	Status   bool   `json:"status"`                      // 当前商品是否已经发布
}

// farm type
type AFarmServerJson struct {
	FarmID string `json:"farmId"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}
