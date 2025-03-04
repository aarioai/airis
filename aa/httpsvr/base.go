package httpsvr

import (
	"context"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/kataras/iris/v12"
)

// 读取json buffer 的时候，会清空掉 r.Body，所以这个使用一次；
func New(ictx iris.Context, multiSizes ...int64) (*request.Request, *response.Writer, context.Context) {
	r := request.New(ictx)
	if len(multiSizes) > 0 {
		r.SetMaxMultipartSize(multiSizes[0])
		if len(multiSizes) > 1 {
			r.SetMaxFormSize(multiSizes[1])
		}
	}
	w := response.NewWriter(ictx, r)
	return r, w, r.Context()
}
