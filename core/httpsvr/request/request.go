package request

import (
	"context"
	"github.com/aarioai/airis/core/airis"
	"github.com/kataras/iris/v12"
	"net/http"
)

// 每个请求独立request，因此几乎不存在并发问题，不用sync.Map
type Request struct {
	ictx            iris.Context
	ctx             context.Context
	r               *http.Request
	contentType     string
	userAgent       string
	bodyParsed      bool
	
	// 程序注入的参数，优先读取（相较于客户端传递的参数）
	injectedHeaders map[string]any
	injectedQueries map[string]any // 程序注入的参数，包括 path params, 以及 SetQueryData 设置的参数
	injectedBodies  map[string]any

	
}

func New(ictx iris.Context) *Request {
	r := ictx.Request()
	req := Request{
		ictx:            ictx,
		r:               r,
		contentType:     "",
		userAgent:       "",
		injectedQueries: nil,
		injectedHeaders: nil,
		injectedBodies:  nil,
	}
	paramsLen := ictx.Params().Len()
	if paramsLen > 0 {
		params := make(map[string]any, paramsLen)
		for _, v := range ictx.Params().Store {
			params[v.Key] = v.ValueRaw
		}
		req.injectedQueries = params
	}
	return &req
}

// 注入header
func (r *Request) InjectHeader(name string, value any) {
	if r.injectedHeaders == nil {
		r.injectedHeaders = make(map[string]any)
	}
	r.injectedHeaders[name] = value
}

// 注入query
func (r *Request) InjectQuery(name string, value any) {
	if r.injectedQueries == nil {
		r.injectedQueries = make(map[string]any)
	}
	r.injectedQueries[name] = value
}

// 注入body
func (r *Request) InjectBody(name string, value any) {
	if r.injectedBodies == nil {
		r.injectedBodies = make(map[string]any)
	}
	r.injectedBodies[name] = value
}

func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return airis.Context(r.ictx)
}
