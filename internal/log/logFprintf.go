package log

import (
	"fmt"
	"os"
	T "shortlink2/internal/types"
	"time"
)

var _ T.ILog = (*LogFprintf)(nil)

type LogFprintf struct {
	loglvl T.LogLevel
	host   string
	svc    string
}

func NewLogFprintf(cfg *T.CfgEnv) *LogFprintf {
	host := cfg.SL_HTTP_IP + cfg.SL_HTTP_PORT
	svc := cfg.SL_APP_NAME + cfg.SL_APP_PROTOCS
	var lvl T.LogLevel
	switch cfg.SL_LOG_LEVEL {
	case T.StrTrace:
		lvl = T.Trace
	case T.StrDebug:
		lvl = T.Debug
	case T.StrInfo:
		lvl = T.Info
	case T.StrWarn:
		lvl = T.Warn
	case T.StrError:
		lvl = T.Error
	case T.StrFatal:
		lvl = T.Fatal
	case T.StrPanic:
		lvl = T.Panic
	default:
		lvl = T.Logsoff
	}
	return &LogFprintf{
		loglvl: lvl,
		host:   host,
		svc:    svc,
	}
}

func logMessage(lvl, host, svc, err, mess string) {
	timenow := time.Now().Format(time.RFC3339Nano)
	fmt.Fprintf(os.Stderr, "{\"T\":\"%s\",\"L\":\"%s\",\"H\":\"%s\",\"S\":\"%s\",\"M\":\"%s\",\"E\":\"%s\"}\n", timenow, lvl, host, svc, mess, err)
}

func (l *LogFprintf) LogTrace(format string, v ...any) {
	if l.loglvl <= T.Trace {
		logMessage(T.StrTrace, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= T.Debug {
		logMessage(T.StrDebug, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= T.Info {
		logMessage(T.StrInfo, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= T.Warn {
		logMessage(T.StrWarn, l.host, l.svc, "", fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogError(err error, format string, v ...any) {
	if l.loglvl <= T.Error {
		logMessage(T.StrError, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogFatal(err error, format string, v ...any) {
	if l.loglvl <= T.Fatal {
		logMessage(T.StrFatal, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogPanic(err error, format string, v ...any) {
	if l.loglvl <= T.Panic {
		logMessage(T.StrPanic, l.host, l.svc, err.Error(), fmt.Sprintf(format, v...))
	}
}
