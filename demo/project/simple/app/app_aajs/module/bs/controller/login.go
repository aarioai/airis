package controller

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/kataras/iris/v12"
)

func (c *Controller) PostLogin(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	account, e1 := r.BodyString("account", `^\d+$`)
	password, e2 := r.BodyString("password")
	state, _ := r.BodyString("state")
	if e := resp.FirstError(e1, e2); e != nil {
		return
	}
	resp.TryWrite(c.s.Login(ctx, account, password, state))
}

func (c *Controller) HeadUserToken(ictx iris.Context) {
	defer ictx.Next()
	ictx.StatusCode(iris.StatusOK)
}

func (c *Controller) GrantUserToken(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	grantType, e0 := r.BodyString("grant_type", `^(authorization_code|refresh_token)$`)
	code, e1 := r.BodyString("code")
	if e := ae.First(e0, e1); e != nil {
		return
	}
	switch grantType {
	case "authorization_code":
	case "refresh_token":
		resp.TryWrite(c.s.RefreshUserToken(ctx, code))
		return
	}
	resp.WriteE(ae.NewBadParam("grant_type"))
}
