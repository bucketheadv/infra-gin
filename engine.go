package infra_gin

import (
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func init() {
	Engine = gin.Default()
}
