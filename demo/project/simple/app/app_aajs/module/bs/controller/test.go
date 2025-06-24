package controller

import (
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/kataras/iris/v12"
)

func (c *Controller) HeadRestful(ictx iris.Context) {
	defer ictx.Next()
	ictx.StatusCode(200)
}
func (c *Controller) PostRestful(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	say, _ := r.BodyString("say", false)

	resp.Write(map[string]string{"say": say})
}

func (c *Controller) PutRestful(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	id, e0 := r.QueryInt("id")
	say, _ := r.BodyString("say", false)
	if e := resp.FirstError(e0); e != nil {
		return
	}

	resp.Write(map[string]any{
		"id":  id,
		"say": say,
	})
}

func (c *Controller) PatchRestful(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	num, _ := r.BodyInt("num", false)

	resp.Write(map[string]int{"num": num})
}

func (c *Controller) GetRestful(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()

	response, _ := r.QueryString("response", false)
	hello, _ := r.QueryString("hello", false)

	resp.Write(map[string]string{
		"response": response,
		"hello":    hello,
	})
}

func (c *Controller) DeleteRestful(ictx iris.Context) {
	_, resp, _ := httpsvr.New(ictx)
	defer resp.Next()

	resp.WriteOK()
}
