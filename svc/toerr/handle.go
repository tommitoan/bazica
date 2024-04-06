package bazica

import (
	"github.com/rollbar/rollbar-go"
	"github.com/rs/zerolog/log"
)

// HandleInternalError handle errors that occur from non api facing functions
// For example, from scheduled or cron jobs , goroutines etc./*
func HandleInternalError(err error, fields map[string]any) {
	rollbar.Critical(err)
	log.Error().Err(err).Fields(fields).Msg("caught internal error")
}
