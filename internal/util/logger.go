package util

import (
	"os"

	"github.com/rs/zerolog"

	"fliqt/config"
)

func NewLogger(cfg *config.Config) *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if cfg.PrettyLog {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		return &logger
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &logger
}
