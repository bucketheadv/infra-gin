package api

import (
	"cmp"
	"github.com/bucketheadv/infra-core/basic"
	"github.com/gin-gonic/gin"
)

var (
	ErrParamInvalid = NewParamError("参数%s错误")
	ErrParamBlank   = NewParamError("参数%s不能为空")
)

func GetQuery[T cmp.Ordered](c *gin.Context, key string) (T, error) {
	var v T
	q, success := c.GetQuery(key)
	if !success {
		var e = ErrParamInvalid.Format(key)
		return v, &e
	}

	if q == "" {
		var e = ErrParamBlank.Format(key)
		return v, &e
	}

	v, err := basic.StringTo[T](q)
	return v, err
}
