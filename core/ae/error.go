package ae

import (
	"fmt"
	"github.com/aarioai/airis/core/airis"
	"github.com/aarioai/airis/pkg/arrmap"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

// Error 定义标准错误结构
type Error struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Caller    string `json:"caller"`
	Detail    string `json:"details"`
	TraceInfo string `json:"trace_info"`
}

// New 使用错误码和消息创建 Error
func New(code int, msgs ...any) *Error {
	e := &Error{
		Code:   code,
		Caller: Caller(2),
	}

	if msg := arrmap.SprintfArgs(msgs); msg != "" {
		e.Msg = msg
	} else {
		e.Msg = CodeText(code)
	}

	return e
}

// NewMsg 使用消息创建 Error
func NewMsg(format string, args ...any) *Error {
	return &Error{
		Code:   InternalServerError,
		Msg:    fmt.Sprintf(format, args...),
		Caller: Caller(2),
	}
}

// NewError 从标准 error 创建 Error
func NewError(err error, details ...string) *Error {
	if err == nil {
		return nil
	}
	return NewMsg(err.Error()).WithCaller(2).withDetail(details...)
}

// WithCaller 添加调用者信息
func (e *Error) WithCaller(skip int) *Error {
	if e == nil {
		return nil
	}
	e.Caller = Caller(skip + 1)
	return e
}
func (e *Error) withDetail(details ...string) *Error {
	if e == nil {
		return nil
	}
	if len(details) == 1 {
		return e.WithDetail(details[0])
	} else if len(details) > 1 {
		args := make([]any, len(details)-1)
		for i := 1; i < len(details); i++ {
			args[i-1] = details[i]
		}
		return e.WithDetail(details[0], args...)
	}
	return e
}

// WithDetail 添加详细信息
func (e *Error) WithDetail(format string, args ...any) *Error {
	if e == nil {
		return nil
	}
	e.Detail = fmt.Sprintf(format, args...)
	return e
}
func (e *Error) WithTraceInfo(ctx iris.Context) *Error {
	if e == nil {
		return nil
	}
	e.TraceInfo = airis.TraceInfo(ctx)
	return e
}

// Text 输出错误信息，最好不要使用 Error，避免跟 error 一致，导致人写的时候发生失误
// $caller {$trace_info} code:$code $msg\n$detail
func (e *Error) Text() string {
	if e == nil {
		return "<nil>"
	}
	var s strings.Builder
	s.Grow(32)
	if e.Caller != "" {
		s.WriteString(e.Caller)
		s.WriteByte(' ')
	}
	if e.TraceInfo != "" {
		s.WriteByte('{')
		s.WriteString(e.TraceInfo)
		s.WriteString("} ")
	}
	s.WriteString("code:")
	s.WriteString(strconv.Itoa(e.Code))
	s.WriteByte(' ')
	s.WriteString(e.Msg)
	if e.Detail != "" {
		s.WriteByte('\n')
		s.WriteString(e.Detail)
	}
	return s.String()
}
func (e *Error) Trace(ctx iris.Context) string {
	return e.WithTraceInfo(ctx).Text()
}

// 状态检查方法
func (e *Error) IsNotFound() bool {
	return e != nil && (e.Code == NotFound || e.Code == NoRowsAvailable || e.Code == Gone)
}

func (e *Error) IsServerError() bool {
	return e != nil && e.Code >= 500 && e.Code <= 599
}

func (e *Error) IsRetryWith() bool {
	return e != nil && e.Code == RetryWith && e.Msg != ""
}
