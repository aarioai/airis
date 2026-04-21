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
	if acontext.ServHealth().IsReady() {
		resp.WriteOK()
		return
	}
	//resp.ErrorAsStatus = true
	resp.StatusCode(ae.ServiceUnavailable)
}

func (c *Controller) HealthReady(ictx iris.Context) {
	_, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	if acontext.ServHealth().IsReady() {
		resp.WriteOK()
		return
	}
	//resp.ErrorAsStatus = true
	resp.StatusCode(ae.ServiceUnavailable)
}
