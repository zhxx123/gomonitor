package middleware

import (
	"github.com/kataras/iris/v12"
)

// RefreshTokenCookie 刷新过期时间
func RefreshTokenCookie(ctx iris.Context) {
	// tokenString := ctx.GetCookie("token")
	// if tokenString != "" {
	// 	ctx.SetCookie(&http.Cookie{
	// 		Name:     "token",
	// 		Value:    tokenString,
	// 		MaxAge:   config.ServerConfig.TokenMaxAge,
	// 		Path:     "/",
	// 		Domain:   "",
	// 		Secure:   true,
	// 		HttpOnly: true,
	// 	})
	// 	if user, err := getUser(ctx); err == nil {
	// 		users.UserToRedis(user)
	// 	}
	// }
	// ctx.Next()

	// old
	// tokenString := ctx.GetCookie("token")
	// if tokenString != "" {

	// 	if user, err := getOauthToken(ctx); err == nil {
	// 		users.UserToRedis(user)
	// 	}
	// }
	ctx.Next()

}
