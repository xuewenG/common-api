package handler

import (
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/xuewenG/common-api/pkg/config"
)

type Handler struct {
	config *config.Config
	logger *zerolog.Logger
}

func newHandler(i do.Injector) (*Handler, error) {
	return &Handler{
		config: do.MustInvoke[*config.Config](i),
		logger: do.MustInvoke[*zerolog.Logger](i),
	}, nil
}
