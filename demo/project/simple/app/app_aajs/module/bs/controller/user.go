package controller

import (
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/kataras/iris/v12"
)

func (c *Controller) GetUsers(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.CloseWith(r)
}
