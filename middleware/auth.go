package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/zhxx123/gomonitor/config"
	"github.com/zhxx123/gomonitor/controller"
	"github.com/zhxx123/gomonitor/model"
	"github.com/zhxx123/gomonitor/service/user"
	users "github.com/zhxx123/gomonitor/service/user"
	// "github.com/zhxx123/gomonitor/controller/common"
)

func getUser(ctx iris.Context) (model.User, error) {
	var user model.User
	tokenString := ctx.GetCookie("token")

	if tokenString == "" {
		return user, errors.New("未登录")
	}

	token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ServerConfig.TokenSecret), nil
	})

	if tokenErr != nil {
		return user, errors.New("未登录")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["id"].(float64))
		var err error
		user, err = users.UserFromRedis(userID)
		if err != nil {
			return user, errors.New("未登录")
		}
		return user, nil
	}
	return user, errors.New("未登录")
}

// SigninRequired 必须是登录用户
func SigninRequired(ctx iris.Context) {
	var user model.User
	var err error
	if user, err = getUser(ctx); err != nil {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(controller.ApiResource(model.STATUS_EXPIRED, nil, "未登录"))
		return
	}
	ctx.Values().Set("user", user)
	ctx.Next()
}

// AdminRequired 必须是管理员
func AdminRequired(ctx iris.Context) {
	var user model.User
	var err error
	if user, err = getUser(ctx); err != nil {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(controller.ApiResource(model.STATUS_EXPIRED, nil, "未登录"))
		return
	}
	if user.RoleID == model.UserRoleAdmin || user.RoleID == model.UserRoleCrawler || user.RoleID == model.UserRoleSuperAdmin {
		ctx.Values().Set("user", user)
		ctx.Next()
	} else {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(controller.ApiResource(model.STATUS_AUTH_ERROR, nil, "没有权限"))
		return

	}
}

//*************************************************** old method
/**
 * 判断 token 是否有效
 * 如果有效 就获取信息并且保存到请求里面
 * @method AuthToken
 * @param  {[type]}  ctx       iris.Context    [description]
 */
func getOauthToken(tokenString string) (model.UserOauth, error) {
	var user model.UserOauth

	if tokenString == "" {
		return user, errors.New("未登录")
	}

	token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ServerConfig.TokenSecret), nil
	})

	if tokenErr != nil {
		return user, errors.New("未登录")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["id"].(float64))
		var err error
		user, err = users.OauthTokenFromRedis(userID)
		if err != nil {
			return user, errors.New("未登录")
		}
		return user, nil
	}
	return user, errors.New("未登录")
}

func AuthToken(ctx iris.Context) {
	u := ctx.Values().Get("jwt").(*jwt.Token) //获取 token 信息
	token, err := getOauthToken(u.Raw)        //获取 access_token 信息
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(controller.ApiResource(model.STATUS_EXPIRED, nil, "请登录"))
		return
	}
	// fmt.Println("+v%", token)
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(controller.ApiResource(model.STATUS_EXPIRED, nil, "登录已失效"))
		return
	} else {
		ctx.Values().Set("auth_user_id", token.UserId)
	}

	// 更新 token 信息
	newExpire := time.Now().Add(time.Hour * time.Duration(2)).Unix() // 2小时
	token.ExpressIn = newExpire
	if newExpire != token.ExpressIn { // 时间不一致时候，更新过期时间
		user.OauthTokenToRedis(token) // 获取 access_token 信息
	}
	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}

func AuthAdminToken(ctx iris.Context) {
	u := ctx.Values().Get("jwt").(*jwt.Token) // 获取 token 信息
	token, err := getOauthToken(u.Raw)        // 获取 access_token 信息
	// fmt.Println(token)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(controller.ApiResource(model.STATUS_EXPIRED, nil, "请登录"))
		return
	}
	// logrus.Debug("autoToken:", token)
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(controller.ApiResource(model.STATUS_EXPIRED, nil, "登录已失效"))
		return
	} else {
		ctx.Values().Set("user_id", token.UserId)
		ctx.Values().Set("role_id", token.RoleId)
	}
	// 更新 token 信息
	newExpire := time.Now().Add(time.Hour * time.Duration(2)).Unix() // 2小时
	token.ExpressIn = newExpire
	if newExpire != token.ExpressIn { // 时间不一致时候，更新过期时间
		user.OauthTokenToRedis(token) // 获取 access_token 信息
	}
	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}
