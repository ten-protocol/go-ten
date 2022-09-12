package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// setups some defaults like timestamp precision and logger
func init() { //nolint:gochecknoinits
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMilli, NoColor: true}).With().Timestamp().Logger()
	logger = logger.Level(WarnLevel)
}

// Errors will always show
const (
	InfoLevel     = zerolog.InfoLevel
	WarnLevel     = zerolog.WarnLevel
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
	Info("Log level set to: %s", level)
}

func OutputToFile(f *os.File) {
	Info("Log output set to: %s", f.Name())
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

func Warn(msg string, args ...interface{}) {
	logger.Warn().Msgf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.Info().Msgf(msg, args...)
}

func Panic(msg string, args ...interface{}) {
	logger.Panic().Msgf(msg, args...)
}

// ParseLevel returns a logging level given the string - defaults to info level if string not parseable
func ParseLevel(levelStr string) zerolog.Level {
	lvl, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		// we purposefully
		Error("Unable to parse log level: %s - defaulting to %s level", levelStr, zerolog.InfoLevel)
		lvl = zerolog.InfoLevel
	}
	return lvl
}
