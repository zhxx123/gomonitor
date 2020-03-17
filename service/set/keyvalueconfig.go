package set

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/sirupsen/logrus"
)

// SetKeyValue 设置key, value
func SetKeyValue(ctx iris.Context, reqData *model.KeyValueData) (map[string]interface{}, string, int) {

	var keyVauleConfig model.KeyValueConfig
	if err := db.DB.Where("key_name = ?", reqData.KeyName).Find(&keyVauleConfig).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logrus.Error(err.Error())
			return nil, "error", model.ErrorCode.ERROR
		}
		var theKeyVauleConfig model.KeyValueConfig
		theKeyVauleConfig.KeyName = reqData.KeyName
		theKeyVauleConfig.Value = reqData.Value
		if err := db.DB.Create(&theKeyVauleConfig).Error; err != nil {
			logrus.Error(err.Error())
			return nil, "error", model.ErrorCode.ERROR
		}
		res := map[string]interface{}{
			"id": theKeyVauleConfig.ID,
		}
		return res, "success", model.ErrorCode.SUCCESS
	}
	keyVauleConfig.Value = reqData.Value
	if err := db.DB.Save(&keyVauleConfig).Error; err != nil {
		logrus.Error(err.Error())
		return nil, "error", model.ErrorCode.ERROR
	}
	res := map[string]interface{}{
		"id": keyVauleConfig.ID,
	}
	return res, "success", model.ErrorCode.SUCCESS
}
