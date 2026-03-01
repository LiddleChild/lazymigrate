package log

import (
	"time"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

type Message struct {
	Time    time.Time
	Level   LogLevel
	Message string
}
