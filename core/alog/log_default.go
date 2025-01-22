package alog

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/core/airis"
	"log"
	"strings"
	"sync"

	"github.com/aarioai/airis/core/ae"
)

type xlog struct {
}

var (
	xlogOnce     sync.Once
	xlogInstance *xlog
)

func NewDefaultLog() LogInterface {
	xlogOnce.Do(func() {
		xlogInstance = &xlog{}
	})
	return xlogInstance
}
func xlogHeader(ctx context.Context, caller string, level ErrorLevel) string {
	traceInfo := airis.TraceInfo(ctx)
	b := strings.Builder{}
	b.Grow(15 + len(traceInfo))
	b.WriteString(caller)
	b.WriteString(traceInfo)
	if level != ErrAll {
		b.WriteString(" [" + level.Name() + "] ")
	}
	return b.String()
}

func xprintf(ctx context.Context, level ErrorLevel, msg string, args ...any) {
	errlevel := errorlevel(ctx)
	if errlevel != ErrAll && errlevel&level == 0 {
		return
	}
	_, caller := ae.CallerMsg(msg, 1)
	head := xlogHeader(ctx, caller, level)
	msg = head + msg
	if len(args) == 0 {
		log.Println(msg)
	} else {
		log.Printf(msg+"\n", args...)
	}
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

func (l *xlog) Assert(ctx context.Context, condition bool, msg string, args ...any) {
	if condition {
		xprintf(ctx, Debug, msg, args...)
	}
}

func (l *xlog) Debug(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Debug, msg, args...)
}

func (l *xlog) Info(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Info, msg, args...)
}

func (l *xlog) Notice(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Notice, msg, args...)
}

func (l *xlog) Warn(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Warn, msg, args...)
}

func (l *xlog) Error(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Error, msg, args...)
}

func (l *xlog) Fatal(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Fatal, msg, args...)
}

func (l *xlog) Alert(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Alert, msg, args...)
}

func (l *xlog) Emerg(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Emerg, msg, args...)
}

func (l *xlog) Println(ctx context.Context, msg ...any) {
	log.Println(xlogHeader(ctx, ae.Caller(1), ErrAll), fmt.Sprint(msg...))
}

func (l *xlog) Trace(ctx context.Context) {
	traceInfo := airis.TraceInfo(ctx)
	log.Printf("[TRACE]%s %s\n", traceInfo, ae.Caller(1))
}
