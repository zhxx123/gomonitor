package controller

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/admin"
)

// GetAdminArticleList 查询文章列表
func GetAdminArticleList(ctx iris.Context) {
	articleJson := new(model.ArticleJson)

	if err := ctx.ReadQuery(articleJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetArticleList(articleJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// QueryAdminArticle 查询单个文章详情
func QueryAdminArticle(ctx iris.Context) {

	articleDetailJson := new(model.ArticleDetailJson)

	if err := ctx.ReadQuery(articleDetailJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleDetailJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	response, msg, status := admin.GetAdminArticleDetail(articleDetailJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// UpdateAdminArticle 更新文章内容
func UpdateAdminArticle(ctx iris.Context) {

	articleUpdateJson := new(model.ArticleUpdateInfo)

	if err := ctx.ReadJSON(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminArticleContent(ctx, articleUpdateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// AddAdminArticle 新增文章
func AddAdminArticle(ctx iris.Context) {
	articleUpdateJson := new(model.ArticleUpdateInfo)

	if err := ctx.ReadJSON(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.AddAdminArticleContent(ctx, articleUpdateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminArticleStatus 更新文章状态
func UpdateAdminArticleStatus(ctx iris.Context) {
	articleStatus := new(model.ArticleStatusJson)

	if err := ctx.ReadJSON(articleStatus); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleStatus); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminArticleStatus(articleStatus)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// UpdateAdminArticleReadCount 更新文章阅读数量
func UpdateAdminArticleReadCount(ctx iris.Context) {
	article_id := ctx.Values().Get("id")
	if article_id == 0 {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_FAILED, nil, "参数错误"))
		return
	}
	articleId := article_id.(int)
	msg, status := admin.UpdateAdminArticleReadCount(articleId)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// GetAdminHelpType 获取帮助中心分类
func GetAdminHelpType(ctx iris.Context) {
	response, msg, status := admin.GetAdminHelpType()
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, response, msg))
}

// AddAdminHelpType 添加帮助文章分类
func AddAdminHelpType(ctx iris.Context) {
	articleUpdateJson := new(model.ArticleCategoryUpdateJson)

	if err := ctx.ReadJSON(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.AddAdminHelpType(ctx, articleUpdateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// AddAdminHelpType 添加帮助文章分类
func UpdateAdminHelpType(ctx iris.Context) {
	articleUpdateJson := new(model.ArticleCategoryUpdateJson)

	if err := ctx.ReadJSON(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateAdminHelpType(ctx, articleUpdateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}

// AddAdminHelpType 添加帮助文章分类
func UpdateStatusAdminHelpType(ctx iris.Context) {
	articleUpdateJson := new(model.ArticleCategoryUpdateJson)

	if err := ctx.ReadJSON(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	if err := validate.Struct(articleUpdateJson); err != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(ApiResource(model.STATUS_PARAM_ERROR, nil, err.Error()))
		return
	}
	msg, status := admin.UpdateStatusAdminHelpType(ctx, articleUpdateJson)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ApiResource(status, nil, msg))
}
