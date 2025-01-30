package request

import (
	"context"
	"github.com/aarioai/airis/core/airis"
	"github.com/kataras/iris/v12"
	"mime/multipart"
	"net/http"
	"sync"
)

// Request 每个请求独立request，因此几乎不存在并发问题，不用sync.Map
// @extend type T interface{Release()error}
type Request struct {
	ictx             iris.Context
	ctx              context.Context
	r                *http.Request
	contentType      string
	userAgent        string
	bodyParsed       bool
	maxMultipartSize int64
	maxFormSize      int64

	// 程序注入的参数，优先读取（相较于客户端传递的参数）
	injectedHeaders map[string]any
	injectedQueries map[string]any // 程序注入的参数，包括 path params, 以及 SetQueryData 设置的参数
	injectedBodies  map[string]any
	injectedFiles   map[string][]*multipart.FileHeader
}

var (
	// 对象池，减少内存分配
	// sync.Pool 通常不需要手动释放对象，当创建的对象，没有引用时会自动回收
	requestPool = sync.Pool{
		New: func() interface{} {
			return new(Request)
		},
	}
)

func New(ictx iris.Context) *Request {
	req := requestPool.Get().(*Request)
	req.ictx = ictx
	req.ctx = airis.Context(ictx)
	req.r = ictx.Request()
	req.contentType = ""
	req.userAgent = ""
	req.bodyParsed = false
	req.maxMultipartSize = 5 << 20 // 5M equals to 500KB * 9 images
	req.maxFormSize = 1 << 20      // 10 MB is a lot of json/form data.
	req.injectedHeaders = nil
	req.injectedQueries = nil
	req.injectedBodies = nil
	req.injectedFiles = nil

	paramsLen := ictx.Params().Len()
	if paramsLen > 0 {
		params := make(map[string]any, paramsLen)
		for _, v := range ictx.Params().Store {
			params[v.Key] = v.ValueRaw
		}
		req.injectedQueries = params
	}
	return req
}
func (r *Request) SetMaxMultipartSize(size int64) {
	r.maxMultipartSize = size
}
func (r *Request) SetMaxFormSize(size int64) {
	r.maxFormSize = size
}

// InjectHeader 注入header
func (r *Request) InjectHeader(name string, value any) {
	if r.injectedHeaders == nil {
		r.injectedHeaders = make(map[string]any)
	}
	r.injectedHeaders[name] = value
}

// InjectQuery 注入query
func (r *Request) InjectQuery(name string, value any) {
	if r.injectedQueries == nil {
		r.injectedQueries = make(map[string]any)
	}
	r.injectedQueries[name] = value
}

// InjectBody 注入body
func (r *Request) InjectBody(name string, value any) {
	if r.injectedBodies == nil {
		r.injectedBodies = make(map[string]any)
	}
	r.injectedBodies[name] = value
}

func (r *Request) Context() context.Context {
	return r.ctx
}

// Release 释放实例到对象池
// 即使这个对象不是从对象池中获取的，也会放入对象池。不影响使用。
func (r *Request) Release() error {
	requestPool.Put(r)
	return nil
}
