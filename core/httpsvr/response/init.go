package response

import (
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/httpsvr/request"
	"github.com/aarioai/airis/core/logger"
	"github.com/kataras/iris/v12"
)

// Body 这个没有方法，单纯结构体，不需要对象池管理
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

var (
	// GlobalHideServerError 全局隐藏服务器错误(code>=500)
	// 优先级：writer.HideServerError > GlobalHideServerError
	globalErrorHandler      func(ictx iris.Context, request *request.Request, contentType string, d Body) (int, error, bool)
	globalBeforeSerialize   []func(ictx iris.Context, contentType string, d Body) Body
	globalSerialize         func(contentType string, d Body) (bytes []byte, newContentType string, err error)
	log                     = logger.NewDefaultLog()
	SerializeTag            = "json" // `json:"key"`
	globalServeContentTypes = []string{"application/json"}
)

// 避免循环调用，避免传递 *Writer
func defaultErrorHandler(ictx iris.Context, request *request.Request, contentType string, d Body) (int, error, bool) {
	if d.Code == ae.NotModified {
		ictx.StatusCode(d.Code)
		return 0, nil, false
	}
	if contentType == "text/html" {
		// 返回状态码，交给route层处理
		//w.ictx.Values().Set(ErrCodeKey, d.Code)
		//w.ictx.Values().Set(ErrMsgKey, d.Msg)
		ictx.StatusCode(d.Code)
		return 0, nil, false
	}
	// 丢给 next 执行
	return 0, nil, true
}

func defaultSerialize(contentType string, d Body) ([]byte, string, error) {
	bytes, err := EncodeJson(d)
	if err != nil {
		return nil, "", err
	}
	return bytes, contentType, nil
}

func RegisterGlobalServeContentTypes(contentTypes []string) {
	if len(contentTypes) == 0 {
		panic("must register at least one content type")
	}
	globalServeContentTypes = contentTypes
}
func RegisterGlobalErrorHandler(f func(ictx iris.Context, request *request.Request, contentType string, d Body) (int, error, bool)) {
	globalErrorHandler = f
}
func RegisterGlobalBeforeSerialize(f func(ictx iris.Context, contentType string, d Body) Body) {
	if globalBeforeSerialize == nil {
		globalBeforeSerialize = make([]func(ictx iris.Context, contentType string, d Body) Body, 0)
	}
	globalBeforeSerialize = append(globalBeforeSerialize, f)
}
func RegisterGlobalSerialize(f func(contentType string, d Body) (bytes []byte, newContentType string, err error)) {
	globalSerialize = f
}

func RegisterLogger(logger logger.LogInterface) {
	log = logger
}
