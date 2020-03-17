package admin

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

// GetArticleList 获取文章列表
func GetArticleList(articleJson *model.ArticleJson) (model.MyMap, string, int) {

	var articleInfoRes []model.ArticleInfoRes
	// 单个查询，并扫描到另外一个表中
	count := 0
	dbs := db.DB
	if articleJson.ArticleId != 0 {
		dbs = dbs.Where("article_infos.article_id = ?", articleJson.ArticleId)
	}
	if articleJson.IsPushed == true {
		dbs = dbs.Where("article_infos.status = ?", articleJson.Status)
	}
	offset, err := db.GetOffset(articleJson.Page, articleJson.Limit)
	if err != nil {
		return nil, "failed", model.STATUS_FAILED
	}

	if err = dbs.Table("article_infos").Where("type = ?", articleJson.Type).
		Count(&count).Offset(offset).Limit(articleJson.Limit).
		Select("article_infos.id,article_infos.article_id,article_infos.category,article_infos.author_id,article_infos.priority,article_infos.status,article_infos.read_count,articles.pushed_at,articles.title").
		Joins("left join articles on articles.article_id = article_infos.article_id").
		Scan(&articleInfoRes).Error; err != nil {
		logrus.Errorf("GetArticleFromDB ArticleInfo err: %s\n", err.Error())
		return nil, "error", model.STATUS_SUCCESS
	}
	logrus.Debugf("GetArticleListFromDB article_id: %d, category: %d, status: %t, ispushed: %t type: %d page: %d limit: %d offset: %d\n",
		articleJson.ArticleId, articleJson.Category, articleJson.Status, articleJson.IsPushed, articleJson.Type, articleJson.Page, articleJson.Limit, offset)
	response := model.MyMap{
		"data":   articleInfoRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminArticleDetail 获取文章详情
func GetAdminArticleDetail(articleJson *model.ArticleDetailJson) (*model.ArticleUpdateInfo, string, int) {

	if err := CheckeArticle(articleJson.ArticleId); err != nil {
		logrus.Errorf("GetArticleDetailFromDB Check Articles err: %s\n", err.Error())
		return nil, err.Error(), model.STATUS_FAILED
	}

	var articleDetail model.ArticleUpdateInfo

	if err := db.DB.Table("article_infos").Where("article_infos.article_id = ?", articleJson.ArticleId).Limit(1).
		Select("article_infos.article_id,article_infos.category,article_infos.author_id,article_infos.priority,article_infos.status,article_infos.read_count,articles.language,articles.pushed_at,articles.announcer,articles.title,articles.summary,articles.content").
		Joins("left join articles on articles.article_id = article_infos.article_id").
		Scan(&articleDetail).Error; err != nil {
		logrus.Errorf("GetAdminArticleDetailFromDB ArticleInfo err: %s\n", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	return &articleDetail, "success", model.STATUS_SUCCESS
}

// UpdateAdminArticleContent 更新文章内容
func UpdateAdminArticleContent(ctx iris.Context, articleJson *model.ArticleUpdateInfo) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	// 更新文章管理表
	articleInfo := &model.ArticleInfo{
		AuthorId:  userId,
		Category:  articleJson.Category, // 文章类型
		Priority:  articleJson.Priority,
		ReadCount: articleJson.ReadCount,
	}
	if err := UpdateArticleInfoToDB(articleJson.ArticleId, articleJson.Type, articleInfo); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	// 更新文章内容表
	articles := &model.Articles{
		Language:  articleJson.Language,
		PushedAt:  articleJson.PushedAt,
		Announcer: articleJson.Announcer,
		Title:     articleJson.Title,
		Summary:   articleJson.Summary,
		Content:   articleJson.Content,
	}
	if err := UpdateArticlesToDB(articleJson.ArticleId, articleJson.Language, articles); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// AddAdminArticleContent 新增文章
func AddAdminArticleContent(ctx iris.Context, articleJson *model.ArticleUpdateInfo) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	// 文章管理表
	articleId := GenNewArticleID()
	articleInfo := &model.ArticleInfo{
		ArticleId: articleId,
		AuthorId:  userId,
		Type:      articleJson.Type,     // 文章分类
		Category:  articleJson.Category, // 文章类型
		Priority:  articleJson.Priority,
		ReadCount: 1,
	}
	if err := db.DB.Create(articleInfo).Error; err != nil {
		logrus.Errorf("AddArticleInfoToDB service.DB.Create err %s", err)
		return "error", model.STATUS_FAILED
	}
	pushedAt := articleJson.PushedAt
	if pushedAt == 0 {
		pushedAt = utils.GetNowTime()
	}
	// 文章内容表
	articles := &model.Articles{
		ArticleId: articleId,
		Language:  articleJson.Language,
		PushedAt:  pushedAt,
		Announcer: articleJson.Announcer,
		Title:     articleJson.Title,
		Summary:   articleJson.Summary,
		Content:   articleJson.Content,
	}
	if err := db.DB.Create(articles).Error; err != nil {
		logrus.Errorf("AddArticlesToDB service.DB.Create err %s", err)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminArticleStatus 更新文章状态
func UpdateAdminArticleStatus(articleStatus *model.ArticleStatusJson) (string, int) {
	articles := new(model.ArticleInfo)
	if err := db.DB.Model(articles).Where("article_id = ? ", articleStatus.ArticleId).
		Update("status", articleStatus.Status).Error; err != nil {
		logrus.Errorf("UpdateAdminArticleStatusToDB failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	logrus.Debugf("UpdateAdminArticleStatusToDB article_id: %d status: %t\n", articleStatus.ArticleId, articleStatus.Status)
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminArticleReadCount 更新文章阅读统计
func UpdateAdminArticleReadCount(articleId int) (string, int) {
	articles := new(model.ArticleInfo)
	if err := db.DB.Model(articles).Where("article_id = ? ", articleId).
		UpdateColumn("read_count", gorm.Expr("read_count + ?", 1)).Error; err != nil {
		logrus.Errorf("UpdateAdminArticleReadCountToDB failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// GetAdminHelpType 获取帮助分类列表
func GetAdminHelpType() (model.MyMap, string, int) {
	var helpArticleCategory []model.ArticleCategory
	count := 0
	if err := db.DB.Model(model.ArticleCategory{}).Count(&count).
		Find(&helpArticleCategory).Error; err != nil {
		logrus.Errorf("GetAdminHelpType failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	// count = len(helpArticleCategory)
	response := model.MyMap{
		"data":   helpArticleCategory,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// AddAdminHelpType 添加帮助中心分类
func AddAdminHelpType(ctx iris.Context, articleJson *model.ArticleCategoryUpdateJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	// 帮助中心
	articleInfo := &model.ArticleCategory{
		AuthorId: userId,
		Type:     articleJson.Type,
		Category: articleJson.Category,
		Name:     articleJson.Name,
		Language: articleJson.Language,
		Status:   false,
	}

	if err := db.DB.Create(articleInfo).Error; err != nil {
		logrus.Errorf("AddAdminHelpType err %s", err)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminHelpType 更新帮助中心分类
func UpdateAdminHelpType(ctx iris.Context, articleJson *model.ArticleCategoryUpdateJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	// 帮助中心
	articleInfo := map[string]interface{}{ // 如果使用结构体更新，对于默认字段，如 0 ,"", false 等将不会更新
		"author_id": userId,
		"type":      articleJson.Type,
		"category":  articleJson.Category,
		"name":      articleJson.Name,
		"language":  articleJson.Language,
		"status":    false,
	}
	if err := db.DB.Model(&model.ArticleCategory{}).Where("id = ?", articleJson.ID).Updates(articleInfo).Error; err != nil {
		logrus.Errorf("AddAdminHelpType err %s", err)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateStatusAdminHelpType 更新帮助中心文章分类
func UpdateStatusAdminHelpType(ctx iris.Context, articleJson *model.ArticleCategoryUpdateJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	// 帮助中心
	articleInfo := map[string]interface{}{ // 如果使用结构体更新，对于默认字段，如 0 ,"", false 等将不会更新
		"author_id": userId,
		"status":    articleJson.Status,
	}
	if err := db.DB.Model(&model.ArticleCategory{}).Where("id = ?", articleJson.ID).Updates(articleInfo).Error; err != nil {
		logrus.Errorf("UpdateStatusAdminHelpType err %s", err)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 判断文章是否存在
func CheckeArticle(articleId int) error {
	articleInfo := new(model.ArticleInfo)
	if err := db.DB.Where("article_id = ?", articleId).First(&articleInfo).RecordNotFound(); err == true {
		return errors.New("文章不存在")
	}
	return nil
}

// 更新文章基础信息
func UpdateArticleInfoToDB(articleId, articleType int, articleInfo *model.ArticleInfo) error {
	articles := new(model.ArticleInfo)
	if err := db.DB.Model(articles).Where("article_id = ? AND type = ? ", articleId, articleType).
		Updates(articleInfo).Update("status", false).Error; err != nil {
		logrus.Errorf("UpdateArticleInfo failed err: %s", err.Error())
		return err
	}
	return nil
}

// 更新文章内容
func UpdateArticlesToDB(articleId, language int, articleJson *model.Articles) error {
	articles := new(model.Articles)
	if err := db.DB.Model(articles).Where("article_id = ? AND language = ? ", articleId, language).
		Updates(articleJson).Error; err != nil {
		logrus.Errorf("UpdateArticlesToDB failed err: %s", err.Error())
		return err
	}
	return nil
}

// 获取最后一篇幅文章ID
func GetArticleLast() int {
	article := new(model.ArticleInfo)
	if err := db.DB.Last(article).Error; err != nil {
		logrus.Errorf("GetUserLast failed err: %s", err.Error())
		return 0
	}
	id := int(article.ID)
	return id
}

// 获取一个文章ID
func GenNewArticleID() int {
	id := GetArticleLast() + 1000
	return id
}
