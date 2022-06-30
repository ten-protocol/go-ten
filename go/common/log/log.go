package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// setups some defaults like timestamp precision and logger
func init() { //nolint:gochecknoinits
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(TraceLevel)
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMilli}).With().Timestamp().Logger()
}

// Errors will always show
const (
	InfoLevel     = zerolog.InfoLevel
	DebugLevel    = zerolog.DebugLevel
	TraceLevel    = zerolog.TraceLevel
	DisabledLevel = zerolog.Disabled // Nothing will show at this level not even errors
)

// singleton usage of the logger
var (
	logger = zerolog.Logger{}
)

func SetLogLevel(level zerolog.Level) {
	logger = logger.Level(level)
}

func OutputToFile(f *os.File) {
	logger = zerolog.New(zerolog.ConsoleWriter{Out: f, TimeFormat: time.StampMilli, NoColor: true}).With().Timestamp().Logger()
}

func Error(msg string, args ...interface{}) {
	logger.Error().Msgf(msg, args...)
}

func Trace(msg string, args ...interface{}) {
	logger.Trace().Msgf(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Debug().Msgf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Info().Msgf(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	logger.Panic().Msgf(msg, args...)
}
