package ae

import (
	"fmt"
	"net/http"
	"strconv"
)

// HTTP 扩展状态码
const (
	// 2xx 成功
	OK        = http.StatusOK        // 请求成功
	NoContent = http.StatusNoContent // 请求成功但无内容返回

	NotModified = http.StatusNotModified

	// 4xx 客户端错误
	BadRequest      = http.StatusBadRequest      // 请求参数错误
	Unauthorized    = http.StatusUnauthorized    // 未授权
	PaymentRequired = http.StatusPaymentRequired // 需要支付
	Forbidden       = http.StatusForbidden       // 禁止访问
	NotFound        = http.StatusNotFound        // 资源不存在
	//MethodNotAllowed      = http.StatusMethodNotAllowed // 方法不允许
	NotAcceptable     = http.StatusNotAcceptable // 不接受的请求
	ProxyAuthRequired = http.StatusProxyAuthRequired
	RequestTimeout    = http.StatusRequestTimeout // 请求超时，注意与429区别
	Conflict          = http.StatusConflict       // 资源冲突
	Gone              = http.StatusGone           // 资源已永久删除
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
	RetryWith       = 491 // 【自定义错误码】需要重试

	// 5xx 服务器错误
	InternalServerError   = http.StatusInternalServerError   // 服务器内部错误
	NotImplemented        = http.StatusNotImplemented        // 未实现
	BadGateway            = http.StatusBadGateway            // 网关错误
	ServiceUnavailable    = http.StatusServiceUnavailable    // 服务不可用
	GatewayTimeout        = http.StatusGatewayTimeout        // 网关超时
	VariantAlsoNegotiates = http.StatusVariantAlsoNegotiates // 服务器内部配置错误

	StatusException = 590 // http 状态码出错，未达到程序阶段。一般由路由层，或nginx等代理层抛出
)

var (
	defaultCodeTexts = map[int]string{
		NoRowsAvailable: "No rows available",
		RetryWith:       "Retry with",
	}

	// 快捷变量，使用时候不需要指定message的，其他的一般都需要指定message供调试或反馈给客户端

	UnauthorizedE          = New(Unauthorized)
	PaymentRequiredE       = New(PaymentRequired)
	ForbiddenE             = New(Forbidden)
	NotFoundE              = New(NotFound)
	NotAcceptableE         = New(NotAcceptable)
	ProxyAuthRequiredE     = New(ProxyAuthRequired)
	TimeoutE               = New(RequestTimeout)
	ConflictE              = New(Conflict)
	GoneE                  = New(Gone)
	RequestEntityTooLargeE = New(RequestEntityTooLarge)
	UnsupportedMediaE      = New(UnsupportedMedia)
	LockedE                = New(Locked)
	TooEarlyE              = New(TooEarly)
	IllegalE               = New(Illegal)
	NoRowsE                = New(NoRowsAvailable) // 自定义状态码

	NotImplementedE     = New(NotImplemented)
	BadGatewayE         = New(BadGateway)
	ServiceUnavailableE = New(ServiceUnavailable)
)

func RetryWithE(redirect string) *Error {
	return New(RetryWith, redirect) // 特殊错误码，msg 用于跳转
}

func BadParamE(param string) *Error {
	return New(BadRequest, "bad param `"+param+"`")
}

func VariantAlsoNegotiatesE(format string, args ...any) *Error {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return New(VariantAlsoNegotiates, msg)
}

func CodeText(code int) string {
	if text, ok := defaultCodeTexts[code]; ok {
		return text
	}
	text := http.StatusText(code)
	if text == "" {
		return strconv.Itoa(code)
	}
	return text
}
