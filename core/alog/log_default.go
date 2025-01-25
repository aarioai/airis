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
func xlogHeader(ctx context.Context, level ErrorLevel, caller string) string {
	traceInfo := airis.TraceInfo(ctx)
	b := strings.Builder{}
	b.Grow(15 + len(traceInfo))

	if level != ErrAll {
		b.WriteString("[")
		b.WriteString(level.Name())
		b.WriteString("] ")
	}

	b.WriteString(caller)
	b.WriteString(traceInfo)

	return b.String()
}

func xprintf(ctx context.Context, level ErrorLevel, caller string, msg string, args ...any) {
	errlevel := errorlevel(ctx)
	if errlevel != ErrAll && errlevel&level == 0 {
		return
	}

	head := xlogHeader(ctx, level, caller)
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
func (l *xlog) E(ctx context.Context, e *ae.Error, msg ...any) {
	s := e.Text()
	if len(msg) > 0 {
		s = fmt.Sprint(msg...)
	}
	xprintf(ctx, Error, e.Caller, s)
}
func (l *xlog) Debug(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Debug, utils.Caller(1), msg, args...)
}

func (l *xlog) Info(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Info, utils.Caller(1), msg, args...)
}

func (l *xlog) Notice(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Notice, utils.Caller(1), msg, args...)
}

func (l *xlog) Warn(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Warn, utils.Caller(1), msg, args...)
}

func (l *xlog) Error(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Error, utils.Caller(1), msg, args...)
}

func (l *xlog) Fatal(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Fatal, utils.Caller(1), msg, args...)
}

func (l *xlog) Alert(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Alert, utils.Caller(1), msg, args...)
}

func (l *xlog) Emerg(ctx context.Context, msg string, args ...any) {
	xprintf(ctx, Emerg, utils.Caller(1), msg, args...)
}

func (l *xlog) Println(ctx context.Context, msg ...any) {
	log.Println(xlogHeader(ctx, ErrAll, utils.Caller(1)), fmt.Sprint(msg...))
}

func (l *xlog) Trace(ctx context.Context) {
	traceInfo := airis.TraceInfo(ctx)
	log.Printf("[TRACE]%s %s\n", traceInfo, utils.Caller(1))
}
