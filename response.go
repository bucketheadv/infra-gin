package infra

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func ApiResponseOk[T any](c *gin.Context, response Response[T]) {
	if response.Code == 0 {
		response.Code = http.StatusOK
	}
	ApiResponse(c, response)
}

func ApiResponseError[T any](c *gin.Context, response Response[T]) {
	if response.Code == 0 {
		response.Code = http.StatusInternalServerError
	}
	ApiResponse(c, response)
}

func ApiResponse[T any](c *gin.Context, response Response[T]) {
	c.JSON(response.Code, response)
}
