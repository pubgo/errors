package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var logger zerolog.Logger

func init() {
	logger = log.With().Str("pkg", "errors").Logger()
}

func InitDebugLog() {
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}
