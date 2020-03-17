package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zhxx123/gomonitor/service/wallet"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/cache"
	"github.com/zhxx123/gomonitor/service/common"
	"github.com/zhxx123/gomonitor/service/db"
	"github.com/zhxx123/gomonitor/utils"
)

// model.UserFromRedis 从redis中取出用户信息
func UserFromRedis(userID int) (model.User, error) {
	loginUser := fmt.Sprintf("%s%d", model.LoginUser, userID)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	userBytes, err := redis.Bytes(RedisConn.Do("GET", loginUser))
	if err != nil {
		logrus.Error(err)
		return model.User{}, errors.New("未登录")
	}
	var user model.User
	bytesErr := json.Unmarshal(userBytes, &user)
	if bytesErr != nil {
		logrus.Error(bytesErr)
		return user, errors.New("未登录")
	}
	return user, nil
}

// UserToRedis 将用户信息存到redis
func UserToRedis(user model.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		logrus.Error(err)
		return errors.New("error")
	}
	loginUserKey := fmt.Sprintf("%s%d", model.LoginUser, user.ID)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	if _, redisErr := RedisConn.Do("SET", loginUserKey, userBytes, "EX", config.ServerConfig.TokenMaxAge); redisErr != nil {
		logrus.Errorf("redis set failed: %s", redisErr.Error())
		return errors.New("error")
	}
	return nil
}

// OauthTokenFromRedis 从redis中取出用户登录信息
func OauthTokenFromRedis(userID int) (model.UserOauth, error) {
	loginUser := fmt.Sprintf("%s%d", model.LoginUser, userID)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	userBytes, err := redis.Bytes(RedisConn.Do("GET", loginUser))
	if err != nil {
		logrus.Error(err)
		return model.UserOauth{}, errors.New("未登录")
	}
	var user model.UserOauth
	bytesErr := json.Unmarshal(userBytes, &user)
	if bytesErr != nil {
		logrus.Error(bytesErr)
		return user, errors.New("未登录")
	}
	return user, nil
}

// OauthTokenRedis 将用户登录信息存到 redis
func OauthTokenToRedis(user model.UserOauth) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		logrus.Error(err)
		return errors.New("error")
	}
	loginUserKey := fmt.Sprintf("%s%d", model.LoginUser, user.UserId)

	RedisConn := db.RedisPool.Get()
	defer RedisConn.Close()

	if _, redisErr := RedisConn.Do("SET", loginUserKey, userBytes, "EX", config.ServerConfig.TokenMaxAge); redisErr != nil {
		logrus.Errorf("redis set failed: %s", redisErr.Error())
		return errors.New("error")
	}
	return nil
}

// 校验及登录
func Signin(ctx iris.Context, userJson *model.UserLogin, isAdmin bool) (model.MyMap, string, int) {
	var sql string
	if userJson.LoginType == "email" {
		if emailFormat := utils.VerifyEmail(userJson.SigninInput); emailFormat != true {
			logrus.Errorf("login email error %s", userJson.SigninInput)
			return nil, "参数错误", model.STATUS_FAILED
		}
		sql = "email = ?"
	} else if userJson.LoginType == "phone" {
		if emailFormat := utils.VerifyPhone(userJson.SigninInput); emailFormat != true {
			logrus.Errorf("login phone error %s", userJson.SigninInput)
			return nil, "参数错误", model.STATUS_FAILED
		}
		sql = "phone = ?"
	} else {
		return nil, "参数无效", model.STATUS_FAILED
	}

	// 校验图片验证码
	if err := CheckCaptchaCode(userJson.UID, userJson.VerifyCode); err != nil {
		return nil, err.Error(), model.STATUS_FAILED
	}

	var users model.User
	if err := db.DB.Where(sql, userJson.SigninInput).First(&users).Error; err != nil {
		return nil, "账号不存在", model.STATUS_FAILED
	}

	if ok := users.CheckPassword(userJson.Password); ok != true {
		msg := "用户名或者密码有误"
		logrus.Warnf("login failed u:%s,p:%s\n", userJson.SigninInput, userJson.Password)
		return nil, msg, model.STATUS_AUTH_ERROR
	}
	// fmt.Println("UserStatus:",user.UserStatus,time.Now().Unix())
	if userStatus := users.UserStatus; userStatus > time.Now().Unix() {
		msg := fmt.Sprintf("当前用户已经被锁定，请于 %s 之后尝试", utils.TimeStFormat(userStatus))
		return nil, msg, model.STATUS_EXPIRED
	}
	// 如果是管理员账号登录，需要判断用户角色
	if isAdmin && users.RoleID == model.UserRoleNormal {
		return nil, "当前账号不存在", model.STATUS_FAILED
	}
	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  users.UserID,
		"exp": time.Now().Add(time.Hour * time.Duration(2)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.ServerConfig.TokenSecret))
	if err != nil {
		logrus.Error(err)
		return nil, "内部错误", model.ErrorCode.ERROR
	}
	useAgent := ctx.GetHeader("User-Agent")
	loginMachine := utils.GetMachineType(useAgent)
	// 	X-Forwarded-For
	registerIp := utils.RemoteAddr(ctx)
	var oauthToken model.UserOauth
	nowTime := utils.GetNowTime()
	oauthToken.Token = tokenString
	oauthToken.UserId = users.UserID
	oauthToken.RoleId = users.RoleID
	oauthToken.Secret = config.ServerConfig.TokenSecret
	oauthToken.Revoked = false
	oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(2)).Unix()
	oauthToken.LoginAt = nowTime
	oauthToken.LoginType = loginMachine
	oauthToken.LoginIp = registerIp
	oauthToken.LoginCity = common.GetIpCity(registerIp)

	// 创建登录记录
	if err := db.DB.Create(&oauthToken).Error; err != nil {
		logrus.Error(err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	// 更新用户最后登录记录表
	if err := db.DB.Model(users).Where("user_id = ?", users.UserID).Update("last_login_at", nowTime).Error; err != nil {
		logrus.Error(err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	// times := utils.GetRandNumInt(10000)
	// logrus.Debug(times)
	// time.Sleep(time.Duration(times) * time.Millisecond)
	// 将用户信息加入缓存
	// if err := UserToRedis(users); err != nil {
	// 	return nil, "error", model.STATUS_FAILED
	// }
	// 将用户登录信息存入缓存
	if err := OauthTokenToRedis(oauthToken); err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	// 设置 cookies
	// ctx.SetCookie(&http.Cookie{
	// 	Name:     "token",
	// 	Value:    tokenString,
	// 	MaxAge:   config.ServerConfig.TokenMaxAge,
	// 	Path:     "/",
	// 	Domain:   "",
	// 	HttpOnly: true,
	// 	Secure:   true,
	// })

	response := model.MyMap{
		"access_token": tokenString,
	}

	return response, "success", model.STATUS_SUCCESS
}

// 注册
func Signup(ctx iris.Context, userData *model.UserJson) (string, int) {

	// 校验邮箱验证码
	if err := common.CheckVerifyCode(userData.UID, userData.Email, userData.VerifyCode, false, true); err != nil {
		return err.Error(), model.STATUS_FAILED
	}

	// 校验参数
	// userData.Username = utils.AvoidXSS(userData.Username)
	userData.Phone = strings.TrimSpace(userData.Phone)
	userData.Email = strings.TrimSpace(userData.Email)

	if EmailFormat := utils.VerifyEmail(userData.Email); EmailFormat != true {
		return "邮箱有误", model.STATUS_FAILED
	}
	if passwdFormat := utils.VerifyPassword(userData.Password); passwdFormat != true {
		logrus.Errorf("密码格式错误")
		return "密码必须包含大小字符，特殊字符两种以上，并且长度在6到20位之间", model.STATUS_FAILED
	}

	var user model.User
	if err := db.DB.Where("email = ?", userData.Email).Find(&user).Error; err == nil {
		var msg string
		if userData.Phone != "" && user.Phone == userData.Phone {
			msg = fmt.Sprintf("手机 %s 已被注册", user.Phone)
		} else if user.Email == userData.Email {
			msg = fmt.Sprintf("邮箱 %s 已存在", user.Email)
		}
		return msg, model.STATUS_FAILED
	}
	// 获取last users id
	userId := 1
	if err := db.DB.Last(&user).Error; err == nil {
		userId = int(user.ID) + 1
	}
	// reqHost := ctx.RemoteAddr()
	registerIp := utils.RemoteAddr(ctx)
	var newUser model.User
	newUser.UserID = userId
	newUser.Username = userData.Email
	newUser.Password = newUser.EncryptPassword(userData.Password)
	newUser.Email = userData.Email
	newUser.Phone = userData.Phone
	newUser.RoleID = model.UserRoleNormal
	newUser.UserStatus = model.UserStatusNormal
	newUser.RegisterAt = utils.GetNowTime()
	newUser.RegisterIp = registerIp
	if err := db.DB.Create(&newUser).Error; err != nil {
		logrus.Error(err)
		return "error", model.STATUS_FAILED
	}

	// 创建账户 添加,虚拟货币地址分配
	if err := AddUserAccount(&newUser); err != nil {
		return "error", model.STATUS_FAILED
	}

	return "success", model.STATUS_SUCCESS
}

// ResetPassword 重置密码
func ResetPassword(ctx iris.Context, userJson *model.UserUpdatePwd) (string, int) {

	var sql string
	if userJson.Type == "email" {
		if emailFormat := utils.VerifyEmail(userJson.AddressInput); emailFormat != true {
			logrus.Errorf("%s 邮箱错误", userJson.AddressInput)
			return "参数错误", model.STATUS_FAILED
		}

		// 重置密码邮件
		cacheKey := fmt.Sprintf("%s%s", model.ResetTime, userJson.AddressInput)
		// 校验邮箱验证码
		if err := common.CheckVerifyCode(userJson.UID, cacheKey, userJson.VerifyCode, false, true); err != nil {
			return err.Error(), model.STATUS_FAILED
		}
		sql = "email = ?"
	} else if userJson.Type == "phone" {
		if emailFormat := utils.VerifyPhone(userJson.AddressInput); emailFormat != true {
			logrus.Errorf("%s 手机号码错误", userJson.AddressInput)
			return "参数错误", model.STATUS_FAILED
		}

		// 重置密码邮件
		cacheKey := fmt.Sprintf("%s%s", model.ResetTime, userJson.AddressInput)
		// 校验手机验证码
		if err := common.CheckVerifyCode(userJson.UID, cacheKey, userJson.VerifyCode, false, true); err != nil {
			return err.Error(), model.STATUS_FAILED
		}
		sql = "phone = ?"
	} else {
		return "参数无效", model.STATUS_FAILED
	}

	if passwdFormat := utils.VerifyPassword(userJson.Password); passwdFormat != true {
		logrus.Errorf("密码格式错误")
		return "密码必须包含大小字符，特殊字符两种以上，并且长度在6到20位之间", model.STATUS_FAILED
	}

	var users model.User
	if err := db.DB.Where(sql, userJson.AddressInput).First(&users).Error; err != nil {
		return "账号不存在", model.STATUS_FAILED
	}

	// fmt.Println("UserStatus:",user.UserStatus,time.Now().Unix())
	if userStatus := users.UserStatus; userStatus > time.Now().Unix() {
		msg := fmt.Sprintf("当前用户已经被锁定，请于 %s 之后尝试", utils.TimeStFormat(userStatus))
		return msg, model.STATUS_EXPIRED
	}
	newPass := users.EncryptPassword(userJson.Password)
	if err := db.DB.Model(&users).Update("password", newPass).Error; err != nil {
		logrus.Error(err)
		return "error", model.STATUS_FAILED
	}

	return "success", model.STATUS_SUCCESS
}

// 获取手机或者邮箱验证码
func CreateEmailPhoneCode(verifyCode *model.VerifyCodeJson) (string, int) {
	// 判断缓存是否已经发送验证码
	newAddressKey := fmt.Sprintf("%s%s", model.SendCode, verifyCode.Address)
	_, found := cache.OC.Get(newAddressKey)
	if found {
		return "success", model.STATUS_SUCCESS
	}
	// 检测邮箱或者电话
	if err := CheckEmailPhonCodeOption(verifyCode); err != nil {
		return err.Error(), model.STATUS_FAILED
	}
	var code string
	var err error
	if verifyCode.Type == 1 {
		code, err = common.GetEmailVerifyCodeAndSend(verifyCode.Address, verifyCode.Option)
		if err != nil {
			logrus.Errorf("CreateEmailPhoneCode error: %s", err.Error())
			return err.Error(), model.STATUS_FAILED
		}
		// 手机验证码，暂时未开启
		// } else if verifyCode.Type == 2 {
		// 	code, err = GetPhoneVerifyCodeAndSend(verifyCode.Address)
		// 	if err != nil {
		// 		return nil, "failed", model.STATUS_FAILED
		// 	}
	} else {
		return "参数错误", model.STATUS_PARAM_ERROR
	}
	res := &model.VerifyCodeRes{
		UID:        verifyCode.UID,
		VerifyCode: code,
	}
	// 将当前验证码加入缓存,并且设置有效期
	// 将验证码接受地址作为 key, 验证码结果以及 uid 作为 value
	cache.OC.Set(verifyCode.Address, res, cache.CacheDefaultExpiration)

	// 缓存当前地址，一分钟发送一次
	cache.OC.Set(newAddressKey, 1, cache.SendCodeDefaultExpiration)

	return "success", model.STATUS_SUCCESS
}

// Signout 退出登录
func Signout(ctx iris.Context) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

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

// GetUserAllAmount 获取用户账户信息
func GetUserAllAmount(ctx iris.Context) (*[]model.UserAccountRes, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var userAccount []model.UserAccounts
	var userAccountRes []model.UserAccountRes
	if err := db.DB.Where("user_id = ?", userId).Find(&userAccount).Scan(&userAccountRes).Error; err != nil {
		logrus.Errorf("GetUserAllAmountFromDB service.DB get user accountAmount err %s", err)
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetUserAllAmountFromDB  userTd: %d", userId)
	return &userAccountRes, "success", model.STATUS_SUCCESS
}

// GetUserAssetFlow 获取用户资产流水记录
func GetUserAssetFlow(ctx iris.Context, assetflow *model.AssetflowJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var assetflowRes []model.AssetflowRes
	var userAssetflow []model.UserAssetflow
	count := 0
	dbs := db.DB
	if assetflow.StartAt != 0 {
		dbs = dbs.Where("create_at >= ?", assetflow.StartAt) // 时间戳
	}
	if assetflow.EndAt != 0 {
		dbs = dbs.Where("create_at < ?", assetflow.EndAt)
	}
	if assetflow.CoinType != "" {
		dbs = dbs.Where("coin_type = ?", assetflow.CoinType)
	}
	offset, err := db.GetOffset(assetflow.Page, assetflow.Limit)
	if err != nil {
		return nil, "参数错误", model.STATUS_FAILED
	}
	if err := dbs.Model(model.UserAssetflow{}).Where("user_id = ?", userId).Count(&count).
		Offset(offset).Limit(assetflow.Limit).
		Find(&userAssetflow).Scan(&assetflowRes).Error; err != nil {
		logrus.Errorf("GetUserAssetFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetUserAssetFromDB useid: %d page: %d limit %d offset: %d\n", userId, assetflow.Page, assetflow.Limit, offset)
	response := model.MyMap{
		"data":   userAssetflow,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetUserProfileInfo 获取用户信息
func GetUserProfileInfo(ctx iris.Context) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	users, err := GetUserById(userId)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}

	newEmail := utils.StrWithShelter(users.Email, 1)
	newPhone := utils.StrWithShelter(users.Phone, 2)

	response := model.MyMap{
		// "username":    users.Username,
		"email":       newEmail,
		"user_id":     users.UserID + 10000,
		"phone":       newPhone,
		"register_at": users.RegisterAt,
	}
	return response, "success", model.STATUS_SUCCESS
}

// UpdateUesrName 更新用户名信息
func UpdateUserName(ctx iris.Context, userReqData *model.UserUpdateJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	if emailFormat := utils.VerifyUserName(userReqData.Username); emailFormat != true {
		return "用户名格式错误", model.ErrorCode.ERROR
	}
	username := utils.AvoidXSS(userReqData.Username)
	username = strings.TrimSpace(username)

	var users model.User

	if err := db.DB.Model(&users).Where("user_id", userId).Update("username", username).Error; err != nil {
		logrus.Error(err)
		return "error", model.ErrorCode.ERROR
	}

	return "success", model.ErrorCode.SUCCESS
}

// UpdatePassword 更新用户密码
func UpdatePassword(ctx iris.Context, userData *model.PasswordUpdateData) (string, int) {

	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	if passwdFormat := utils.VerifyPassword(userData.Password); passwdFormat != true {
		logrus.Errorf("密码格式错误")
		return "密码必须包含大小字符，特殊字符两种以上，并且长度在6到20位之间", model.STATUS_FAILED
	}

	users, err := GetUserById(userId)
	if err != nil {
		return "error", model.STATUS_FAILED
	}

	if ok := users.CheckPassword(userData.Password); ok != true {
		return "原密码错误", model.STATUS_FAILED
	}

	users.Password = users.EncryptPassword(userData.NewPwd)
	if err := db.DB.Save(&users).Error; err != nil {
		return "原密码不正确", model.STATUS_FAILED
	}

	return "success", model.ErrorCode.SUCCESS
}

// GetUserLoginList 获取用户登录记录
func GetUserLoginList(ctx iris.Context, userOauthJson *model.UserOauthJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var oauthList []model.UserOauth
	var oauthRes []model.UserOauthRes
	count := 0
	offset, err := db.GetOffset(userOauthJson.Page, userOauthJson.Limit)
	if err != nil {
		return nil, "error", model.STATUS_FAILED
	}
	var userOauth model.UserOauth
	if err := db.DB.Model(userOauth).Order("id desc").Where("user_id =?", userId).
		Offset(offset).Limit(userOauthJson.Limit).Find(&oauthList).Scan(&oauthRes).Error; err != nil {
		logrus.Errorf("GetUserLoginListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	count = len(oauthRes)
	logrus.Debugf("GetUserMessageListFromDB useId: %d page: %d limit: %d\n", userOauthJson.UserId, userOauthJson.Page, userOauthJson.Limit)
	response := model.MyMap{
		"data":   oauthRes,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// 获取用户消息列表
func GetUserMessages(ctx iris.Context, message *model.UserMessageJson) (model.MyMap, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var messsageList []model.UserMessage
	var msgRes []model.UserMessageRes
	count := 0
	offset, err := db.GetOffset(message.Page, message.Limit)
	if err != nil {
		return nil, "参数无效", model.STATUS_FAILED
	}
	if err := db.DB.Where("user_id = ? AND status = ?", userId, true).Count(&count).
		Offset(offset).Limit(message.Limit).
		Find(&messsageList).Scan(&msgRes).Error; err != nil {
		logrus.Errorf("GetUserMessageListFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetUserMessageListFromDB useId: %d page: %d limit %d offset: %d\n", userId, message.Page, message.Limit, offset)
	response := model.MyMap{
		"data":   messsageList,
		"length": count,
	}
	return response, "success", model.STATUS_SUCCESS
}

// GetUserMessageDetail 获取消息详情
func GetUserMessageDetail(ctx iris.Context, msgDetailJson *model.MessageDetailJson) (*model.MessageDetailRes, string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return nil, "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	var message model.UserMessage
	var msgDetail model.MessageDetailRes
	// Select("article_id,pushed_at,announcer,title,summary,readed")
	if err := db.DB.Where("user_id = ? AND status = ? AND article_id = ?", userId, true, msgDetailJson.ArticleId).First(&message).Scan(&msgDetail).Error; err != nil {
		logrus.Errorf("GetUserMessageDetailFromDB failed err: %s", err.Error())
		return nil, "error", model.STATUS_FAILED
	}
	logrus.Debugf("GetUserMessageDetailFromDB useId: %d\n", userId)

	return &msgDetail, "success", model.STATUS_SUCCESS
}

// UpdateUserMessageStatus 更新用户消息状态
func UpdateUserMessageStatus(ctx iris.Context, msgDetailJson *model.MessageDetailJson) (string, int) {
	userInter := ctx.Values().Get("auth_user_id")
	if userInter == nil {
		return "未登录", model.STATUS_FAILED
	}
	userId := userInter.(int)

	nowTime := utils.GetNowTime()

	updateData := map[string]interface{}{
		"readed": true,
		"time":   nowTime,
	}
	message := new(model.UserMessage)
	if err := db.DB.Model(message).Where("user_id = ? AND readed = ? AND article_id = ?", userId, false, msgDetailJson.ArticleId).Updates(updateData).Error; err != nil {
		logrus.Errorf("UpdateUserMessageStatus falied Err:%s", err.Error())
		return "error", model.STATUS_FAILED
	}
	return "success", model.STATUS_SUCCESS
}

// 判断验证码业务类型
func CheckEmailPhonCodeOption(verfiyCode *model.VerifyCodeJson) error {
	switch verfiyCode.Option {
	case 1: // 用户注册，校验邮箱未使用
		if status := CheckUserNotExist(verfiyCode.Address); status == false {
			return errors.New("用户已存在")
		}
		return nil
	case 2: // 用户重置密码，校验用户邮箱存在
		if CheckUserNotExist(verfiyCode.Address) == true {
			return errors.New("用户不存在")
		}
		return nil
	default:
		return errors.New("参数错误")
	}
}

// 校验图片验证码
func CheckCaptchaCode(uid, verifyCode string) error {
	if model.CHECK_VERIFY_CODE == false {
		logrus.Info("captcha code not check")
		return nil
	}
	if uid == "" {
		return errors.New("no uid")
	}
	if err := common.CaptchaVerify(uid, verifyCode); err != nil {
		return errors.New("错误的验证码")
	}
	return nil
}

// 添加用户账户地址
func AddUserAccount(user *model.User) error {
	// 1. CNY
	coinType := "CNY"
	userAccount := &model.UserAccounts{
		UserId:        user.UserID,
		CoinAmount:    "0",
		CoinType:      coinType,
		VirtualAmount: "0",
	}
	UserAddAccount(userAccount)

	// 2. MGD
	coinType = "MGD"
	coinAddr := GetCoinAddressFromDB(coinType)
	mgdUserAccount := &model.UserAccounts{
		UserId:        user.UserID,
		CoinAmount:    "0",
		CoinAddr:      coinAddr,
		CoinType:      coinType,
		VirtualAmount: "0",
	}
	if err := UserAddAccount(mgdUserAccount); err == nil { // 初始化帐号
		// 更新地址分配记录
		if coinAddr != "" {
			UpdateCoinAddressToDB(coinType, coinAddr, user.UserID)

			// 添加地址到 mgd 同步线程 !!!
			if err := wallet.UpdateWalletAddress(coinType, coinAddr, user.UserID); err != nil {
				logrus.Errorf("AddUserAccount UpdateWalletAddress cointype %s coinaddr %s userid %d err: %s", coinType, coinAddr, user.UserID, err.Error())
			}
		}

	}
	// 3. ETH
	coinType = "ETH"
	coinAddr = GetCoinAddressFromDB(coinType)
	ethUserAccount := &model.UserAccounts{
		UserId:        user.UserID,
		CoinAmount:    "0",
		CoinAddr:      coinAddr,
		CoinType:      coinType,
		VirtualAmount: "0",
	}
	if err := UserAddAccount(ethUserAccount); err == nil { // 初始化帐号
		// 更新地址分配记录
		if coinAddr != "" {
			UpdateCoinAddressToDB(coinType, coinAddr, user.UserID)
			// 添加地址到 eth 同步线程 !!!
			if err := wallet.UpdateWalletAddress(coinType, coinAddr, user.UserID); err != nil {
				logrus.Errorf("AddUserAccount UpdateWalletAddress cointype %s coinaddr %s userid %d err: %s", coinType, coinAddr, user.UserID, err.Error())
			}
		}
	}
	return nil
}

/**
检测 用户是否已经存在
*/
func CheckUserNotExist(userAccount string) bool {
	user := new(model.User)
	status := db.DB.Where("email = ? OR phone = ?", userAccount, userAccount).First(user).RecordNotFound()
	return status
}

/**
 * 检测并创建用户账户
 */
func UserAddAccount(userAccount *model.UserAccounts) error {
	userAcc := new(model.UserAccounts)
	if err := db.DB.Where("user_id = ? AND coin_type = ?", userAccount.UserId, userAccount.CoinType).First(&userAcc).RecordNotFound(); err != true {
		logrus.Infof("UserAddAccount  already exist: %+v", userAccount)
		return errors.New("already exist")
	}
	if err := db.DB.Create(userAccount).Error; err != nil {
		logrus.Errorf("UserAddAccount failed err: %s", err.Error())
		return err
	}
	logrus.Debugf("UserAddAccount add account: %+v", userAccount)
	return nil
}

// 获取未分配地址
func GetCoinAddressFromDB(coinType string) string {
	wtAddress := new(model.WalletAddress)
	if err := db.DB.Where("coin_type = ? AND allocated = ?", coinType, false).First(&wtAddress).Error; err != nil {
		logrus.Errorf("GetCoinAddressFromDB coinType: %s err: %s", coinType, err.Error())
		return ""
	}
	return wtAddress.Address
}

// 更新未分配地址表
func UpdateCoinAddressToDB(coinType, coinAddr string, userId int) error {
	wtAddress := new(model.WalletAddress)
	if err := db.DB.Model(wtAddress).Where("coin_type = ? AND address = ?  AND allocated = ?", coinType, coinAddr, false).Update(map[string]interface{}{"user_id": userId, "allocated": true}).Error; err != nil {
		logrus.Errorf("UpdateCoinAddressToDB falied Err:%s", err.Error())
		return err
	}
	return nil
}

// 获取用户详情
func GetUserById(id int) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("user_id = ?", id).First(&user).Error; err != nil {
		logrus.Errorf("GetUserById failed err: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

// GetUserOauthByToken 获取touken
func GetUserOauthByToken(token string) *model.UserOauth {
	var ot model.UserOauth
	if err := db.DB.Where("token =  ?", token).First(&ot).Error; err != nil {
		logrus.Errorf("UserOauthCreate failed err: %s", err.Error())
		return nil
	}
	return &ot
}

// UpdateUserOauthByToken 通过token 更新 用户登录超时时间
func UpdateUserOauthByToken(newExpirTime int64, token string) error {
	ot := new(model.UserOauth)
	if err := db.DB.Model(ot).Where("token =  ?", token).Update("express_in", newExpirTime).Error; err != nil {
		logrus.Errorf("UpdateUserOauthByToken failed err: %s", err.Error())
		return nil
	}
	return nil
}

/* ************************old ************************************************ */

// // admin
// // AllList 查询用户列表，只有管理员才能调此接口
// func AllList(ctx iris.Context, userData *model.AdminUserListData) (map[string]interface{}, string, int) {
// 	role, _ := strconv.Atoi(userData.Role)
// 	allUserRole := []int{
// 		model.UserRoleNormal,
// 		model.UserRoleEditor,
// 		model.UserRoleAdmin,
// 		model.UserRoleCrawler,
// 		model.UserRoleSuperAdmin,
// 	}
// 	foundRole := false
// 	for _, r := range allUserRole {
// 		if r == role {
// 			foundRole = true
// 			break
// 		}
// 	}

// 	var startTime string
// 	var endTime string

// 	if startAt, err := strconv.Atoi(userData.StartAt); err != nil {
// 		startTime = time.Unix(0, 0).Format("2006-01-02 15:04:05")
// 	} else {
// 		startTime = time.Unix(int64(startAt/1000), 0).Format("2006-01-02 15:04:05")
// 	}

// 	if endAt, err := strconv.Atoi(userData.EndAt); err != nil {
// 		endTime = time.Now().Format("2006-01-02 15:04:05")
// 	} else {
// 		endTime = time.Unix(int64(endAt/1000), 0).Format("2006-01-02 15:04:05")
// 	}

// 	pageNo, pageNoErr := strconv.Atoi(userData.PageNo)
// 	if pageNoErr != nil {
// 		pageNo = 1
// 	}
// 	if pageNo < 1 {
// 		pageNo = 1
// 	}

// 	offset := (pageNo - 1) * model.PageSize
// 	pageSize := model.PageSize

// 	var users []model.User
// 	var totalCount int
// 	if foundRole {
// 		if err := db.DB.Model(&model.User{}).Where("created_at >= ? AND created_at < ? AND role = ?", startTime, endTime, role).
// 			Count(&totalCount).Error; err != nil {
// 			fmt.Println(err.Error())
// 			return nil, "error", model.ErrorCode.ERROR
// 		}
// 		if err := db.DB.Where("created_at >= ? AND created_at < ? AND role = ?", startTime, endTime, role).
// 			Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
// 			fmt.Println(err.Error())
// 			return nil, "error", model.ErrorCode.ERROR
// 		}
// 	} else {
// 		if err := db.DB.Model(&model.User{}).Where("created_at >= ? AND created_at < ?", startTime, endTime).
// 			Count(&totalCount).Error; err != nil {
// 			logrus.Error(err.Error())
// 			return nil, "error", model.ErrorCode.ERROR
// 		}
// 		if err := db.DB.Where("created_at >= ? AND created_at < ?", startTime, endTime).Order("created_at DESC").Offset(offset).
// 			Limit(pageSize).Find(&users).Error; err != nil {
// 			logrus.Error(err.Error())
// 			return nil, "error", model.ErrorCode.ERROR
// 		}
// 	}
// 	var results []interface{}
// 	for i := 0; i < len(users); i++ {
// 		results = append(results, map[string]interface{}{
// 			"id":          users[i].ID,
// 			"name":        users[i].Name,
// 			"email":       users[i].Email,
// 			"role":        users[i].Role,
// 			"status":      users[i].Status,
// 			"createdAt":   users[i].CreatedAt,
// 			"activatedAt": users[i].ActivatedAt,
// 		})
// 	}
// 	res := map[string]interface{}{
// 		"users":      results,
// 		"pageNo":     pageNo,
// 		"pageSize":   pageSize,
// 		"totalCount": totalCount,
// 	}
// 	return res, "success", model.ErrorCode.SUCCESS
// }
