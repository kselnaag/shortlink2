/*
	Fprintf log module:

- universal DI interface (see types/log.go)
- structured log to JSON
- manual key positions into JSON object
- 8 Log levels (trace, debug, info, warn, error, panic, fatal, nolog)
- stack trace in Panic and Fatal log messages (os.Exit(1) on Fatal)
- multi-target message sending with io.Writer interface (if empty - os.Stderr)
- log batching with timeout (if 0 - no batching)
*/
package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	T "shortlink2/internal/types"
	"strings"
	"sync"
	"time"
)

var _ T.ILog = (*LogFprintf)(nil)

const (
	StrTrace = "TRACE"
	StrDebug = "DEBUG"
	StrInfo  = "INFO"
	StrWarn  = "WARN"
	StrError = "ERROR"
	StrPanic = "PANIC"
	StrFatal = "FATAL"
	StrNoLog = "NOLOG"
)

type LogLevel int8

const (
	Trace LogLevel = iota
	Debug
	Info
	Warn
	Error
	Panic
	Fatal
	NoLog
)

type LogFprintf struct {
	loglvl    LogLevel
	host      string
	svc       string
	targets   []io.Writer
	batchTime time.Duration
	logbuf    []string
	mu        sync.Mutex
}

func NewLogFprintf(cfg T.ICfg, batchTime time.Duration, targets ...io.Writer) *LogFprintf {
	debug.SetTraceback("all")
	if len(targets) == 0 {
		targets = append(targets, os.Stderr)
	}
	host := cfg.GetVal(T.SL_HTTP_IP) + cfg.GetVal(T.SL_HTTP_PORT)
	svc := cfg.GetVal(T.SL_APP_NAME)
	var lvl LogLevel
	switch cfg.GetVal(T.SL_LOG_LEVEL) {
	case StrTrace:
		lvl = Trace
	case StrDebug:
		lvl = Debug
	case StrInfo:
		lvl = Info
	case StrWarn:
		lvl = Warn
	case StrError:
		lvl = Error
	case StrPanic:
		lvl = Panic
	case StrFatal:
		lvl = Fatal
	default:
		lvl = NoLog
	}
	return &LogFprintf{
		loglvl:    lvl,
		host:      host,
		svc:       svc,
		targets:   targets,
		batchTime: batchTime,
		logbuf:    make([]string, 0, 100),
	}
}

func (l *LogFprintf) writeBatch() {
	l.mu.Lock()
	if len(l.logbuf) != 0 {
		batchstr := strings.Join(l.logbuf, "")
		for _, point := range l.targets {
			fmt.Fprintf(point, batchstr)
		}
		l.logbuf = l.logbuf[:0]
	}
	l.mu.Unlock()
}

func (l *LogFprintf) Start() func() {
	var wg sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())
	if l.batchTime != 0 {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-time.After(l.batchTime):
					l.writeBatch()
				case <-ctx.Done():
					l.writeBatch()
					wg.Done()
					return
				}
			}
		}()
	}
	return func() {
		ctxCancel()
		wg.Wait()
	}
}

func (l *LogFprintf) logMessage(lvl, host, svc, mess string) {
	timenow := time.Now().Format(time.RFC3339Nano)
	formatstr := `{"T":"%s","L":"%s","H":"%s","S":"%s","M":"%s"}` + "\n"
	if l.batchTime == 0 {
		for _, point := range l.targets {
			fmt.Fprintf(point, formatstr, timenow, lvl, host, svc, mess)
		}
	} else {
		l.mu.Lock()
		l.logbuf = append(l.logbuf, fmt.Sprintf(formatstr, timenow, lvl, host, svc, mess))
		l.mu.Unlock()
	}
}

func (l *LogFprintf) LogTrace(format string, v ...any) {
	if l.loglvl <= Trace {
		l.logMessage(StrTrace, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= Debug {
		l.logMessage(StrDebug, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= Info {
		l.logMessage(StrInfo, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= Warn {
		l.logMessage(StrWarn, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogError(err error) {
	if l.loglvl <= Error {
		l.logMessage(StrError, l.host, l.svc, err.Error())
	}
}

func (l *LogFprintf) LogPanic(err error) {
	if l.loglvl <= Panic {
		l.logMessage(StrPanic, l.host, l.svc, fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()))
	}
}

func (l *LogFprintf) LogFatal(err error) {
	if l.loglvl <= Fatal {
		l.logMessage(StrFatal, l.host, l.svc, fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()))
		if l.batchTime != 0 {
			l.writeBatch()
		}
		os.Exit(1)
	}
}
