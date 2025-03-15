package middlewares

import (
	"errors"
	"github.com/bucketheadv/infra-core/modules/logger"
	"github.com/bucketheadv/infra-gin/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegErrorHandler(e *gin.Engine) {
	e.Use(globalPanicHandler())
	e.Use(globalErrorHandler())
	e.NoRoute(func(c *gin.Context) {
		var response = api.Response[any]{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
		api.ApiResponseError(c, response)
	})
}

func resolveCodeError(r any) api.BizError {
	switch v := r.(type) {
	case error:
		var e *api.BizError
		if errors.As(v, &e) {
			return *e
		}
		return api.BizError{
			Code:    http.StatusInternalServerError,
			Message: v.Error(),
		}
	default:
		return api.BizError{
			Code:    http.StatusInternalServerError,
			Message: v.(string),
		}
	}
}

func globalPanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			var err = resolveCodeError(r)
			var response = api.Response[any]{
				Code:    err.Code,
				Message: err.Message,
			}
			logger.Errorf("中间件全局Panic捕获: %s\n", err.Message)
			api.ApiResponseError(c, response)
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

		var err = resolveCodeError(c.Errors[0])
		var response = api.Response[any]{
			Code:    err.Code,
			Message: c.Errors.String(),
		}
		logger.Errorf("中间件全局Error捕获: %s\n", c.Errors.String())
		api.ApiResponseError(c, response)
		c.Abort()
	}
}
