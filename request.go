package infra_gin

import (
	"cmp"
	"errors"
	"github.com/bucketheadv/infra-core/basic"
	"github.com/gin-gonic/gin"
)

var (
	ParamError        = errors.New("参数错误")
	ParamConvertError = errors.New("参数转换错误")
	ParamBlankError   = errors.New("参数为空")
)

func GetQuery[T cmp.Ordered](c *gin.Context, key string) (T, error) {
	var v T
	q, success := c.GetQuery(key)
	if !success {
		return v, ParamError
	}

	if q == "" {
		return v, ParamBlankError
	}

	err := basic.ConvertStringTo(q, &v)
	return v, err
}
