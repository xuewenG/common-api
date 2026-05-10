package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/xuewenG/common-api/pkg/bark"
	"github.com/xuewenG/common-api/pkg/config"
	"github.com/xuewenG/common-api/pkg/engine"
	"github.com/xuewenG/common-api/pkg/handler"
	"github.com/xuewenG/common-api/pkg/logger"
	"github.com/xuewenG/common-api/pkg/router"
)

var mode string
var version string
var commitId string

func main() {
	if mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	injector := do.New()

	do.Provide(injector, config.NewConfig)
	do.Provide(injector, logger.NewLogger)

	cfg := do.MustInvoke[*config.Config](injector)
	logger := do.MustInvoke[*zerolog.Logger](injector)

	logger.Info().Msgf(
		"Server is starting:\nmode: %s\nversion: %s\ncommitId: %s\nport: %s\n",
		mode,
		version,
		commitId,
		cfg.Port,
	)

	do.Provide(injector, bark.NewBarkClient)

	do.Provide(injector, engine.NewEngine)
	do.Provide(injector, router.NewRouter)
	do.Provide(injector, handler.NewHealthHandler)
	do.Provide(injector, handler.NewHolidayHandler)
	do.Provide(injector, handler.NewBiliRecHandler)
	do.Provide(injector, handler.NewIPHandler)

	ginEngine := do.MustInvoke[*gin.Engine](injector)
	r := do.MustInvoke[*router.Router](injector)

	r.Bind(ginEngine)
	ginEngine.Run(fmt.Sprintf(":%s", cfg.Port))
}
