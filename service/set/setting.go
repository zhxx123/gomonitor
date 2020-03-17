package set

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/sirupsen/logrus"
)

var (
	rwL     = sync.RWMutex{}
	confMap map[string]interface{}
)

func InitSysSeting() {
	setList, err := GetAllSettingFromDB()
	if err != nil {
		logrus.Errorf("InitSysSeting GetAllSettingFromDB error %s", err.Error())
	}
	if len(setList) == 0 {
		logrus.Infof("InitSysSeting GetAllSettingFromDB empty")
		return
	}
	rwL.Lock()
	confMap = make(map[string]interface{})
	for _, value := range setList {
		var mapResult map[string]interface{}
		jsonStr := value.Value
		err := json.Unmarshal([]byte(jsonStr), &mapResult)
		if err != nil {
			logrus.Errorf("InitSysSeting json.Unmarshal err: %s", err.Error())
		}
		for k, v := range mapResult {
			newKey := fmt.Sprintf("%s_%s_%s", value.Category, value.Name, k)
			confMap[newKey] = v
		}
	}
	rwL.Unlock()
}

func GetConfMapKeyValue(key string, defValue interface{}) (interface{}, bool) {
	if key == "" {
		return defValue, false
	}
	rwL.RLock()
	defer rwL.RUnlock()
	for k, v := range confMap {
		if k == key {
			return v, true
		}
	}
	return defValue, false
}

// 更新安全设置
func UpdateUserSafeSeting() bool {
	verifycode, ok := GetConfMapKeyValue("sys_safe_verifycode", false)
	if ok && reflect.TypeOf(verifycode).String() == "bool" {
		model.CHECK_VERIFY_CODE = verifycode.(bool)
	}
	pausetrade, ok := GetConfMapKeyValue("sys_safe_verifycode", false)
	if ok && reflect.TypeOf(pausetrade).String() == "bool" {
		model.PAUSE_TRADE = pausetrade.(bool)
	}
	logrus.Infof("UpdateUserSafeSeting verifycode: %t, pausetrade %t", model.CHECK_VERIFY_CODE, model.PAUSE_TRADE)
	return true
}

// 获取系统设置记录
func GetAdminSettingsFromDB(settingJson *model.ASettingsJson) (*model.Settings, error) {
	settings := new(model.Settings)
	if err := db.DB.Where("category = ? AND name = ?", settingJson.Category, settingJson.Name).First(settings).Error; err != nil {
		logrus.Errorf("GetAdminSettingsFromDB failed err: %s", err.Error())
		return nil, err
	}
	return settings, nil
}

// 获取所有系统设置表
func GetAllSettingFromDB() ([]model.Settings, error) {
	settings := new(model.Settings)
	var setingList []model.Settings
	if err := db.DB.Model(settings).Find(&setingList).Error; err != nil {
		logrus.Errorf("GetAdminSettingsFromDB failed err: %s", err.Error())
		return nil, err
	}
	return setingList, nil
}
