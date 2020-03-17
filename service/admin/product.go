package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
	"github.com/sirupsen/logrus"
)

// GetAdminProductList 获取商品列表
func GetAdminProductList(productJson *model.AProductsJson) (model.MyMap, string, int) {
	var productList []model.Products
	count := 0
	dbs := db.DB
	if productJson.GoodsType != 0 {
		dbs = dbs.Where("goods_type = ?", productJson.GoodsType)
	}
	offset, err := db.GetOffset(productJson.Page, productJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	products := new(model.Products)
	if err := dbs.Model(products).Count(&count).
		Offset(offset).Limit(productJson.Limit).
		Find(&productList).Error; err != nil {
		logrus.Errorf("GetAdminProductListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	lens := len(productList)
	for i := 0; i < lens; i++ {
		if err := db.DB.Model(&productList[i]).Related(&productList[i].Details, "product_details").Error; err != nil {
			logrus.Error(err.Error())
			// return nil, "error", model.STATUS_FAILED
		}
	}
	logrus.Debugf("GetAdminProductList count: %d page: %d limit: %d offset: %d\n", count, productJson.Page, productJson.Limit, offset)

	response := model.MyMap{
		"data":   productList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// CreateAdminProducts 新增商品
func CreateAdminProducts(ctx iris.Context, productJson *model.AUpdateProductsJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	status := GetFarmServersFromDB(productJson.FarmID, productJson.MinerGoodsType)
	if status != true {
		return "Farm 不存在", model.STATUS_FAILED
	}

	goodsId := utils.GenGoodsID(productJson.GoodsType)
	// 添加更新时间
	pubshedAt := utils.GetNowTime()
	products := &model.Products{
		AuthorId:       userId,
		GoodsType:      productJson.GoodsType,
		GoodsId:        goodsId,
		GoodsName:      productJson.GoodsName,
		OrgPrice:       productJson.OrgPrice,
		CurPrice:       productJson.CurPrice,
		Quantity:       productJson.Quantity,      // 可以用数量
		TotalQuantity:  productJson.TotalQuantity, // 总数量
		Unit:           productJson.Unit,          // 商品单位
		Description:    productJson.Description,
		ImageUri:       productJson.ImageUri, // btc图片的地址
		PushedAt:       pubshedAt,
		FarmID:         productJson.FarmID,
		MinerGoodsType: productJson.MinerGoodsType,
		Status:         false,
	}
	// 添加商品表
	if err := db.DB.Create(products).Error; err != nil {
		logrus.Errorf("UserOauthCreate failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	details := productJson.Details
	// 添加商品详情表
	productDetails := &model.ProductDetails{
		ProductsID:      products.ID,
		Label:           details.Label,
		CoinType:        details.CoinType,
		RoiSt:           details.RoiSt,
		RoiEnd:          details.RoiEnd,
		Mining:          details.Mining,
		MiningUnit:      details.MiningUnit,
		DesirableOutput: details.DesirableOutput,
		Power:           details.Power,
		PowerPrice:      details.PowerPrice,
		OccupyPrice:     details.OccupyPrice,
		MangerPrice:     details.MangerPrice,
		LeaseTime:       details.LeaseTime,
		StartTime:       details.StartTime,
	}
	if err := db.DB.Create(productDetails).Error; err != nil {
		logrus.Errorf("CreateAdminProducts productDetails failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminProducts 更新商品信息
func UpdateAdminProducts(ctx iris.Context, productJson *model.AUpdateProductsJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	status := GetFarmServersFromDB(productJson.FarmID, productJson.MinerGoodsType)
	if status != true {
		return "Farm 不存在", model.STATUS_FAILED
	}
	productJson.AuthorId = userId
	productJson.Status = false //修改商品之后，用户不能查看
	// 添加更新时间
	productJson.PushedAt = utils.GetNowTime()
	product := new(model.Products)
	if err := db.DB.Model(product).Where("goods_id = ? AND goods_type = ?", productJson.GoodsId, productJson.GoodsType).
		Updates(productJson).Update("status", false).Error; err != nil {
		logrus.Errorf("UpdateAdminProducts failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	// 更新商品详情
	details := productJson.Details
	var productDetails model.ProductDetails
	if err := db.DB.Model(&productDetails).Where("products_id = ?", details.ProductsID).Updates(details).Error; err != nil {
		logrus.Errorf("UpdateAdminProducts failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminProductStatus 更新商品状态
func UpdateAdminProductStatus(ctx iris.Context, productJson *model.AUpdateProductStatusJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	updateData := map[string]interface{}{
		"author_id": productJson.AuthorId,
		"status":    productJson.Status,
	}
	product := new(model.Products)
	if err := db.DB.Model(product).Where("goods_id = ?", productJson.GoodsId).
		Updates(updateData).Error; err != nil {
		logrus.Errorf("UpdateAdminProductStatusToDB failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// GetAdminProductDetail 获取商品详情
func GetAdminProductDetail(goodsId string) (*model.Products, string, int) {
	var product model.Products
	if err := db.DB.Where("goods_id = ?", goodsId).First(&product).Error; err != nil {
		logrus.Errorf("UpdateAdminProductToDB failed err: %s", err.Error())
		return nil, err.Error(), model.STATUS_FAILED
	}
	return &product, "success", model.STATUS_SUCCESS
}

// GetAdminFarmsList 获取矿场机器列表
func GetAdminFarmsList(productJson *model.AFarmServerJson) (model.MyMap, string, int) {
	var farmServerList []model.FarmServer
	count := 0
	dbs := db.DB
	if productJson.FarmID != "" {
		dbs = dbs.Where("farm_id = ?", productJson.FarmID)
	}
	offset, err := db.GetOffset(productJson.Page, productJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	farms := new(model.FarmServer)
	if err := dbs.Model(farms).Count(&count).
		Offset(offset).Limit(productJson.Limit).
		Find(&farmServerList).Error; err != nil {
		logrus.Errorf("GetAdminFarmsList failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	lens := len(farmServerList)
	for i := 0; i < lens; i++ {
		if err := db.DB.Model(&farmServerList[i]).Related(&farmServerList[i].PriceList).Error; err != nil {
			logrus.Error(err.Error())
			// return nil, "error", model.STATUS_FAILED
		}
	}
	logrus.Debugf("GetAdminFarmsList count: %d page: %d limit: %d offset: %d\n", count, productJson.Page, productJson.Limit, offset)

	response := model.MyMap{
		"data":   farmServerList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取矿场机器详情
func GetFarmServersFromDB(farmID, minerType string) bool {
	farmSever := new(model.FarmServer)
	if err := db.DB.Where("farm_id = ? AND miner_type = ? ", farmID, minerType).First(farmSever).RecordNotFound(); err == true {
		return false
	}
	return true
}
