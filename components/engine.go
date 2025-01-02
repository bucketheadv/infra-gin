package components

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var engine *gin.Engine

var lock sync.RWMutex

func GetDefaultEngine() *gin.Engine {
	lock.Lock()
	defer lock.Unlock()
	if engine == nil {
		engine = gin.Default()
	}
	return engine
}
