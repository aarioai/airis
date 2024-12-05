package request

import (
	"context"
	"github.com/aarioai/airis/core/airis"
	"github.com/kataras/iris/v12"
	"net/http"
)

type Request struct {
	ictx            iris.Context
	ctx             context.Context
	r               *http.Request
	contentType     string
	userAgent       string
	partialQueries  map[string]any // 每个请求独立request，因此几乎不存在并发问题，不用sync.Map
	partialHeaders  map[string]any
	partialBodyData map[string]any
	bodyParsed      bool
}

func New(ictx iris.Context) *Request {
	r := ictx.Request()
	req := Request{
		ictx:            ictx,
		r:               r,
		contentType:     "",
		userAgent:       "",
		partialQueries:  nil,
		partialHeaders:  nil,
		partialBodyData: nil,
	}
	paramsLen := ictx.Params().Len()
	if paramsLen > 0 {
		params := make(map[string]any, paramsLen)
		for _, v := range ictx.Params().Store {
			params[v.Key] = v.ValueRaw
		}
	}
	return &req
}
func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return airis.Context(r.ictx)
}
