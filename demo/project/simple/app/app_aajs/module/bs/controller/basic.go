package controller

import "github.com/kataras/iris/v12"

func (c *Controller) GetPing(ictx iris.Context) {
	defer ictx.Next()
	ictx.WriteString("PONG")
}
