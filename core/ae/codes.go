package ae

import "net/http"

// HTTP 扩展状态码
const (
	// 2xx 成功
	CodeOK        = http.StatusOK        // 请求成功
	CodeNoContent = http.StatusNoContent // 请求成功但无内容返回

	CodeNotModified = http.StatusNotModified

	// 4xx 客户端错误
	CodeBadRequest            = http.StatusBadRequest       // 请求参数错误
	CodeUnauthorized          = http.StatusUnauthorized     // 未授权
	CodePaymentRequired       = http.StatusPaymentRequired  // 需要支付
	CodeForbidden             = http.StatusForbidden        // 禁止访问
	CodeNotFound              = http.StatusNotFound         // 资源不存在
	CodeMethodNotAllowed      = http.StatusMethodNotAllowed // 方法不允许
	CodeNotAcceptable         = http.StatusNotAcceptable    // 不接受的请求
	CodeProxyAuthRequired     = http.StatusProxyAuthRequired
	CodeRequestTimeout        = http.StatusRequestTimeout        // 请求超时
	CodeConflict              = http.StatusConflict              // 资源冲突
	CodeGone                  = http.StatusGone                  // 资源已永久删除
	CodeLengthRequired        = http.StatusLengthRequired        // 需要Content-Length
	CodePreconditionFailed    = http.StatusPreconditionFailed    // 前置条件验证失败
	CodeRequestEntityTooLarge = http.StatusRequestEntityTooLarge // 请求实体过大
	CodeRequestURITooLong     = http.StatusRequestURITooLong     // 请求URI过长
	CodeUnsupportedMedia      = http.StatusUnsupportedMediaType  // 不支持的媒体类型
	CodeTooManyRequests       = http.StatusTooManyRequests       // 请求过多

	// 自定义状态码
	CodeNoRows    = 444 // 【自定义错误码】无数据记录
	CodeRetryWith = 449 // 需要重试

	// 5xx 服务器错误
	CodeInternalError      = http.StatusInternalServerError // 服务器内部错误
	CodeNotImplemented     = http.StatusNotImplemented      // 未实现
	CodeBadGateway         = http.StatusBadGateway          // 网关错误
	CodeServiceUnavailable = http.StatusServiceUnavailable  // 服务不可用
	CodeGatewayTimeout     = http.StatusGatewayTimeout      // 网关超时
	CodeBandwidthLimit     = 509                            // 带宽限制
	CodeStatusException    = 555                            // http 状态码出错，未达到程序阶段
)

// 预定义错误
var (
	// 2xx
	OK        = New(CodeOK, http.StatusText(CodeOK))
	NoContent = New(CodeNoContent, http.StatusText(CodeNoContent))

	// 4xx
	BadRequest             = New(CodeBadRequest, http.StatusText(CodeBadRequest))
	Unauthorized           = New(CodeUnauthorized, http.StatusText(CodeUnauthorized))
	PaymentRequired        = New(CodePaymentRequired, http.StatusText(CodePaymentRequired))
	Forbidden              = New(CodeForbidden, http.StatusText(CodeForbidden))
	NotFound               = New(CodeNotFound, http.StatusText(CodeNotFound))
	MethodNotAllowed       = New(CodeMethodNotAllowed, http.StatusText(CodeMethodNotAllowed))
	NotAcceptable          = New(CodeNotAcceptable, http.StatusText(CodeNotAcceptable))
	ProxyAuthRequiredError = New(CodeProxyAuthRequired, http.StatusText(CodeProxyAuthRequired))
	Timeout                = New(CodeRequestTimeout, http.StatusText(CodeRequestTimeout))
	Conflict               = New(CodeConflict, http.StatusText(CodeConflict))
	Gone                   = New(CodeGone, http.StatusText(CodeGone))
	LengthRequired         = New(CodeLengthRequired, http.StatusText(CodeLengthRequired))
	PreconditionFailed     = New(CodePreconditionFailed, http.StatusText(CodePreconditionFailed))
	RequestEntityTooLarge  = New(CodeRequestEntityTooLarge, http.StatusText(CodeRequestEntityTooLarge))
	RequestURITooLong      = New(CodeRequestURITooLong, http.StatusText(CodeRequestURITooLong))
	UnsupportedMedia       = New(CodeUnsupportedMedia, http.StatusText(CodeUnsupportedMedia))
	TooManyRequests        = New(CodeTooManyRequests, http.StatusText(CodeTooManyRequests))
	NoRows                 = New(CodeNoRows, "No Rows")       // 自定义状态码
	RetryWithError         = New(CodeRetryWith, "Retry With") // 自定义状态码

	// 5xx
	InternalServerError = New(CodeInternalError, http.StatusText(CodeInternalError))
	NotImplemented      = New(CodeNotImplemented, http.StatusText(CodeNotImplemented))
	BadGateway          = New(CodeBadGateway, http.StatusText(CodeBadGateway))
	ServiceUnavailable  = New(CodeServiceUnavailable, http.StatusText(CodeServiceUnavailable))
	GatewayTimeout      = New(CodeGatewayTimeout, http.StatusText(CodeGatewayTimeout))
	BandwidthLimit      = New(CodeBandwidthLimit, "Bandwidth Limit Exceeded") // 自定义状态码
	StatusException     = New(CodeStatusException, "Server Status Exception") // 自定义状态码
)
