package controller

import (
	"github.com/aarioai/airis/aa/atype/aenum"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/kataras/iris/v12"
)

func (c *Controller) GetUsersWithPaging(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	paging := r.Paging()
	sex, e0 := r.QuerySex("sex", false)
	if e := resp.FirstError(e0); e != nil {
		return
	}

	if sex == aenum.NilSex {
		resp.TryWrite(c.s.Users(ctx, paging))
		return
	}

	resp.TryWrite(c.s.QueryUsersBySex(ctx, sex, paging))

}

func (c *Controller) GetUsers(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	uids, e0 := r.QueryUint64s("uid", true, false)
	if e := resp.FirstError(e0); e != nil {
		return
	}

	resp.TryWrite(c.s.QueryUsers(ctx, uids))
}

func (c *Controller) GetUser(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	if e := resp.FirstError(e0); e != nil {
		return
	}

	resp.TryWrite(c.s.User(ctx, uid))
}

func (c *Controller) PostUser(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	username, e0 := r.BodyString("username", `^\w+$`, false)
	age, e1 := r.BodyInt("age", false)
	if e := resp.FirstError(e0, e1); e != nil {
		return
	}

	resp.TryWrite(c.s.PostUser(ctx, username, age))
}

func (c *Controller) PutUser(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	username, e1 := r.BodyString("username", `^\w+$`)
	age, e2 := r.BodyInt("age", false)
	if e := resp.FirstError(e0, e1, e2); e != nil {
		return
	}

	resp.WriteE(c.s.PutUser(ctx, uid, username, age))
}

func (c *Controller) PatchUser(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	age, e1 := r.BodyInt("age")
	if e := resp.FirstError(e0, e1); e != nil {
		return
	}

	resp.WriteE(c.s.PatchUser(ctx, uid, age))
}

func (c *Controller) DeleteUser(ictx iris.Context) {
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.Next()
	uid, e0 := r.QueryUint64("uid")
	if e := resp.FirstError(e0); e != nil {
		return
	}

	resp.WriteE(c.s.DeleteUser(ctx, uid))
}
