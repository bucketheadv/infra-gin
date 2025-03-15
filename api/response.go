package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func ResponseOk[T any](c *gin.Context, response Response[T]) {
	ResponseJSON(c, response, http.StatusOK)
}

func ResponseError[T any](c *gin.Context, response Response[T], httpStatus int) {
	ResponseJSON(c, response, httpStatus)
}

func ResponseJSON[T any](c *gin.Context, response Response[T], httpStatus int) {
	c.JSON(httpStatus, response)
}
