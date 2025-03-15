package api

import (
	"cmp"
	"github.com/bucketheadv/infra-core/basic"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrParamInvalid = NewError(http.StatusBadRequest, "参数错误")
	ErrParamBlank   = NewError(http.StatusBadRequest, "参数为空")
)

func GetQuery[T cmp.Ordered](c *gin.Context, key string) (T, error) {
	var v T
	q, success := c.GetQuery(key)
	if !success {
		return v, ErrParamInvalid
	}

	if q == "" {
		return v, ErrParamBlank
	}

	v, err := basic.StringTo[T](q)
	return v, err
}
