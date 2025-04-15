package acontext

import (
	"github.com/kataras/iris/v12"
	"strings"
)

func clientIP(ictx iris.Context) string {
	// X-Forwarded-For each agent will append its ip to tail and join with a comma. the first one is the client ip
	forwards := ictx.GetHeader("X-Forwarded-For")
	if forwards != "" {
		ips := strings.Split(forwards, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				return ip
			}
		}
	}
	// X-Real-IP from Nginx config `proxy_set_header X-Real-IP $remote_add;`
	ip := ictx.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	// RemoteAddr only available to client connect to this server directly (without any agent)
	return ictx.RemoteAddr()
}

// ClientIP get client real IP
func ClientIP(ictx iris.Context) string {
	ip := ictx.Values().GetString(CtxClientIP)
	if ip != "" {
		return ip
	}
	ip = clientIP(ictx)
	if ip == "" {
		return ""
	}
	ictx.Values().Set(CtxClientIP, ip)
	return ip
}
