package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func NewEngine(i do.Injector) (*gin.Engine, error) {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	return engine, nil
}
