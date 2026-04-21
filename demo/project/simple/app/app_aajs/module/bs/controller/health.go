package controller

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/kataras/iris/v12"
)

func (c *Controller) Ping(ictx iris.Context) {
	defer ictx.Next()
	ictx.WriteString("PONG")
}

func (c *Controller) PingRedis(ictx iris.Context) {
	defer ictx.Next()
	_, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	ictx.WriteString(c.h.Ping(ctx))
}

func (c *Controller) PingMySQL(ictx iris.Context) {
	_, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	ictx.WriteString(c.m.Ping(ctx))
}

func (c *Controller) Health(ictx iris.Context) {
	_, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	// 建议将API以http 状态码输出错误，还是以body code 方式输出错误的选择权交给客户端决定，而不是由服务端决定。
	// 见 https://github.com/aarioai/airis/blob/main/README_%E8%AF%B4%E6%98%8E.md
	//resp.ErrorAsStatus = true
	if acontext.ServHealth().IsReady() {
		resp.WriteOK()
		return
	}
	resp.WriteCode(ae.ServiceUnavailable)
}

func (c *Controller) HealthReady(ictx iris.Context) {
	_, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	// 建议将API以http 状态码输出错误，还是以body code 方式输出错误的选择权交给客户端决定，而不是由服务端决定。
	// 见 https://github.com/aarioai/airis/blob/main/README_%E8%AF%B4%E6%98%8E.md
	//resp.ErrorAsStatus = true
	if acontext.ServHealth().IsReady() {
		resp.WriteOK()
		return
	}
	resp.WriteCode(ae.ServiceUnavailable)
}
