package admin

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/service/user"
	"github.com/zhxx123/gomonitor/utils"

	"github.com/zhxx123/gomonitor/model"
	"github.com/sirupsen/logrus"
)

// CheckAndGetUserList 获取所有用户列表
func CheckAndGetUserList(aul *model.AUserJson) (model.MyMap, string, int) {

	var userList []model.User
	var userListRes []model.AUserRes
	count := 0
	dbs := db.DB
	if len(aul.Email) > 0 {
		dbs = dbs.Where("email = ?", aul.Email)
	}
	if aul.RoleType != 0 {
		dbs = dbs.Where("role_id = ?", aul.RoleType)
	}
	offset, err := db.GetOffset(aul.Page, aul.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if err := dbs.Model(model.User{}).Count(&count).
		Offset(offset).Limit(aul.Limit).
		Find(&userList).Scan(&userListRes).Error; err != nil {
		logrus.Errorf("GetAdminUserListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetUserMessageListFromDB userEmail: %s page: %d limit %d offset: %d\n", aul.Email, aul.Page, aul.Limit, offset)
	response := model.MyMap{
		"data":   userListRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetAdminUserInfo 获取单个用户详细信息
func GetAdminUserInfo(ctx iris.Context) (*model.AUserInfoRes, string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return nil, "error", model.STATUS_FAILED
	}
	return GetUserInfoWithID(userId)
}

// GetAdminUserInfomation 获取指定用户的详情信息
func GetAdminUserInfomation(ctx iris.Context) (*model.AUserInfoRes, string, int) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	return GetUserInfoWithID(id)
}

// UpdateAdminUserStatusRole 更新用户状态角色信息
func UpdateAdminUserStatusRole(aul *model.AUserUpdateStatusJson) (string, int) {
	user := new(model.User)
	if aul.RoleID == 0 { // 更新用户状态
		if err := db.DB.Model(user).Where("user_id = ?", aul.UserId).
			Update("user_status", aul.UserStatus).Error; err != nil {
			logrus.Errorf("UpdateAdminUseStatusToDB failed err: %s", err.Error())
			return "error", model.STATUS_FAILED
		}
		return "success", model.STATUS_SUCCESS
	}
	if err := db.DB.Model(user).Where("user_id = ?", aul.UserId).
		Update("role_id", aul.RoleID).Error; err != nil {
		logrus.Errorf("UpdateAdminUseRoleToDB failed err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	logrus.Debugf("UpdateAdminUserRoleToDB %d %d\n", aul.UserId, aul.RoleID)
	return "success", model.STATUS_SUCCESS
}

// CheckAndDeleteUser 删除用户
func CheckAndDeleteUser(ctx iris.Context) (string, int) {
	userId, err := ctx.Params().GetInt("id")
	if err != nil {
		return "error", model.STATUS_FAILED
	}

	user := new(model.User)
	if err := db.DB.Where("user_id = ?", userId).Delete(user).Error; err != nil {
		logrus.Errorf("DeleteUserById err: %s user_id: %d\n", err.Error(), userId)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// CheckAndLogout 退出登录
func CheckAndLogout(ctx iris.Context) (string, int) {
	userId, err := ctx.Params().GetInt("id")
	if err != nil {
		return "error", model.STATUS_FAILED
	}
	return UserLogOut(userId)
}

// GetAdminUserLoginLogsList 获取用户登录记录
func GetAdminUserLoginLogsList(loginLogJson *model.ALoginLogsJson) (model.MyMap, string, int) {
	var userOauthList []model.UserOauth
	var userOauthListRes []model.AUserOauthRes
	count := 0
	dbs := db.DB
	if loginLogJson.UserId != 0 {
		dbs = dbs.Where("user_id = ?", loginLogJson.UserId)
	}
	if loginLogJson.IsRevoked != false {
		dbs = dbs.Where("revoked = ?", loginLogJson.Revoked)
	}
	offset, err := db.GetOffset(loginLogJson.Page, loginLogJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := dbs.Model(model.UserOauth{}).Order("id desc").Count(&count).
		Offset(offset).Limit(loginLogJson.Limit).
		Find(&userOauthList).Scan(&userOauthListRes).Error; err != nil {
		logrus.Errorf("GetAdminLoginLogsListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminLoginLogsListFromDB loginLogJson: %d %t %t count: %d page: %d limit: %d offset: %d\n", loginLogJson.UserId, loginLogJson.Revoked, loginLogJson.IsRevoked, count, loginLogJson.Page, loginLogJson.Limit, offset)
	response := model.MyMap{
		"data":   userOauthListRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// AdminUserLogout 管理员退出登录
func AdminUserLogout(ctx iris.Context) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}
	var ot model.UserOauth
	if err := db.DB.Model(&ot).Where("user_id = ? AND revoked = ?", userId, false).Update("revoked", true).Error; err != nil {
		logrus.Errorf("UpdateUserOauthByUserId falied Err:%s", err.Error())
		return "error", model.STATUS_FAILED
	}

	// new 未完待更新
	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("DEL", fmt.Sprintf("%s%d", model.LoginUser, userId)); err != nil {
		logrus.Errorf("redis delelte failed: %s", err.Error())
	}
	return "success", model.STATUS_SUCCESS
}

// 获取用户消息列表
func GetAdminUserMessages(message *model.AUserMessageJson) (model.MyMap, string, int) {
	var messsageList []model.UserMessage
	count := 0
	dbs := db.DB
	if message.UserId != 0 {
		dbs = dbs.Where("user_id = ?", message.UserId)
	}
	if message.IsPushed == true {
		dbs = dbs.Where("status = ?", message.Status)
	}
	offset, err := db.GetOffset(message.Page, message.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := dbs.Model(model.UserMessage{}).Count(&count).
		Offset(offset).Limit(message.Limit).
		Find(&messsageList).Error; err != nil {
		logrus.Errorf("GetAdminUserMessageListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetAdminUserMessageListFromDB useId: %d page: %d limit %d offset: %d\n", message.UserId, message.Page, message.Limit, offset)
	response := model.MyMap{
		"data":   messsageList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// AddAdminUserMessages 添加用户消息
func AddAdminUserMessages(ctx iris.Context, message *model.AEditUserMessageJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	if message.UserId == 0 || len(message.Title) == 0 {
		return "faild add", model.STATUS_FAILED
	}
	articleID := GetLastAdminUserMessageRecord()
	pushedAt := message.PushedAt
	if pushedAt == 0 {
		pushedAt = utils.GetNowTime()
	}
	userMessage := &model.UserMessage{
		UserId:    message.UserId,
		ArticleId: int(articleID) + 1000,
		Category:  message.Category,
		Title:     message.Title,
		Summary:   message.Summary,
		Content:   message.Content,
		PushedAt:  pushedAt,
		Announcer: message.Announcer,
		AuthorId:  userId,
	}
	if err := db.DB.Create(userMessage).Error; err != nil {
		logrus.Errorf("AddAdminUserMessageToDB Create UserMessage err %s", err)
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminUserMessageStatus 更新用户消息状态
func UpdateAdminUserMessageStatus(msgDetail *model.AUpdateUserMessageJson) (string, int) {
	message := new(model.UserMessage)
	if err := db.DB.Model(message).Where("user_id = ? AND article_id = ?", msgDetail.UserId, msgDetail.ArticleId).
		Update("status", msgDetail.Status).Error; err != nil {
		logrus.Errorf("UpdateAdminUserMessageToDB failed err: %s", err.Error())
		return "failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// UpdateAdminUserMessage 更新用户消息
func UpdateAdminUserMessage(ctx iris.Context, msgDetail *model.AEditUserMessageJson) (string, int) {
	userInter := ctx.Values().Get("user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)
	if userId == 0 {
		return "error", model.STATUS_FAILED
	}

	if msgDetail.UserId == 0 || msgDetail.ArticleId == 0 || len(msgDetail.Title) == 0 {
		return "failed", model.STATUS_FAILED
	}
	userMsgDetail := map[string]interface{}{
		"title":     msgDetail.Title,
		"summary":   msgDetail.Summary,
		"content":   msgDetail.Content,
		"announcer": msgDetail.Announcer,
		"authorId":  userId,
		"status":    false,
	}
	message := new(model.UserMessage)
	if err := db.DB.Model(message).Where("user_id = ? AND article_id = ?", msgDetail.UserId, msgDetail.ArticleId).Updates(userMsgDetail).Error; err != nil {
		logrus.Errorf("UpdateAdminUserMessageContentToDB failed err: %s", err.Error())
		return "failed", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// GetAdminUserMangerList 管理员列表
func GetAdminUserMangerList(aul *model.AUserJson) (model.MyMap, string, int) {

	var userList []model.User
	var userListRes []model.AUserRes
	count := 0
	dbs := db.DB
	if len(aul.Email) > 0 {
		dbs = dbs.Where("email = ?", aul.Email)
	}
	if aul.RoleType != 0 {
		dbs = dbs.Where("role_id = ?", aul.RoleType)
	} else {
		dbs = dbs.Where("role_id != 1")
	}
	offset, err := db.GetOffset(aul.Page, aul.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	if err := dbs.Model(model.User{}).Count(&count).
		Offset(offset).Limit(aul.Limit).
		Find(&userList).Scan(&userListRes).Error; err != nil {
		logrus.Errorf("GetAdminUserListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_SUCCESS
	}
	logrus.Debugf("GetUserMessageListFromDB userEmail: %s page: %d limit %d offset: %d\n", aul.Email, aul.Page, aul.Limit, offset)
	response := model.MyMap{
		"data":   userListRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 用户注销
func UserLogOut(userId int) (string, int) {
	ot := new(model.UserOauth)
	if err := db.DB.Model(ot).Where("user_id = ? AND revoked = ?", userId, false).Update("revoked", true).Error; err != nil {
		logrus.Errorf("UpdateUserOauthByUserId falied Err: %s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 获取用户详情
func GetUserInfoWithID(userId int) (*model.AUserInfoRes, string, int) {
	user, err := user.GetUserById(userId)
	if err != nil {
		return nil, "no user", model.STATUS_FAILED
	}
	// 查询管理员权限表，获取管理员权限
	//
	var roleList []int
	roleList = append(roleList, user.RoleID)
	userInfo := &model.AUserInfoRes{
		ID:       user.ID,
		UserId:   userId,
		Username: user.Username,
		Email:    user.Email,
		// Password:   user.Password,
		RegisterAt: user.RegisterAt,
		Phone:      user.Phone,
		Roles:      roleList,
		RegisterIp: user.RegisterIp,
		UserStatus: user.UserStatus,
	}
	return userInfo, "success", model.STATUS_SUCCESS
}

// 获取最后一条消息ID
func GetLastAdminUserMessageRecord() uint {
	message := new(model.UserMessage)
	if err := db.DB.Last(message).Error; err != nil {
		logrus.Errorf("GetLastAdminUserMessageRecord failed err: %s", err.Error())
		return 0
	}
	return message.ID
}
