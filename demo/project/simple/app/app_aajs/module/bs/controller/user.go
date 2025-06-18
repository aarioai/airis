package controller

import (
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/kataras/iris/v12"
)

func (c *Controller) GetUsers(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	paging := r.Paging()

	resp.TryWrite(c.s.Users(paging))
}

func (c *Controller) GetUser(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	if e := resp.FirstError(e0); e != nil {
		return
	}

	resp.TryWrite(c.s.User(uid))
}

func (c *Controller) PostUser(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	username, e0 := r.BodyString("username", `^\w+$`, false)
	age, e1 := r.BodyInt("age", false)
	if e := resp.FirstError(e0, e1); e != nil {
		return
	}

	resp.TryWrite(c.s.PostUser(username, age))
}

func (c *Controller) PutUser(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	username, e1 := r.BodyString("username", `^\w+$`)
	age, e2 := r.BodyInt("age", false)
	if e := resp.FirstError(e0, e1, e2); e != nil {
		return
	}

	resp.WriteE(c.s.PutUser(uid, username, age))
}

func (c *Controller) PatchUser(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	age, e1 := r.BodyInt("age")
	if e := resp.FirstError(e0, e1); e != nil {
		return
	}

	resp.WriteE(c.s.PatchUser(uid, age))
}

func (c *Controller) DeleteUser(ictx iris.Context) {
	r, resp, _ := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	if e := resp.FirstError(e0); e != nil {
		return
	}

	resp.WriteE(c.s.DeleteUser(uid))
}
