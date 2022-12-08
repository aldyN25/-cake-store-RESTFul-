package infra

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

func newLogger() zerolog.Logger {

	buildInfo, _ := debug.ReadBuildInfo()

	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.Main.Version).
		Logger()

}
