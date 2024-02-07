package types

type ILog interface {
	LogTrace(format string, v ...any)
	LogDebug(format string, v ...any)
	LogInfo(format string, v ...any)
	LogWarn(format string, v ...any)
	LogError(err error, format string, v ...any)
	LogFatal(err error, format string, v ...any)
	LogPanic(err error, format string, v ...any)
}

const (
	StrTrace   = "TRACE"
	StrDebug   = "DEBUG"
	StrInfo    = "INFO"
	StrWarn    = "WARN"
	StrError   = "ERROR"
	StrFatal   = "FATAL"
	StrPanic   = "PANIC"
	StrLOGSOFF = "LOGSOFF"
)

type LogLevel int8

const (
	Trace LogLevel = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Panic
	Logsoff
)
