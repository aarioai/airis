package acontext

import (
	"context"
	"strings"
)

type requestInfo struct {
	TraceID    string
	RemoteAddr string
	RemoteUser string
}

type ctxKey uint8

const keyRequestInfo ctxKey = 1

// Bind copies the client address, user, and request id onto the request context.
// A single context value holds the metadata to avoid repeated WithValue allocations.
func Bind(c RequestCarrier) context.Context {
	remoteAddr := ClientIP(c)
	if remoteAddr == "" {
		remoteAddr = c.RemoteAddr()
	}

	info := &requestInfo{
		TraceID:    c.ValueString(CtxRequestID),
		RemoteAddr: remoteAddr,
		RemoteUser: c.ValueString(CtxRemoteUser),
	}
	ctx := context.WithValue(c.Request().Context(), keyRequestInfo, info)
	c.ResetRequest(c.Request().WithContext(ctx))
	return ctx
}

func infoFrom(ctx context.Context) *requestInfo {
	if ctx == nil {
		return nil
	}
	if info, ok := ctx.Value(keyRequestInfo).(*requestInfo); ok {
		return info
	}
	return nil
}

// TraceID returns the request id stored by Bind, or the legacy context key.
func TraceID(ctx context.Context) string {
	if info := infoFrom(ctx); info != nil {
		return info.TraceID
	}
	if ctx == nil {
		return ""
	}
	s, _ := ctx.Value(CtxTraceId).(string)
	return s
}

// RemoteAddrValue returns the remote address stored by Bind.
func RemoteAddrValue(ctx context.Context) string {
	if info := infoFrom(ctx); info != nil {
		return info.RemoteAddr
	}
	if ctx == nil {
		return ""
	}
	s, _ := ctx.Value(CtxRemoteAddr).(string)
	return s
}

// RemoteUserValue returns the remote user stored by Bind.
func RemoteUserValue(ctx context.Context) string {
	if info := infoFrom(ctx); info != nil {
		return info.RemoteUser
	}
	if ctx == nil {
		return ""
	}
	s, _ := ctx.Value(CtxRemoteUser).(string)
	return s
}

func TraceInfo(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := TraceID(ctx)
	remoteAddr := RemoteAddrValue(ctx)
	user := RemoteUserValue(ctx)
	if traceId == "" && remoteAddr == "" && user == "" {
		return ""
	}
	var s strings.Builder
	s.Grow(len(traceId) + len(remoteAddr) + len(user) + 40)
	writeTrace := func(key, val string) {
		if val == "" {
			return
		}
		if s.Len() > 0 {
			s.WriteString(", ")
		}
		s.WriteString(key)
		s.WriteByte(':')
		s.WriteString(val)
	}
	writeTrace("trace_id", traceId)
	writeTrace("remote_addr", remoteAddr)
	writeTrace("user", user)
	return " {" + s.String() + "}"
}
