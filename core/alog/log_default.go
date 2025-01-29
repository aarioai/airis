package alog

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/core/airis"
	"github.com/aarioai/airis/pkg/utils"
	"log"
	"strings"
	"sync"

	"github.com/aarioai/airis/core/ae"
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
func format(ctx context.Context, level ErrorLevel, caller, msg string) string {
	traceInfo := airis.TraceInfo(ctx)
	b := strings.Builder{}
	b.Grow(15 + len(traceInfo))

	if level != ErrAll {
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

func (l *xlog) print(ctx context.Context, level ErrorLevel, caller string, msg string, args ...any) {
	if level > l.Level {
		return
	}
	log.Println(format(ctx, level, caller, fmt.Sprintf(msg, args...)))
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

func (l *xlog) Debug(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Debug, utils.Caller(1), msg, args...)
}

func (l *xlog) Info(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Info, utils.Caller(1), msg, args...)
}

func (l *xlog) Notice(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Notice, utils.Caller(1), msg, args...)
}

func (l *xlog) Warn(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Warn, utils.Caller(1), msg, args...)
}

func (l *xlog) Error(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Error, utils.Caller(1), msg, args...)
}
func (l *xlog) E(ctx context.Context, e *ae.Error, msg ...any) {
	s := e.String()
	if len(msg) > 0 {
		s = fmt.Sprint(msg...)
	}
	l.print(ctx, Error, e.Caller, s)
}

func (l *xlog) Fatal(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Fatal, utils.Caller(1), msg, args...)
}

func (l *xlog) Alert(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Alert, utils.Caller(1), msg, args...)
}

func (l *xlog) Emerg(ctx context.Context, msg string, args ...any) {
	l.print(ctx, Emerg, utils.Caller(1), msg, args...)
}

func (l *xlog) Print(ctx context.Context, msg string, args ...any) {
	log.Println(format(ctx, ErrAll, utils.Caller(1), fmt.Sprintf(msg, args...)))
}

func (l *xlog) Trace(ctx context.Context) {
	traceInfo := airis.TraceInfo(ctx)
	log.Printf("[TRACE]%s %s\n", traceInfo, utils.Caller(1))
}
