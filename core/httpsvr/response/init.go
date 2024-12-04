package response

import (
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/logger"
	"github.com/kataras/iris/v12"
)

var (
	// GlobalHideServerError 全局隐藏服务器错误(code>=500)
	// 优先级：writer.HideServerError > GlobalHideServerError
	globalErrorHandler      func(w *Writer, contentType string, d Response) (int, error)
	globalBeforeSerialize   []func(ictx iris.Context, contentType string, d Response) Response
	globalSerialize         func(contentType string, d Response) (bytes []byte, newContentType string, err error)
	log                     = logger.NewDefaultLog()
	SerializeTag            = "json" // `json:"key"`
	globalServeContentTypes = []string{"application/json"}
)

func defaultErrorHandler(w *Writer, contentType string, d Response) (int, error) {
	if d.Code == ae.CodeNotModified {
		w.StatusCode(d.Code)
		return 0, nil
	}
	if contentType == "text/html" {
		// 返回状态码，交给route层处理
		//w.ictx.Values().Set(ErrCodeKey, d.Code)
		//w.ictx.Values().Set(ErrMsgKey, d.Msg)
		w.StatusCode(d.Code)
		return 0, nil
	}
	return w.WriteMsg(d.Code, d.Msg)
}

func defaultSerialize(contentType string, d Response) ([]byte, string, error) {
	bytes, err := EncodeJson(d)
	if err != nil {
		return nil, "", err
	}
	return bytes, contentType, nil
}
func init() {
}

func RegisterGlobalServeContentTypes(contentTypes []string) {
	if len(contentTypes) == 0 {
		panic("must register at least one content type")
	}
	globalServeContentTypes = contentTypes
}
func RegisterGlobalErrorHandler(f func(ictx iris.Context, contentType string, d Response) (int, error)) {
	globalErrorHandler = f
}
func RegisterGlobalBeforeSerialize(f func(ictx iris.Context, contentType string, d Response) Response) {
	if globalBeforeSerialize == nil {
		globalBeforeSerialize = make([]func(ictx iris.Context, contentType string, d Response) Response, 0)
	}
	globalBeforeSerialize = append(globalBeforeSerialize, f)
}
func RegisterGlobalSerialize(f func(contentType string, d Response) (bytes []byte, newContentType string, err error)) {
	globalSerialize = f
}

func RegisterLogger(logger logger.LogInterface) {
	log = logger
}
