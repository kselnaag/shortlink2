package log

import (
	"context"
	"fmt"
	"io"
	"os"
	T "shortlink2/internal/types"
	"strings"
	"sync"
	"time"
)

var _ T.ILog = (*LogFprintf)(nil)

type LogFprintf struct {
	loglvl    T.LogLevel
	host      string
	svc       string
	targets   []io.Writer
	batchTime time.Duration
	logbuf    []string
	mu        sync.Mutex
}

func NewLogFprintf(cfg T.ICfg, batchTime time.Duration, targets ...io.Writer) *LogFprintf {
	if len(targets) == 0 {
		targets = append(targets, os.Stderr)
	}

	host := cfg.GetVal(T.SL_HTTP_IP) + cfg.GetVal(T.SL_HTTP_PORT)
	svc := cfg.GetVal(T.SL_APP_NAME)
	var lvl T.LogLevel
	switch cfg.GetVal(T.SL_LOG_LEVEL) {
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
	case T.StrPanic:
		lvl = T.Panic
	case T.StrFatal:
		lvl = T.Fatal
	default:
		lvl = T.NoLog
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

func (l *LogFprintf) Start() func() {
	var wg sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())
	if l.batchTime != 0 {
		wg.Add(1)
		go func() {
			writeBatch := func() {
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
			for {
				select {
				case <-time.After(l.batchTime):
					writeBatch()
				case <-ctx.Done():
					writeBatch()
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
	if l.loglvl <= T.Trace {
		l.logMessage(T.StrTrace, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= T.Debug {
		l.logMessage(T.StrDebug, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= T.Info {
		l.logMessage(T.StrInfo, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= T.Warn {
		l.logMessage(T.StrWarn, l.host, l.svc, fmt.Sprintf(format, v...))
	}
}

func (l *LogFprintf) LogError(err error) {
	if l.loglvl <= T.Error {
		l.logMessage(T.StrError, l.host, l.svc, err.Error())
	}
}

func (l *LogFprintf) LogPanic(err error) {
	if l.loglvl <= T.Panic {
		l.logMessage(T.StrPanic, l.host, l.svc, err.Error())
	}
}

func (l *LogFprintf) LogFatal(err error) {
	if l.loglvl <= T.Fatal {
		l.logMessage(T.StrFatal, l.host, l.svc, err.Error())
	}
}
