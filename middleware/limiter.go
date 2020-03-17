package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/kataras/iris"
	"github.com/kataras/iris/v12/context"
	"github.com/zhxx123/gomonitor/controller"
	"github.com/zhxx123/gomonitor/model"
)

// LimitHandler is a middleware that performs
// rate-limiting given a "limiter" configuration.
func LimitHandler(lmt *limiter.Limiter) context.Handler {
	return func(ctx context.Context) {
		httpError := tollbooth.LimitByRequest(lmt, ctx.ResponseWriter(), ctx.Request())
		if httpError != nil {
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(controller.ApiResource(model.STATUS_FREQUENT_LIMIT, nil, httpError.Message))
			return
		}
		ctx.Next()
	}
}
