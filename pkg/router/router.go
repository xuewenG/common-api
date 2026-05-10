package router

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xuewenG/common-api/pkg/config"
	"github.com/xuewenG/common-api/pkg/handler"
)

type Router struct {
	i do.Injector
}

func NewRouter(i do.Injector) (*Router, error) {
	return &Router{
		i: i,
	}, nil
}

func (r *Router) Bind(engine *gin.Engine) {
	cfg := do.MustInvoke[*config.Config](r.i)
	engine.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(cfg.CorsOrigin, ","),
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	// 健康检查
	healthHandler := do.MustInvoke[*handler.HealthHandler](r.i)
	engine.GET("/health", healthHandler.Check)

	// 假期相关接口
	holidayHandler := do.MustInvoke[*handler.HolidayHandler](r.i)
	engine.GET("/holiday/next", holidayHandler.GetNextHoliday)

	// BiliRec 相关接口
	biliRecHandler := do.MustInvoke[*handler.BiliRecHandler](r.i)
	engine.POST("/bilirec/event", biliRecHandler.OnEvent)

	// IP 相关接口
	ipHandler := do.MustInvoke[*handler.IPHandler](r.i)
	engine.GET("/ip/echo", ipHandler.Echo)
}
