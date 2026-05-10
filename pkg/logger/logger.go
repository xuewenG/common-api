package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/xuewenG/common-api/pkg/config"
)

func NewLogger(i do.Injector) (*zerolog.Logger, error) {
	config := do.MustInvoke[*config.Config](i)

	var level = zerolog.InfoLevel
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	err := level.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
		},
	).With().Timestamp().Logger()

	return &logger, nil
}
