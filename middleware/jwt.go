package middleware

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/zhxx123/gomonitor/config"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() *jwtmiddleware.Middleware {
	secrct := config.ServerConfig.TokenSecret
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(secrct), nil
		},

		SigningMethod: jwt.SigningMethodHS256,
	})

}
