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
		api.ResponseError(c, response, http.StatusNotFound)
	})
}

func resolveCodeError(r any) (api.BizError, int) {
	switch v := r.(type) {
	case error:
		return resolveErrorType(v)
	default:
		return api.BizError{
			Code:    http.StatusInternalServerError,
			Message: v.(string),
		}, http.StatusInternalServerError
	}
}

func resolveErrorType(r error) (api.BizError, int) {
	var e *api.BizError
	if errors.As(r, &e) {
		return *e, http.StatusOK
	}

	var e2 *api.ParamError
	if errors.As(r, &e2) {
		return api.BizError{
			Code:    http.StatusBadRequest,
			Message: e2.Message,
		}, http.StatusBadRequest
	}

	return api.BizError{
		Code:    http.StatusInternalServerError,
		Message: r.Error(),
	}, http.StatusOK
}

func globalPanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			var err, httpStatus = resolveCodeError(r)
			var response = api.Response[any]{
				Code:    err.Code,
				Message: err.Message,
			}
			logger.Errorf("中间件全局Panic捕获: %s\n", err.Message)
			c.AbortWithStatusJSON(httpStatus, response)
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

		var err, httpStatus = resolveCodeError(c.Errors[0])
		var response = api.Response[any]{
			Code:    err.Code,
			Message: c.Errors.String(),
		}
		logger.Errorf("中间件全局Error捕获: %s\n", c.Errors.String())
		c.AbortWithStatusJSON(httpStatus, response)
	}
}
