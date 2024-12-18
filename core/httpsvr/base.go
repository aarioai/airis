package httpsvr

import (
	"context"
	"github.com/aarioai/airis/core/httpsvr/request"
	"github.com/aarioai/airis/core/httpsvr/response"
	"github.com/kataras/iris/v12"
)

// 读取json buffer 的时候，会清空掉 r.Body，所以这个使用一次；
func New(ictx iris.Context) (*request.Request, *response.Writer, context.Context) {
	r := request.New(ictx)
	w := response.NewWriter(ictx, r)
	return r, w, r.Context()
}
