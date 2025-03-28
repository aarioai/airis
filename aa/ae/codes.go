package ae

import (
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/types"
	"net/http"
	"strings"
)

// HTTP 扩展状态码
const (
	// 2xx 成功
	OK             = http.StatusOK             // 请求成功。可用于 GET/POST/DELETE/PUT/PATCH 等，data 返回标准（如json）数据
	Created        = http.StatusCreated        // 创建成功，body 会返回新建的对象（非标准数据类型）。通常用于POST，data 返回特殊字节流
	Accepted       = http.StatusAccepted       // 创建的请求操作成功（不代表创建成功），异步创建过程开始，通常用于POST后异步创建
	NoContent      = http.StatusNoContent      // 请求成功但无内容返回，通常用于DELETE/PUT/PATCH
	ResetContent   = http.StatusResetContent   // 请求成功，服务端做了更新，要求客户端刷新页面，重新拉取接口
	PartialContent = http.StatusPartialContent // 断点续传会使用到

	NotModified = http.StatusNotModified

	// 4xx 客户端错误
	BadRequest      = http.StatusBadRequest      // 请求参数错误
	Unauthorized    = http.StatusUnauthorized    // 未授权
	PaymentRequired = http.StatusPaymentRequired // 需要支付
	Forbidden       = http.StatusForbidden       // 禁止访问
	NotFound        = http.StatusNotFound        // 资源不存在
	//MethodNotAllowed      = http.StatusMethodNotAllowed // 方法不允许
	NotAcceptable     = http.StatusNotAcceptable     // 不接受的请求
	ProxyAuthRequired = http.StatusProxyAuthRequired // 如微信等第三方认证前置
	RequestTimeout    = http.StatusRequestTimeout    // 请求超时，注意与429区别
	Conflict          = http.StatusConflict          // 资源冲突
	Gone              = http.StatusGone              // 资源已永久删除
	//LengthRequired        = http.StatusLengthRequired        // 需要Content-Length
	PreconditionFailed    = http.StatusPreconditionFailed    // 前置条件验证失败，注意与424区别
	RequestEntityTooLarge = http.StatusRequestEntityTooLarge // 请求实体过大
	//RequestURITooLong     = http.StatusRequestURITooLong     // 请求URI过长
	UnsupportedMedia = http.StatusUnsupportedMediaType // 不支持的媒体类型
	Locked           = http.StatusLocked
	FailedDependency = http.StatusFailedDependency           // 之前发生错误，注意与 412 区别
	TooEarly         = http.StatusTooEarly                   // 表示服务器不愿意冒险处理可能被重播的请求。
	TooManyRequests  = http.StatusTooManyRequests            // 请求过多，注意与408区别
	Illegal          = http.StatusUnavailableForLegalReasons // 该请求因政策法律原因不可用。
	// 自定义状态码
	NoRowsAvailable = 490 // 【自定义错误码】无数据记录
	RetryWith       = 491 // 【自定义错误码】需要重试，msg 是 redirect
	ConflictWith    = 492 // 【自定义错误码】数据冲突，msg 是冲突的有效信息

	// 5xx 服务器错误
	InternalServerError   = http.StatusInternalServerError   // 服务器内部错误
	NotImplemented        = http.StatusNotImplemented        // 未实现
	BadGateway            = http.StatusBadGateway            // 网关错误
	ServiceUnavailable    = http.StatusServiceUnavailable    // 服务不可用
	GatewayTimeout        = http.StatusGatewayTimeout        // 网关超时
	VariantAlsoNegotiates = http.StatusVariantAlsoNegotiates // 服务器内部配置错误
	LoopDetected          = http.StatusLoopDetected          // 一般用于内部业务分享，检测到死循环存在
	StatusException       = 590                              // http 状态码出错，未达到程序阶段。一般由路由层，或nginx等代理层抛出
)

var (
	// 自定义状态码
	defaultCodeTexts = map[int]string{
		NoRowsAvailable: "No Rows Available",
		RetryWith:       "Retry With",
		ConflictWith:    "Conflict With",
	}

	// 快捷变量，使用时候不需要指定message的，其他的一般都需要指定message供调试或反馈给客户端
	// ErrorXXX/ErrXXX  都应被视为常量，不应修改

	ErrorNoContent   = New(NoContent).Lock()
	ErrorNotModified = New(NotModified).Lock()

	ErrorUnauthorized    = New(Unauthorized).Lock()
	ErrorPaymentRequired = New(PaymentRequired).Lock()
	ErrorForbidden       = New(Forbidden).Lock()
	ErrorNotFound        = New(NotFound).Lock()

	ErrorTimeout               = New(RequestTimeout).Lock()
	ErrorGone                  = New(Gone).Lock()
	ErrorRequestEntityTooLarge = New(RequestEntityTooLarge).Lock()
	ErrorLocked                = New(Locked).Lock()
	ErrorTooEarly              = New(TooEarly).Lock()
	ErrorIllegal               = New(Illegal).Lock()
	ErrorNoRows                = New(NoRowsAvailable).Lock() // 自定义状态码

	ErrorNotImplemented     = New(NotImplemented).Lock()
	ErrorBadGateway         = New(BadGateway).Lock()
	ErrorServiceUnavailable = New(ServiceUnavailable).Lock()
)

func NewRetryWith(redirect string) *Error {
	return New(RetryWith, redirect) // 特殊错误码，msg 用于跳转
}
func NewConflictWith(format string, args ...any) *Error {
	return New(ConflictWith, afmt.Sprintf(format, args...))
}

func NewBadParam(param string, tips ...string) *Error {
	msg := "bad param `" + param + "`"
	if len(tips) > 0 {
		msg += ": " + strings.Join(tips, " ")
	}
	return New(BadRequest, msg)
}
func NewNotAcceptable(msg ...any) *Error {
	return New(NotAcceptable).AppendMsg(msg...)
}
func NewProxyAuthRequired(msg ...any) *Error {
	return New(ProxyAuthRequired).AppendMsg(msg...)
}

func NewConflict(name string) *Error {
	return New(Conflict).AppendMsg(name)
}

func NewPreconditionFailed(msg ...any) *Error {
	return New(PreconditionFailed).AppendMsg(msg...)
}

func NewUnsupportedMedia(wants ...string) *Error {
	e := New(UnsupportedMedia)
	if len(wants) > 0 {
		e.AppendMsg("want " + strings.Join(wants, " or "))
	}
	return e
}
func NewFailedDependency(msg ...any) *Error {
	return New(FailedDependency).AppendMsg(msg...)
}
func NewVariantAlsoNegotiates(format string, args ...any) *Error {
	return New(VariantAlsoNegotiates, afmt.Sprintf(format, args...))
}
func NewLoopDetected(msg ...any) *Error {
	return New(LoopDetected).AppendMsg(msg...)
}

// Text 获取错误码对应的文本描述
func Text(code int) string {
	if text, ok := defaultCodeTexts[code]; ok {
		return text
	}
	text := http.StatusText(code)
	if text == "" {
		return types.Itoa(code)
	}
	return text
}
