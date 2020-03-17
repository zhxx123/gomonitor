package article

import (
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/sirupsen/logrus"
)

// GetHelpCenterTypeList 获取帮助中心问题类型列表
func GetHelpCenterTypeList() (model.MyMap, string, int) {
	var helpArticleCategory []model.ArticleCategory
	var artilceCategoryList []model.ArticleCategoryRes
	count := 0
	if err := db.DB.Model(model.ArticleCategory{}).Count(&count).
		Find(&helpArticleCategory).Scan(&artilceCategoryList).Error; err != nil {
		logrus.Errorf("GetHelpCenterTypeList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	// count = len(helpArticleCategory)
	response := model.MyMap{
		"data":   helpArticleCategory,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取文章列表
func GetArticleList(articleJson *model.AritlceUserJson) (model.MyMap, string, int) {

	var articleInfoRes []model.ArticleHelpInfoRes
	// 单个查询，并扫描到另外一个表中
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(articleJson.Page, articleJson.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if articleJson.Category != 0 {
		dbs = dbs.Where("category = ?", articleJson.Category)
	}
	if err = dbs.Table("article_infos").Where("article_infos.type = ? AND article_infos.status = ", articleJson.Type, true).
		Count(&count).Offset(offset).Limit(articleJson.Limit).
		Select("article_infos.article_id,article_infos.category,article_infos.priority,articles.title,articles.content").
		Joins("left join articles on articles.article_id = article_infos.article_id").
		Scan(&articleInfoRes).Error; err != nil {
		logrus.Errorf("GetArticleFromDB ArticleInfo err: %s\n", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	response := model.MyMap{
		"data":   articleInfoRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取帮助中心文章列表
func GetArticleHelpList(articleJson *model.AritlceUserJson) (model.MyMap, string, int) {

	var articleInfoRes []model.ArticleInfoRes
	// 单个查询，并扫描到另外一个表中
	count := 0
	dbs := db.DB
	offset, err := db.GetOffset(articleJson.Page, articleJson.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if articleJson.Category != 0 {
		dbs = dbs.Where("category = ?", articleJson.Category)
	}
	if err = dbs.Table("article_infos").Where("article_infos.type = ? AND article_infos.status = ", articleJson.Type, true).
		Count(&count).Offset(offset).Limit(articleJson.Limit).
		Select("article_infos.id,article_infos.article_id,article_infos.category,article_infos.author_id,article_infos.priority,article_infos.read_count,articles.pushed_at,articles.title").
		Joins("left join articles on articles.article_id = article_infos.article_id").
		Scan(&articleInfoRes).Error; err != nil {
		logrus.Errorf("GetArticleFromDB ArticleInfo err: %s\n", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	response := model.MyMap{
		"data":   articleInfoRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

func GetArticleDetail(articleJson *model.ArticleDetailJson) (*model.ArticleDetailRes, string, int) {
	var articleInfo model.ArticleInfo
	if err := db.DB.Where("article_id = ?", articleJson.ArticleId).First(&articleInfo).RecordNotFound(); err == true {
		return nil, "文章不存在", model.STATUS_FAILED
	}
	var articleDetail model.ArticleDetailRes
	var articles model.Articles
	if err := db.DB.Where("article_id = ? AND language = ? ", articleJson.ArticleId, articleJson.Language).First(&articles).Scan(&articleDetail).Error; err != nil {
		logrus.Errorf("GetArticleDetailFromDB Articles err: %s\n", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	return &articleDetail, "success", model.STATUS_SUCCESS
}
