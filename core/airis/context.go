package airis

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"runtime"
	"strings"
)

func GetTraceId(ictx iris.Context) string {
	traceId := ictx.GetHeader("X-Request-Id")
	if traceId != "" {
		return traceId
	}
	return "NO_TRACE_ID"
}

// Context
// 后面可以使用 ictx.Request().Context() 直接访问
func Context(ictx iris.Context) context.Context {
	traceId := GetTraceId(ictx)
	remoteAddr := ictx.RemoteAddr()
	user := ictx.Values().GetString(CtxRemoteUser)
	ctx := context.WithValue(ictx.Request().Context(), CtxTraceId, traceId)
	ctx = context.WithValue(ctx, CtxRemoteAddr, remoteAddr)
	ctx = context.WithValue(ctx, CtxRemoteUser, user)
	ictx.ResetRequest(ictx.Request().WithContext(ctx))
	return ctx
}

// JobContext 后台任务
// Job 后台任务；Task 往往需要交互的任务
func JobContext(parent context.Context) context.Context {
	pc, file, line, _ := runtime.Caller(2)
	// 这个和log 的caller并不相同。这里仅表示context位置，并非日志caller调用位置
	traceId := fmt.Sprintf("%s:%d.%s", file, line, runtime.FuncForPC(pc).Name())
	return context.WithValue(parent, CtxTraceId, traceId)
}

func TraceInfo(ctx context.Context) string {
	traceId, _ := ctx.Value(CtxTraceId).(string)
	remoteAddr, _ := ctx.Value(CtxRemoteAddr).(string)
	user, _ := ctx.Value(CtxRemoteUser).(string)
	var s strings.Builder
	s.Grow(len(traceId) + len(remoteAddr) + len(user))
	if traceId != "" {
		s.WriteString("trace_id:" + traceId)
	}
	if remoteAddr != "" {
		if s.Len() > 0 {
			s.WriteString(", ")
		}
		s.WriteString(" remote_addr:")
	}
	if user != "" {
		if s.Len() > 0 {
			s.WriteString(", ")
		}
		s.WriteString(" user:" + user)
	}
	if s.Len() == 0 {
		return ""
	}
	return " {" + s.String() + "}"
}
