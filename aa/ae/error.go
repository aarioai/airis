package ae

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/airis/pkg/utils"
	"github.com/kataras/iris/v12"
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
func New(code int, message ...any) *Error {
	e := &Error{
		Code:   code,
		Caller: utils.Caller(1),
	}

	if msg := afmt.SprintfArgs(message); msg != "" {
		e.Msg = msg
	} else {
		e.Msg = Text(code)
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
func (e *Error) Clone() *Error {
	return &Error{
		Code:      e.Code,
		Msg:       e.Msg,
		Caller:    e.Caller,
		Detail:    e.Detail,
		TraceInfo: e.TraceInfo,
	}
}

func (e *Error) handleLock() *Error {
	if e.locked {
		return e.Clone()
	}
	return e
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
	newE := e.handleLock()
	newE.Msg = afmt.Sprintf(format, args...)
	return newE
}

// AppendMsg 尝试添加消息
func (e *Error) AppendMsg(message ...any) *Error {
	msg := afmt.SprintfArgs(message)
	if msg == "" {
		return e
	}

	newE := e.handleLock()
	newE.Msg += " - " + msg
	return newE
}

// WithCaller 添加调用者信息
func (e *Error) WithCaller(skip int) *Error {
	newE := e.handleLock()
	newE.Caller = utils.Caller(skip)
	return newE
}

func (e *Error) WithDetail(detail ...any) *Error {
	newE := e.handleLock()
	s := afmt.SprintfArgs(detail)
	newE.Detail = s
	return newE
}

func (e *Error) WithTraceInfo(ctx iris.Context) *Error {
	newE := e.handleLock()
	newE.TraceInfo = acontext.TraceInfo(ctx)
	return newE
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

func (e *Error) IsNotMatch() bool {
	return e.Code == NotFound || e.Code == Gone || e.Code == NoRowsAvailable
}

func (e *Error) IsServerError() bool {
	return e.Code >= 500 && e.Code <= 599
}

func (e *Error) IsFailedAndSeeOther() bool {
	return e.Code == FailedAndSeeOther && e.Msg != ""
}

func (e *Error) ExceptNotFound() *Error {
	if e == nil || e.IsNotMatch() {
		return nil
	}
	return e
}
