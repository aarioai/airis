package acontext

import (
	"context"
	"github.com/kataras/iris/v12"
	"strings"
)

// Context
// 后面可以使用 ictx.Request().Context() 直接访问
func FromIris(ictx iris.Context) context.Context {
	// Nginx proxy_set_header X-Request-Id $request_id;
	traceId := ictx.GetHeader("X-Request-Id")
	// Nginx proxy_set_header X-Remote-Addr $remote_addr;
	remoteAddr := ictx.GetHeader("X-Remote-Addr")
	if remoteAddr == "" {
		remoteAddr = ictx.RemoteAddr()
	}

	user := ictx.Values().GetString(CtxRemoteUser)
	ctx := context.WithValue(ictx.Request().Context(), CtxTraceId, traceId)
	ctx = context.WithValue(ctx, CtxRemoteAddr, remoteAddr)
	ctx = context.WithValue(ctx, CtxRemoteUser, user)
	ictx.ResetRequest(ictx.Request().WithContext(ctx))
	return ctx
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
