package middlewares

import (
	"github.com/bucketheadv/infragin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegErrorHandler(e *gin.Engine) {
	e.Use(globalPanicHandler())
	e.Use(globalErrorHandler())
	e.NoRoute(func(c *gin.Context) {
		var response = infragin.Response[any]{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
		infragin.ApiResponseError(c, response)
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
			if r := recover(); r != nil {
				var response = infragin.Response[any]{
					Code:    http.StatusInternalServerError,
					Message: errorToString(r),
				}
				infragin.ApiResponseError(c, response)
			}
		}()
		c.Next()
	}
}

func globalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			var response = infragin.Response[any]{
				Code:    http.StatusInternalServerError,
				Message: c.Errors.String(),
			}
			infragin.ApiResponseError(c, response)
			c.Abort()
			return
		}
	}
}
