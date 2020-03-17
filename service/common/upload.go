package common

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/cache"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Upload 文件上传
func UploadImage(ctx iris.Context) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	file, info, err := ctx.FormFile("upFile")

	if err != nil {
		return nil, "参数无效", model.STATUS_FAILED
	}

	var filename = info.Filename
	var index = strings.LastIndex(filename, ".")

	if index < 0 {
		return nil, "无效的文件名", model.STATUS_FAILED
	}

	var ext = filename[index:]
	if len(ext) == 1 {
		return nil, "无效的扩展名", model.STATUS_FAILED
	}
	var mimeType = mime.TypeByExtension(ext)

	logrus.Debugf("filename %s, index %d, ext %s, mimeType %s\n", filename, index, ext, mimeType)
	if mimeType == "" && ext == ".jpeg" {
		mimeType = "image/jpeg"
	}
	if mimeType == "" {
		return nil, "无效的图片类型", model.STATUS_FAILED
	}

	imgUploadedInfo := GenerateImgUploadedInfo(ext, userId)

	logrus.Debug(imgUploadedInfo.UploadDir)

	if err := os.MkdirAll(imgUploadedInfo.UploadDir, 0777); err != nil {
		logrus.Error(err)
		return nil, "error", model.STATUS_FAILED
	}

	if err := saveUploadedFile(file, imgUploadedInfo.UploadFilePath); err != nil {
		logrus.Error(err)
		return nil, "error save", model.STATUS_FAILED
	}

	image := model.Image{
		Title:        imgUploadedInfo.Filename,
		OrignalTitle: filename,
		URL:          imgUploadedInfo.ImgURL,
		Width:        0,
		Height:       0,
		Mime:         mimeType,
	}

	if err := db.DB.Create(&image).Error; err != nil {
		logrus.Error(err)
		return nil, "image error", model.STATUS_FAILED
	}

	res := model.MyMap{
		// "id":       image.ID,
		"url": imgUploadedInfo.ImgURL,
		// "title":    imgUploadedInfo.Filename, //新文件名
		// "original": filename,                 //原始文件名
		// "type":     mimeType,                 //文件类型
	}
	return res, "success", model.STATUS_SUCCESS
}
func GetWorkID(ctx iris.Context, workOrderJson *model.WorkOrderIDJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return nil, "错误", model.STATUS_FAILED
	}
	if workOrderJson.UID == "" {
		return nil, "error", model.STATUS_FAILED
	}

	orderKey := fmt.Sprintf("%s%d%s", model.WorkOrderKey, userId, workOrderJson.UID)
	workOrderID, err := checkOrderAndGetCache(userId, orderKey)
	if err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}
	if len(workOrderID) == 0 {
		workOrderID = utils.GenWorkOrderNo(userId)
		// 将用户提交的订和 uuid 存入缓存
		cache.OC.Set(orderKey, workOrderID, cache.CacheDefaultExpiration*2)
	}
	response := model.MyMap{
		"work_id": workOrderID,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UploadWorkOrder 提交工单并保存
func UploadWorkOrder(ctx iris.Context, workOrderJson *model.WorkOrderJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	orderKey := fmt.Sprintf("%s%d%d", model.WorkOrderKey, userId, workOrderJson.UID)
	value, found := cache.OC.Get(orderKey)
	if found == false || value.(string) != workOrderJson.WorkId {
		return "工单错误,请刷新重试", model.STATUS_FAILED
	}
	// 待定，检查缓存工单状态，只返回工单状态 ！！！
	if _, err := checkWorkOrderCache(workOrderJson.WorkId); err != nil {
		return err.Error(), model.STATUS_FAILED
	}

	nowtime := utils.GetNowTime()
	workOrder := &model.WorkOrder{
		UserId:      userId,
		WorkId:      workOrderJson.WorkId,
		IssueType:   workOrderJson.IssueType,
		Description: workOrderJson.Description,
		Email:       workOrderJson.Email,
		ImgUri:      workOrderJson.ImgUri,
		CreateAt:    nowtime,
		Status:      model.WORKORDER_PROCESS_STATUS,
	}
	if err := db.DB.Create(workOrder).Error; err != nil {
		logrus.Errorf("UploadWorkOrder service.DB.Create err %s, userId: %d, workId: %s", err.Error(), workOrder.UserId, workOrder.WorkId)
		return "error", model.STATUS_FAILED
	}

	// 加入缓存,只缓存工单状态
	cache.OC.Set(workOrderJson.WorkId, model.WORKORDER_PROCESS_STATUS, cache.CacheDefaultExpiration)

	return "success", model.STATUS_SUCCESS
}

// 保存文件
func saveUploadedFile(file multipart.File, fname string) error {
	out, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, file)
	return nil
}

// 创建一个ImageUploadedInfo
func GenerateImgUploadedInfo(ext string, userId int) model.ImageUploadedInfo {
	sep := string(os.PathSeparator)
	uploadImgDir := config.ServerConfig.UploadImgDir
	length := utf8.RuneCountInString(uploadImgDir)
	lastChar := uploadImgDir[length-1:]
	ymStr := utils.GetTodayYM(sep)

	var uploadDir string
	if lastChar != sep {
		uploadDir = uploadImgDir + sep + ymStr
	} else {
		uploadDir = uploadImgDir + ymStr
	}
	uuid := uuid.NewV4()
	uuidName := uuid.String()
	filename := fmt.Sprintf("%s%05d%s", uuidName, userId, ext)
	uploadFilePath := uploadDir + sep + filename
	imgURL := strings.Join([]string{
		"https://" + config.ServerConfig.ImgHost + config.ServerConfig.ImgPath,
		ymStr,
		filename,
	}, "/")
	res := model.ImageUploadedInfo{
		ImgURL:         imgURL,
		UUIDName:       uuidName,
		Filename:       filename,
		UploadDir:      uploadDir,
		UploadFilePath: uploadFilePath,
	}
	return res
}

// 工单缓存
func checkOrderAndGetCache(userId int, orderKey string) (string, error) {
	value, found := cache.OC.Get(orderKey)
	if found { // 当前workId已经存在,可能是用户重复请求
		logrus.Debugf("checkOrderAndGetCache workId: %s found: %t value: %s\n", orderKey, found, value.(string))
		return value.(string), nil
	}
	return "success", nil
}

// 工单提交缓存
func checkWorkOrderCache(outTradeNo string) (int, error) {
	if outTradeNo == "" {
		return 0, errors.New("参数错误")
	}
	value, found := cache.OC.Get(outTradeNo)
	if found {
		if status, ok := (value).(int); ok {
			return status, errors.New("工单已经创建")
		}
		return 0, errors.New("error")
	}
	return 0, nil
}
