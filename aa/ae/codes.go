package ae

import (
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/types"
	"net/http"
	"strings"
)

const (
	Continue           = 100
	SwitchingProtocols = 101
	Processing         = 102
	EarlyHints         = 103

	OK                   = 200
	Created              = 201
	Accepted             = 202
	NonAuthoritativeInfo = 203
	NoContent            = 204
	ResetContent         = 205 // @diff with 303, POST/PUT/PATCH success, the server add some special changes, need client to refresh this page
	PartialContent       = 206 // for breakpoint resume
	MultiStatus          = 207 // WebDAV
	AlreadyReported      = 208 // WebDAV
	IMUsed               = 226

	MultipleChoices   = 300
	MovePermanently   = 301 // 301 redirect, and change method to GET
	Found             = 302 // same as MovedTemporarily, 302 redirect, and change method to GET
	SeeOther          = 303 // @diff with 205
	NotModified       = 304 // @notice diff with 200, success but using cache
	UseProxy          = 305
	_                 = 306
	TemporaryRedirect = 307 // @notice diff with 302, redirect with the same method (e.g. POST/PUT)
	PermanentRedirect = 308 // @notice diff with 301, redirect with the same method (e.g. POST/PUT)
	FailedAndSeeOther = 391 // [NEW] failed and let user redirect to other page

	BadRequest            = 400
	Unauthorized          = 401
	PaymentRequired       = 402
	Forbidden             = 403
	NotFound              = 404
	MethodNotAllowed      = 405
	NotAcceptable         = 406
	ProxyAuthRequired     = 407 // e.g. wechat auth required
	RequestTimeout        = 408 // @notice diff with 429
	Conflict              = 409
	Gone                  = 410 // deleted permanently, most servers use 404 instead 410
	LengthRequired        = 411 // require content length, buffer length or other length
	PreconditionFailed    = 412 // @notice diff with 424
	RequestEntityTooLarge = 413 // post data oversize
	// RequestURIInvalid is alias to RequestURITooLong
	RequestURIInvalid            = 414 // invalid URI
	UnsupportedMedia             = 415 // e.g. required json, but offers xml
	RequestedRangeNotSatisfiable = 416 // e.g. user requires data from 1st page to 100th page, but only 3 pages available
	ExpectationFailed            = 417
	PageExpired                  = 419 // [NEW] Laravel Framework accepted this code. CSRF and other token missing or expired
	EnhanceYourCalm              = 420 // [NEW] Twitter and trends API when the client is being rate limited
	Locked                       = 423
	FailedDependency             = 424 // @notice diff with 412
	TooEarly                     = 425 // reject handle Early Data to defend replay attack
	UpgradeRequired              = 426
	_                            = 427
	PreconditionRequired         = 428
	TooManyRequests              = 429 // @notice diff with 408
	_                            = 430
	RequestHeaderFieldsTooLarge  = 431
	LoginTimeout                 = 440
	NoResponse                   = 444
	UnavailableForLegalReasons   = 451
	NoRowsAvailable              = 494 // [NEW] other servers may use {code:200, text:[]} for empty list
	ConflictWith                 = 499 // [NEW] msg contains conflict data

	InternalServerError           = 500
	NotImplemented                = 501
	BadGateway                    = 502
	ServiceUnavailable            = 503
	GatewayTimeout                = 504
	_                             = 505 // includes HTTPVersionNotSupported
	VariantAlsoNegotiates         = 506
	InsufficientStorage           = 507
	LoopDetected                  = 508
	_                             = 509
	NotExtended                   = 510
	NetworkAuthenticationRequired = 511
	Exception                     = 590 // [NEW] exceptions
)

var (
	newCodeTexts = map[int]string{
		FailedAndSeeOther: "Failed And See Other",
		PageExpired:       "Page Expired",
		EnhanceYourCalm:   "Enhance Your Calm",
		NoRowsAvailable:   "No Rows Available",
		ConflictWith:      "Conflict With",
		Exception:         "Exception",
	}
)

// constant errors
var (
	ErrorContinue        = New(Continue).Lock()
	ErrorSwitchProtocols = New(SwitchingProtocols).Lock()
	ErrorProcessing      = New(Processing).Lock()
	ErrorEarlyHints      = New(EarlyHints).Lock()

	ErrorOK               = New(OK).Lock()
	ErrorCreated          = New(Created).Lock()
	ErrorAccepted         = New(Accepted).Lock()
	ErrorNonAuthoritative = New(NonAuthoritativeInfo).Lock()
	ErrorNoContent        = New(NoContent).Lock()
	ErrorResetContent     = New(ResetContent).Lock()
	ErrorPartialContent   = New(PartialContent).Lock()
	ErrorMultiStatus      = New(MultiStatus).Lock()
	ErrorAlreadyReported  = New(AlreadyReported).Lock()
	ErrorIMUsed           = New(IMUsed).Lock()

	ErrorMultipleChoices   = New(MultipleChoices).Lock()
	ErrorMovePermanently   = New(MovePermanently).Lock()
	ErrorFound             = New(Found).Lock()
	ErrorSeeOther          = New(SeeOther).Lock()
	ErrorNotModified       = New(NotModified).Lock()
	ErrorUseProxy          = New(UseProxy).Lock()
	ErrorTemporaryRedirect = New(TemporaryRedirect).Lock()
	ErrorPermanentRedirect = New(PermanentRedirect).Lock()

	ErrorBadRequest                   = New(BadRequest).Lock()
	ErrorUnauthorized                 = New(Unauthorized).Lock()
	ErrorPaymentRequired              = New(PaymentRequired).Lock()
	ErrorForbidden                    = New(Forbidden).Lock()
	ErrorNotFound                     = New(NotFound).Lock()
	ErrorMethodNotAllowed             = New(MethodNotAllowed).Lock()
	ErrorNotAcceptable                = New(NotAcceptable).Lock()
	ErrorProxyAuthRequired            = New(ProxyAuthRequired).Lock()
	ErrorRequestTimeout               = New(RequestTimeout).Lock()
	ErrorConflict                     = New(Conflict).Lock()
	ErrorGone                         = New(Gone).Lock()
	ErrorLengthRequired               = New(LengthRequired).Lock()
	ErrorPreconditionFailed           = New(PreconditionFailed).Lock()
	ErrorRequestEntityTooLarge        = New(RequestEntityTooLarge).Lock()
	ErrorRequestURIInvalid            = New(RequestURIInvalid).Lock()
	ErrorUnsupportedMedia             = New(UnsupportedMedia).Lock()
	ErrorRequestedRangeNotSatisfiable = New(RequestedRangeNotSatisfiable).Lock()
	ErrorExpectationFailed            = New(ExpectationFailed).Lock()
	ErrorPageExpired                  = New(PageExpired).Lock()
	ErrorEnhanceYourCalm              = New(EnhanceYourCalm).Lock()
	ErrorLocked                       = New(Locked).Lock()
	ErrorFailedDependency             = New(FailedDependency).Lock()
	ErrorTooEarly                     = New(TooEarly).Lock()
	ErrorUpgradeRequired              = New(UpgradeRequired).Lock()
	ErrorPreconditionRequired         = New(PreconditionRequired).Lock()
	ErrorTooManyRequests              = New(TooManyRequests).Lock()
	ErrorRequestHeaderFieldsTooLarge  = New(RequestHeaderFieldsTooLarge).Lock()
	ErrorLoginTimeout                 = New(LoginTimeout).Lock()
	ErrorNoResponse                   = New(NoResponse).Lock()
	ErrorUnavailableForLegalReasons   = New(UnavailableForLegalReasons).Lock()
	ErrorNoRowsAvailable              = New(NoRowsAvailable).Lock()

	ErrorInternalServerError           = New(InternalServerError).Lock()
	ErrorNotImplemented                = New(NotImplemented).Lock()
	ErrorBadGateway                    = New(BadGateway).Lock()
	ErrorServiceUnavailable            = New(ServiceUnavailable).Lock()
	ErrorGatewayTimeout                = New(GatewayTimeout).Lock()
	ErrorVariantAlsoNegotiates         = New(VariantAlsoNegotiates).Lock()
	ErrorInsufficientStorage           = New(InsufficientStorage).Lock()
	ErrorLoopDetected                  = New(LoopDetected).Lock()
	ErrorNotExtended                   = New(NotExtended).Lock()
	ErrorNetworkAuthenticationRequired = New(NetworkAuthenticationRequired).Lock()
	ErrorException                     = New(Exception).Lock()
)

func NewFailedAndSeeOther(redirect string) *Error {
	return New(FailedAndSeeOther, redirect) // 特殊错误码，msg 用于跳转
}
func NewConflictWith(format string, args ...any) *Error {
	return New(ConflictWith, afmt.Sprintf(format, args...))
}

func NewBadParam(param string, tips ...string) *Error {
	msg := fmt.Sprintf(BadParameterFormat, param)
	if len(tips) > 0 {
		msg += Separator + " " + strings.Join(tips, Separator)
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

// Text returns a text for the code. It returns the empty string if the code is unknown.
func Text(code int) string {
	if text, ok := newCodeTexts[code]; ok {
		return text
	}
	text := http.StatusText(code)
	if text == "" {
		return types.Itoa(code)
	}
	return text
}
