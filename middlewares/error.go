package middlewares

import (
	"github.com/bucketheadv/infra-core/modules/logger"
	"github.com/bucketheadv/infra-gin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegErrorHandler(e *gin.Engine) {
	e.Use(globalPanicHandler())
	e.Use(globalErrorHandler())
	e.NoRoute(func(c *gin.Context) {
		var response = infra_gin.Response[any]{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
		infra_gin.ApiResponseError(c, response)
	})
}

func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return v.(string)
	}
}

func globalPanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			var response = infra_gin.Response[any]{
				Code:    http.StatusInternalServerError,
				Message: errorToString(r),
			}
			logger.Errorf("中间件全局Panic捕获: %s\n", errorToString(r))
			infra_gin.ApiResponseError(c, response)
		}()
		c.Next()
	}
}

func globalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) <= 0 {
			return
		}

		var response = infra_gin.Response[any]{
			Code:    http.StatusInternalServerError,
			Message: c.Errors.String(),
		}
		logger.Errorf("中间件全局Error捕获: %s\n", errorToString(c.Errors.String()))
		infra_gin.ApiResponseError(c, response)
		c.Abort()
	}
}
