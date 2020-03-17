package model

import "github.com/jinzhu/gorm"

const (
	ARTICLE_TYPE_HELP   = 1
	ARTICLE_TYPE_NOTICE = 2
)

type ArticleInfo struct {
	gorm.Model
	AuthorId  int  `gorm:"not null; default 0; type:int(10)" json:"author_id"`  //发布者ID
	Type      int  `gorm:"not null; default 0; type:int(2)" json:"type"`        //文章类型，1:帮助文章，2:公告文章
	ArticleId int  `gorm:"not null; default 0; type:int(10)" json:"article_id"` //文章唯一标识符号
	Category  int  `gorm:"not null; default 0; type:int(2)" json:"category"`    //文章类别，对于帮助文章1:order2:account,3:safe
	Priority  int  `gorm:"not null; default 0; type:int(2)" json:"priority"`    //文章优先级，用于公告
	ReadCount int  `gorm:"not null; default 0; type:int(10)" json:"read_count"` //文章阅读量
	Status    bool `gorm:"not null; default 0; type:tinyint(1)" json:"status"`  //当前消息是否已经发布
}

type ArticleCategory struct {
	gorm.Model
	AuthorId int    `gorm:"not null; default 0; type:int(10)" json:"author_id"` //发布者ID
	Type     int    `gorm:"not null; default 0; type:int(2)" json:"type"`       //文章类型，1:帮助文章，2:公告文章
	Category int    `gorm:"not null; default 0; type:int(2)" json:"category"`   //文章类别，对于帮助文章1:order2:account,3:safe
	Name     string `gorm:"not null; default ''; type:varchar(20)" json:"name"` //类型名称
	Language int    `gorm:"not null; default 0; type:int(2)" json:"language"`   //语言类型，1:中文，2:英文
	Status   bool   `gorm:"not null; default 0; type:tinyint(1)" json:"status"` //当前条目是否已经发布，可见
}

type Articles struct {
	gorm.Model
	ArticleId int    `gorm:"not null; default 0; type:int(10)" json:"article_id"`     //文章唯一标识符号
	Language  int    `gorm:"not null; default 0; type:int(2)" json:"language"`        //文章语言0:中文1:英文
	PushedAt  int64  `gorm:"not null; default ''; type:BIGINT(20)" json:"pushed_at"`  //发布时间
	Announcer string `gorm:"not null; default ''; type:varchar(20)" json:"announcer"` //发布者标签xxx团队
	Title     string `gorm:"not null; default ''; type:varchar(128)" json:"title"`    //文章标题
	Summary   string `gorm:"not null; default ''; type:varchar(128)" json:"summary"`  //文章摘要
	Content   string `gorm:"not null; default ''; type:text" json:"content"`          //文章内容
}

type ArticleDetailJson struct {
	ArticleId int `json:"articleId"` // 文章唯一标识符号
	Language  int `json:"language"`  // 文章语言
}
type ArticleDetailRes struct {
	ArticleId int    `json:"article_id"` // 文章唯一标识符号
	Language  int    `json:"language"`   // 文章语言 0:中文 1: 英文
	PushedAt  int64  `json:"pushed_at"`  // 发布时间
	Announcer string `json:"announcer"`  // 发布者标签 xxx团队
	Title     string `json:"title"`      // 文章标题
	Summary   string `json:"summary"`    // 文章摘要
	Content   string `json:"content"`    // 文章内容
}
type AritlceUserJson struct {
	Category int `json:"category"`
	Language int `json:"language"`
	Type     int `json:"type"`
	Page     int `json:"page"`
	Limit    int `json:"limit"`
}

// 文章标识
type ArticleJson struct {
	ArticleId int  `json:"articleId"` // 文章唯一标识符号
	Category  int  `json:"category"`
	Status    bool `json:"status"`
	IsPushed  bool `json:"ispushed"`
	Type      int  `json:"type" validate:"required,number,min=0,max=2"`
	Page      int  `json:"page" validate:"required,number,min=1,max=100"`
	Limit     int  `json:"limit" validate:"required,number,min=20,max=20"`
}

// 获取文章列表
type ArticleInfoRes struct {
	ID        uint   `json:"id"`
	ArticleId int    `json:"article_id"` // 文章唯一标识符号
	Category  int    `json:"category"`   //文章类别，对于帮助文章 1:order 2:account, 3:safe
	AuthorId  int    `json:"author_id"`  // 发布者ID
	Priority  int    `json:"priority"`   // 文章优先级，用于公告
	Status    bool   `json:"status"`     // 当前消息是否已经发布
	ReadCount int    `json:"read_count"` // 文章阅读量
	Title     string `json:"title"`      // 文章标题
	PushedAt  int64  `json:"pushed_at"`  // 发布时间
}

// 获取帮助中心文章列表
type ArticleHelpInfoRes struct {
	ID        uint   `json:"id"`
	ArticleId int    `json:"article_id"` // 文章唯一标识符号
	Category  int    `json:"category"`   //文章类别，对于帮助文章 1:order 2:account, 3:safe
	Priority  int    `json:"priority"`   // 文章优先级，用于公告
	ReadCount int    `json:"read_count"` // 文章阅读量
	Title     string `json:"title"`      // 文章标题
	Content   string `json:"content"`    // 文章内容
}

// 新增或者修改文章
type ArticleUpdateInfo struct {
	ArticleId int    `json:"articleId"` // 文章唯一标识符号
	Category  int    `json:"category"`  //文章类别，对于帮助文章 1:order 2:account, 3:safe
	AuthorId  int    `json:"authorId"`  // 发布者ID
	Type      int    `json:"type"`      // 文章类型，1:帮助文章，2:公告文章
	Priority  int    `json:"priority"`  // 文章优先级，用于公告
	Status    bool   `json:"status"`    // 当前消息是否已经发布
	ReadCount int    `json:"readCount"` // 文章阅读量
	Language  int    `json:"language"`  // 文章语言 0:中文 1: 英文
	PushedAt  int64  `json:"pushedAt"`  // 发布时间
	Announcer string `json:"announcer"` // 发布者标签 xxx团队
	Title     string `json:"title"`     // 文章标题
	Summary   string `json:"summary"`   // 文章摘要
	Content   string `json:"content"`   // 文章内容
}

// 修改文章状态
type ArticleStatusJson struct {
	ArticleId int  `json:"articleId"` // 文章唯一标识符号
	Status    bool `json:"status"`
}

// 新增文章
type ArticleCategoryUpdateJson struct {
	ID       uint   `json:"id"`
	Type     int    `json:"type"`     // 文章类型，1:帮助文章，2:公告文章
	Category int    `json:"category"` //文章类别，对于帮助文章 1:order 2:account, 3:safe
	Name     string `json:"name"`     // 类型名称
	Language int    `json:"language"` // 语言类型， 1:中文，2:英文
	Status   bool   `json:"status"`   // 当前条目是否已经发布，可见
}

// 帮助中心类型列表
type ArticleCategoryRes struct {
	Category int    `json:"category"` // 文章类别 1:order 2:account, 3:safe
	Name     string `json:"name"`     // 类型名称
	// Language int    `json:"language"` // 语言类型， 1:中文，2:英文
}
