package services

import (
	"fmt"
)

type LogLevel int

const (
	CriticalLevel LogLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (l LogLevel) String() string {
	switch l {
	case CriticalLevel:
		return "CRITICAL"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "TRACE"
	default:
		return "INFO"
	}
}

func Audit(services *Services, level LogLevel, msg string, params ...any) {
	safeParams := make([]any, len(params))
	for i, p := range params {
		if p == nil {
			safeParams[i] = "<nil>"
		} else {
			safeParams[i] = p
		}
	}

	formattedMsg := fmt.Sprintf(msg, safeParams...)

	switch level {
	case CriticalLevel:
		services.Logger().Crit(formattedMsg)
	case ErrorLevel:
		services.Logger().Error(formattedMsg)
	case WarnLevel:
		services.Logger().Warn(formattedMsg)
	case InfoLevel:
		services.Logger().Info(formattedMsg)
	case DebugLevel:
		services.Logger().Debug(formattedMsg)
	case TraceLevel:
		services.Logger().Trace(formattedMsg)
	default:
		services.Logger().Info(formattedMsg)
	}
}
