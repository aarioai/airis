package acontext

import "strings"

func clientIP(ctx RequestCarrier) string {
	// X-Forwarded-For: each agent appends its address. The first one is the client.
	forwards := ctx.GetHeader("X-Forwarded-For")
	if forwards != "" {
		ips := strings.Split(forwards, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				return ip
			}
		}
	}
	// X-Real-IP from Nginx `proxy_set_header X-Real-IP $remote_addr`.
	ip := ctx.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	return ctx.RemoteAddr()
}

// ClientIP returns the client address and caches it on the request.
func ClientIP(ctx RequestCarrier) string {
	ip := ctx.ValueString(CtxClientIP)
	if ip != "" {
		return ip
	}
	ip = clientIP(ctx)
	if ip == "" {
		return ""
	}
	ctx.SetValue(CtxClientIP, ip)
	return ip
}
