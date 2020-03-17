package product

import (
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/sirupsen/logrus"
)

// 获取商品详情列表
func GetProductsList(productJson *model.ProductsQueryJson) (model.MyMap, string, int) {
	var productList []model.Products
	var productListRes []model.ProductsQueryRes
	count := 0
	dbs := db.DB
	if productJson.GoodsType == 0 {
		return nil, "参数错误", model.STATUS_FAILED
	}
	offset, err := db.GetOffset(productJson.Page, productJson.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	var products model.Products
	if err := dbs.Model(&products).
		Where("goods_type = ? AND status = ?", productJson.GoodsType, true).
		Count(&count).Offset(offset).Limit(productJson.Limit).
		Find(&productList).Scan(&productListRes).Error; err != nil {
		logrus.Errorf("GetProductListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}

	lens := len(productList)
	for i := 0; i < lens; i++ {
		if err := db.DB.Model(&productListRes[i]).Related(&productListRes[i].Details, "product_details", "products_id").Error; err != nil {
			logrus.Error(err.Error())
			// return nil, "error", model.STATUS_FAILED
		}
	}

	count = len(productListRes)
	logrus.Debugf("GetProductListFromDB  goods_type: %d count: %d page: %d limit: %d offset: %d\n", productJson.GoodsType, count, productJson.Page, productJson.Limit, offset)

	response := model.MyMap{
		"data":   productListRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}
