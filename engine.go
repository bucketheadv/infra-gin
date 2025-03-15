package infra_gin

import (
	"github.com/bucketheadv/infra-gin/middlewares"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func init() {
	Engine = gin.Default()
	middlewares.RegErrorHandler(Engine)
}
