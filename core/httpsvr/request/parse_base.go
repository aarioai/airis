package request

import (
	"encoding/json"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/atype"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	maxMultiSize = 32 << 20 // 32M
	maxInt64     = int64(1<<63 - 1)
	maxFormSize  = 10 << 20 // 10 MB is a lot of json/form data.
)

type userDataInterface interface {
	Get(key string) string
}
type RawValue struct {
	*atype.Atype
	name string
}

func newRawValue(name string, value any) *RawValue {
	return &RawValue{
		Atype: atype.New(value),
		name:  name,
	}
}

// Filter 验证并过滤值
// @param pattern  e.g. `[[:word:]]+` `\w+`
//	Filter(pattern string, required bool)
//	Filter(required bool)
//	Filter(pattern string)
//	Filter(default atype.Atype)
func (v *RawValue) Filter(patterns ...any) *ae.Error {
	required := true
	pattern := ""

	for _, pat := range patterns {
		switch p := pat.(type) {
		case string:
			pattern = p
		case bool:
			required = p
		case *atype.Atype:
			if v.String() == "" {
				v.Reload(p.Raw())
			}
		}
	}

	if v.String() == "" {
		if required {
			return ae.BadParam(v.name)
		}
		return nil
	}
	if pattern != "" {
		re, _ := regexp.Compile(pattern)
		m := re.FindStringSubmatch(v.String())
		if m == nil || len(m) < 1 {
			return ae.BadParam(v.name)
		}
	}
	return nil
}

func (r *Request) contentMediaType() (mediatype string, params map[string]string, err error) {
	// @see http.parsePostForm()
	ct := r.r.Header.Get("Content-Type")
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if ct == "" {
		ct = CtOctetStream.String()
	}
	return mime.ParseMediaType(ct)
}

// ContentType
// Request是每个请求独立内存，基本不存在并发情况，没有加锁的必要性
func (r *Request) ContentType() string {
	if r.contentType != "" {
		return r.contentType
	}
	r.contentType, _, _ = r.contentMediaType()
	return r.contentType
}

func (r *Request) Method() string {
	return r.ictx.Method()
}

// Origin e.g. https://luexu.com  80口省略端口，非80口，带上 :$port
// ${scheme}//${host}[:$port]
func (r *Request) Origin() string {
	scheme := r.r.URL.Scheme
	if scheme == "" {
		if r.r.TLS != nil {
			scheme = "https:"
		} else {
			scheme = "http:"
		}
	}
	// 由于是通过接口做跳转，所以不可行，只会对这个接口结果跳转， 不会对页面跳转！！还需要客户端自行处理
	host := r.r.Host
	h := scheme + "//" + host
	return h
}

func (r *Request) query(programData map[string]any, userData userDataInterface, name string, patterns ...any) (*RawValue, *ae.Error) {
	raw := newRawValue(name, "")
	// 读取程序设置的header
	if programData != nil && programData[name] != nil {
		raw.Reload(programData[name])
		return raw, raw.Filter(patterns...)
	}
	// 从ictx.Header读取用户端传递的header
	if userData != nil {
		v := userData.Get(name)
		raw.Reload(v)
		return raw, raw.Filter(patterns...)
	}
	return raw, raw.Filter(patterns...)
}

// Headers 获取所有headers
// 这个读取少，直接每次独立解析即可
func (r *Request) queries(programData map[string]any, userData map[string][]string) map[string]any {
	data := make(map[string]any)
	// 优先级最低，读取用户的header
	if r.r != nil {
		rhs := r.r.Header
		for k, vs := range rhs {
			if len(vs) != 0 && vs[0] != "" {
				data[k] = vs[0]
			}
		}
	}
	// 优先级高，读取程序设置的header
	if r.partialHeaders != nil {
		for k, v := range r.partialHeaders {
			data[k] = v
		}
	}
	if len(data) > 0 {
		return data
	}
	return nil
}

// Query 从URL参数获取值
func (r *Request) Query(name string, patterns ...any) (*RawValue, *ae.Error) {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.query(r.partialQueries, userData, name, patterns...)
}

func (r *Request) Queries() map[string]any {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.queries(r.partialHeaders, userData)
}
func (r *Request) setPartialBodyData(data map[string][]string) {
	if len(data) == 0 {
		return
	}
	if r.partialBodyData == nil {
		r.partialBodyData = make(map[string]any, len(data))
	}
	for k, v := range data {
		if len(v) > 0 {
			r.partialBodyData[k] = v[0]
		}
	}
}

// parseBodyStream 解析请求体
// @see http.parsePostForm()
func (r *Request) parseBodyStream() *ae.Error {
	defer func() { r.bodyParsed = true }()
	// body 可以不传数据
	if r.r == nil || r.r.Body == nil {
		return nil
	}

	ct, params, err := r.contentMediaType()
	if err != nil {
		return ae.New(ae.CodeBadRequest, "invalid content media type")
	}
	switch ct {
	case CtJSON.String(), CtOctetStream.String(), CtForm.String():
		return r.parseSimpleBody()
	case CtFormData.String():
		return r.parseMultipartBody(params["boundary"])
	}
	return nil
}

// parseSimpleBody 解析简单请求体
func (r *Request) parseSimpleBody() *ae.Error {
	var reader io.Reader = r.r.Body
	maxSize := maxInt64
	if _, ok := reader.(*maxBytesReader); !ok {
		maxSize = maxFormSize
		reader = io.LimitReader(r.r.Body, maxSize+1)
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return ae.NewError(err)
	}
	if int64(len(b)) > maxSize {
		return ae.New(ae.CodeRequestEntityTooLarge, "body too large")
	}
	switch r.ContentType() {
	case CtForm.String():
		values, err := url.ParseQuery(string(b))
		if err != nil {
			return ae.New(ae.CodeBadRequest, "invalid form data")
		}
		r.setPartialBodyData(values)
	default:
		if err := json.Unmarshal(b, &r.partialBodyData); err != nil {
			return ae.New(ae.CodeBadRequest, "invalid json")
		}
	}
	// @see http.ParseMultipartForm
	// .MultipartForm "multipart/form-data" ||  "multipart/mixed"
	// .PostFormValue  "application/x-www-form-urlencoded" url.ParseQuery(body) + .MultipartForm
	// .FormValue() 调用 .Form  =  url.ParseQuery(r.URL.RawQuery) + .PostFormValue
	return nil
}

// parseMultipartBody 解析multipart请求体
func (r *Request) parseMultipartBody(boundary string) *ae.Error {
	if boundary == "" {
		return ae.New(ae.CodeBadRequest, "missing boundary")
	}
	form, err := multipart.NewReader(r.r.Body, boundary).ReadForm(maxMultiSize)
	if err != nil {
		return ae.New(ae.CodeBadRequest, "body should encode in multipart form")
	}
	r.setPartialBodyData(form.Value)
	return nil
}

func (r *Request) Body(name string, patterns ...any) (*RawValue, *ae.Error) {
	if !r.bodyParsed {
		err := r.parseBodyStream()
		if err != nil {
			return nil, err
		}
	}
	raw := newRawValue(name, "")
	if r.partialBodyData != nil {
		if v, ok := r.partialBodyData[name]; ok {
			raw.Reload(v)
		}
	}
	return raw, raw.Filter(patterns...)
}
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	return r.r.Cookie(name)
}

func (r *Request) AddCookie(c *http.Cookie) {
	r.r.AddCookie(c)
}

func (r *Request) Cookies() []*http.Cookie {
	return r.r.Cookies()
}

// @warn 尽量不要通过自定义header传参，因为可能某个web server会基于安全禁止某些无法识别的header
func (r *Request) GetHeader(name string, patterns ...any) (*RawValue, *ae.Error) {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.query(r.partialHeaders, userData, name, patterns...)
}

func (r *Request) Header(name string) *RawValue {
	value, _ := r.GetHeader(name)
	return value
}

// HeaderWild 读取 HTTP Header（包括标准格式和 X- 前缀格式）
//  1. 原始格式 如 name, Name, user_agent
//  2. 标准格式  如 Referer, User-Agent,
//  3. X-前缀格式  如 X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
//
// @warn 尽量不要通过自定义header传参，因为可能某个web server会基于安全禁止某些无法识别的header
func (r *Request) HeaderWild(name string) *RawValue {
	// 1. 原始格式
	value := r.Header(name)
	if value.NotEmpty() {
		return value
	}

	// 	2. 标准格式
	key := cases.Title(language.English).String(strings.ReplaceAll(name, "_", "-"))
	if key != name {
		value = r.Header(key)
		if value.NotEmpty() {
			return value
		}
	}

	// 3. X-前缀格式
	return r.Header("X-" + key)
}

func (r *Request) UserAgent() string {
	ua := r.userAgent
	if ua == "" {
		ua = r.Header("User-Agent").String()
		r.userAgent = ua
	}
	return ua
}
func (r *Request) Accept() string {
	return r.Header("Accept").String()
}

// Headers 获取所有headers
// 这个读取少，直接每次独立解析即可
func (r *Request) Headers() map[string]any {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.queries(r.partialHeaders, userData)
}

// QueryWild 尝试从URL参数、Header（包括标准格式和X-前缀格式）、Cookie读取参数值
//  1. URL参数
//  2. HTTP头部 (支持标准和X-前缀格式)
//  3. Cookie
//
// e.g.  csrf_token: in url params? -> Csrf-Token: in header?  X-Csrf-Token: in header-> csrf_token: in cookie
func (r *Request) QueryWild(name string, patterns ...any) (*RawValue, *ae.Error) {
	// 1. URL参数直接查询模式
	v, e := r.Query(name)
	if e == nil && v.NotEmpty() {
		return v, v.Filter(patterns...)
	}
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	key := strings.ToLower(strings.ReplaceAll(name, "-", "_"))
	if key != name {
		v, e = r.Query(key, patterns...)
		if e == nil && v.NotEmpty() {
			return v, v.Filter(patterns...)
		}
	}

	// 2. HTTP头部（包括标准格式和 X- 前缀）
	v = r.HeaderWild(name)
	if v.NotEmpty() {
		return v, v.Filter(patterns...)
	}

	// 3. Cookie
	if cookie, err := r.Cookie(name); err == nil {
		v = newRawValue(cookie.Name, cookie.Value)
		return v, v.Filter(patterns...)
	}
	if key != name {
		if cookie, err := r.Cookie(key); err == nil {
			v = newRawValue(cookie.Name, cookie.Value)
			return v, v.Filter(patterns...)
		}
	}
	// 返回空值
	return v, v.Filter(patterns...)
}
func (r *Request) Wild(name string) *RawValue {
	v, _ := r.QueryWild(name)
	return v
}
