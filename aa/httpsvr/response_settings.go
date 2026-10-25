package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
)

// Body has no methods. It is a plain struct and does not need object-pool management.
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// ErrorHandler may write the response itself.
// The bool tells the writer whether to continue with the normal serializer.
type ErrorHandler func(ctx *Context, contentType string, d Body) (int, error, bool)

// BeforeSerialize can replace the body just before encoding.
type BeforeSerialize func(ctx *Context, contentType string, d Body) Body

// Settings is the per-application HTTP response policy.
// Attach it with Router.UseHooks so two apps do not share handlers.
type Settings struct {
	ServeContentTypes []string
	ErrorHandler      ErrorHandler
	BeforeSerialize   []BeforeSerialize
	Serialize         func(contentType string, d Body) (bytes []byte, newContentType string, err error)
}

var SerializeTag = "json" // `json:"key"`

func NewSettings() *Settings {
	return &Settings{ServeContentTypes: []string{"application/json"}}
}

func (s *Settings) SetServeContentTypes(contentTypes []string) {
	if s == nil {
		return
	}
	if len(contentTypes) == 0 {
		ae.Panic("must register at least one content type")
	}
	s.ServeContentTypes = contentTypes
}

func (s *Settings) SetErrorHandler(f ErrorHandler) {
	if s == nil {
		return
	}
	s.ErrorHandler = f
}

func (s *Settings) AddBeforeSerialize(f BeforeSerialize) {
	if s == nil || f == nil {
		return
	}
	s.BeforeSerialize = append(s.BeforeSerialize, f)
}

func (s *Settings) SetSerialize(f func(contentType string, d Body) (bytes []byte, newContentType string, err error)) {
	if s == nil {
		return
	}
	s.Serialize = f
}

// SettingsFrom reads the settings attached to the request by Router.UseHooks.
func SettingsFrom(ctx *Context) *Settings {
	if ctx == nil || ctx.Hooks() == nil {
		return nil
	}
	s, _ := ctx.Hooks().(*Settings)
	return s
}

func defaultErrorHandler(ctx *Context, contentType string, d Body) (int, error, bool) {
	if d.Code == ae.NotModified {
		ctx.StatusCode(d.Code)
		return 0, nil, false
	}
	if contentType == "text/html" {
		ctx.StatusCode(d.Code)
		return 0, nil, false
	}
	return 0, nil, true
}

func defaultSerialize(contentType string, d Body) ([]byte, string, error) {
	bytes, err := EncodeJson(d)
	if err != nil {
		return nil, "", err
	}
	return bytes, contentType, nil
}
