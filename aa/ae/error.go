package ae

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/airis/pkg/utils"
	"github.com/kataras/iris/v12"
	"log"
	"strings"
)

// Error 定义标准错误结构
type Error struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Caller    string `json:"caller"`
	Detail    string `json:"details"`
	TraceInfo string `json:"trace_info"`
	locked    bool
}

// New 使用错误码和消息创建 Error
func New(code int, msgs ...any) *Error {
	e := &Error{
		Code:   code,
		Caller: utils.Caller(1),
	}

	if msg := afmt.SprintfArgs(msgs); msg != "" {
		e.Msg = msg
	} else {
		e.Msg = CodeText(code)
	}

	return e
}

// NewE 使用消息创建 Error
func NewE(format string, args ...any) *Error {
	return &Error{
		Code:   InternalServerError,
		Msg:    afmt.Sprintf(format, args...),
		Caller: utils.Caller(1),
	}
}

// NewError 从标准 error 创建 Error
func NewError(err error, details ...any) *Error {
	if err == nil {
		return nil
	}
	return NewE(err.Error()).WithCaller(1).WithDetail(details...)
}

// Lock 锁定后不可再修改或复用，一般作为常量使用
func (e *Error) Lock() *Error {
	e.locked = true
	return e
}
func (e *Error) Unlock() *Error {
	e.locked = false
	return e
}
func (e *Error) WithMsg(format string, args ...any) *Error {
	if e.locked {
		Panic("unable change locked error")
		return e
	}
	e.Msg = afmt.Sprintf(format, args...)
	return e
}

// AppendMsg 尝试添加消息
func (e *Error) AppendMsg(msgs ...any) *Error {
	msg := afmt.SprintfArgs(msgs)
	if e.locked {
		log.Printf("[error] failed to append message %s to a locked error\n", msg)
		return e
	}

	if msg != "" {
		e.Msg += " - " + msg
	}
	return e
}

// WithCaller 添加调用者信息
func (e *Error) WithCaller(skip int) *Error {
	if e.locked {
		log.Printf("[error] failed to change caller(%d) to a locked error\n", skip)
		return e
	}
	e.Caller = utils.Caller(skip)
	return e
}

func (e *Error) WithDetail(detail ...any) *Error {
	s := afmt.SprintfArgs(detail)
	if e.locked {
		log.Printf("[error] failed to change locked Error with detail: %s\n", s)
		return e
	}
	e.Detail = s
	return e
}
func (e *Error) WithTraceInfo(ctx iris.Context) *Error {
	if e.locked {
		log.Println("[error] failed to change trace info to a locked error")
		return e
	}
	e.TraceInfo = acontext.TraceInfo(ctx)
	return e
}

// String 输出错误信息，最好不要使用 Error，避免跟 error 一致，导致人写的时候发生失误
// $caller {$trace_info} code:$code $msg\n$detail
func (e *Error) String() string {
	capacity := len(e.TraceInfo) + len(e.Msg) + len(e.Detail) + 10
	var s strings.Builder
	s.Grow(capacity)
	if e.TraceInfo != "" {
		s.WriteString(e.TraceInfo)
		s.WriteByte(' ')
	}

	s.WriteString("code:")
	s.WriteString(types.FormatInt(e.Code))
	s.WriteByte(' ')
	s.WriteString(e.Msg)
	if e.Detail != "" {
		s.WriteByte('\n')
		s.WriteString(e.Detail)
	}
	return s.String()
}
func (e *Error) Text() string {
	return e.Caller + " " + e.String()
}

func (e *Error) Trace(ctx iris.Context) string {
	return e.WithTraceInfo(ctx).String()
}

func (e *Error) IsNotFound() bool {
	return e.Code == NotFound || e.Code == NoRowsAvailable || e.Code == Gone
}

func (e *Error) IsServerError() bool {
	return e.Code >= 500 && e.Code <= 599
}

func (e *Error) IsRetryWith() bool {
	return e.Code == RetryWith && e.Msg != ""
}

func (e *Error) ExceptNotFound() *Error {
	if e == nil || e.IsNotFound() {
		return nil
	}
	return e
}
