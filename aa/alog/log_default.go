package alog

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/pkg/utils"
	"log"
	"strings"
	"sync"

	"github.com/aarioai/airis/aa/ae"
)

type xlog struct {
	Level ErrorLevel
}

var (
	xlogOnce     sync.Once
	xlogInstance *xlog
)

func NewDefaultLog(level ErrorLevel) LogInterface {
	xlogOnce.Do(func() {
		xlogInstance = &xlog{
			Level: level,
		}
	})
	return xlogInstance
}
func packMessage(ctx context.Context, level ErrorLevel, caller, msg string) string {
	traceInfo := acontext.TraceInfo(ctx)
	b := strings.Builder{}
	b.Grow(15 + len(traceInfo))

	if level != LevelErrAll {
		b.WriteString("[")
		b.WriteString(level.Name())
		b.WriteString("] ")
	}

	b.WriteString(caller)
	b.WriteByte(' ')
	b.WriteString(msg)

	b.WriteString(traceInfo)
	return b.String()
}

func (l *xlog) print(ctx context.Context, level ErrorLevel, caller string, msg string) {
	if level > l.Level {
		return
	}
	log.Println(packMessage(ctx, level, caller, msg))
}
func (l *xlog) New(prefix string, f func(context.Context, string, ...any), suffix ...string) func(context.Context, string, ...any) {
	var s string
	if len(suffix) > 0 {
		s = " " + suffix[0]
	}
	if prefix != "" {
		prefix += " "
	}
	return func(ctx context.Context, msg string, args ...any) {
		f(ctx, prefix+msg+s, args...)
	}
}
func (l *xlog) Debug(ctx context.Context, msg string) {
	l.print(ctx, LevelDebug, utils.Caller(1), msg)
}
func (l *xlog) Debugf(ctx context.Context, msg string, args ...any) {
	l.print(ctx, LevelDebug, utils.Caller(1), fmt.Sprintf(msg, args...))
}

func (l *xlog) Info(ctx context.Context, msg string) {
	l.print(ctx, LevelInfo, utils.Caller(1), msg)
}

func (l *xlog) InfoF(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelInfo, utils.Caller(1), fmt.Sprintf(format, args...))
}
func (l *xlog) Notice(ctx context.Context, msg string) {
	l.print(ctx, LevelNotice, utils.Caller(1), msg)
}
func (l *xlog) Noticef(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelNotice, utils.Caller(1), fmt.Sprintf(format, args...))
}
func (l *xlog) Warn(ctx context.Context, format string) {
	l.print(ctx, LevelWarn, utils.Caller(1), format)
}
func (l *xlog) Warnf(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelWarn, utils.Caller(1), fmt.Sprintf(format, args...))
}

func (l *xlog) Error(ctx context.Context, msg string) {
	l.print(ctx, LevelError, utils.Caller(1), msg)
}
func (l *xlog) Errorf(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelError, utils.Caller(1), fmt.Sprintf(format, args...))
}
func (l *xlog) E(ctx context.Context, e *ae.Error, msg ...any) {
	s := e.String()
	if len(msg) > 0 {
		s = fmt.Sprint(msg...)
	}
	l.print(ctx, LevelError, e.Caller, s)
}

func (l *xlog) Fatal(ctx context.Context, msg string) {
	l.print(ctx, LevelFatal, utils.Caller(1), msg)
}

func (l *xlog) Fatalf(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelFatal, utils.Caller(1), fmt.Sprintf(format, args...))
}

func (l *xlog) Alert(ctx context.Context, msg string) {
	l.print(ctx, LevelAlert, utils.Caller(1), msg)
}

func (l *xlog) Alertf(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelAlert, utils.Caller(1), fmt.Sprintf(format, args...))
}

func (l *xlog) Emerg(ctx context.Context, msg string) {
	l.print(ctx, LevelEmerg, utils.Caller(1), msg)
}

func (l *xlog) Emergf(ctx context.Context, format string, args ...any) {
	l.print(ctx, LevelEmerg, utils.Caller(1), fmt.Sprintf(format, args...))
}

func (l *xlog) Print(ctx context.Context, msg string) {
	log.Println(packMessage(ctx, LevelErrAll, utils.Caller(1), msg))
}

func (l *xlog) Printf(ctx context.Context, format string, args ...any) {
	log.Println(packMessage(ctx, LevelErrAll, utils.Caller(1), fmt.Sprintf(format, args...)))
}

func (l *xlog) Trace(ctx context.Context) {
	traceInfo := acontext.TraceInfo(ctx)
	log.Printf("[TRACE]%s %s\n", traceInfo, utils.Caller(1))
}
