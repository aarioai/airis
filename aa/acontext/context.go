package acontext

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type Context context.Context // global context with cancel

func Background() Context {
	return context.Background()
}
func TODO() Context {
	return context.TODO()
}

func WithCancel(parent Context) (Context, context.CancelFunc) {
	return context.WithCancel(parent)
}
func WithCancelCause(parent Context) (Context, context.CancelCauseFunc) {
	return context.WithCancelCause(parent)
}
func Cause(ctx Context) error {
	return context.Cause(ctx)
}
func AfterFunc(ctx Context, f func()) (stop func() bool) {
	return context.AfterFunc(ctx, f)
}
func WithoutCancel(parent Context) Context {
	return context.WithoutCancel(parent)
}
func WithDeadline(parent Context, d time.Time) (Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}
func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, context.CancelFunc) {
	return context.WithDeadlineCause(parent, d, cause)
}
func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, context.CancelFunc) {
	return context.WithTimeoutCause(parent, timeout, cause)
}
func WithValue(parent Context, key, val any) Context {
	return context.WithValue(parent, key, val)
}

func WithTrace(parent Context) interface{} {
	pc, file, line, _ := runtime.Caller(2)
	traceId := fmt.Sprintf("%s:%d.%s", file, line, runtime.FuncForPC(pc).Name())
	return context.WithValue(parent, CtxTraceId, traceId)
}
