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

// onlyRequired 参数是否只包含 required
// 1. 不传参数等同于传了 true，即 required = true
// 2. 只传了一个参数，且该参数为bool类型，则该参数值为required值
func onlyRequired(patterns []any) (only bool, required bool) {
	switch len(patterns) {
	case 0:
		return true, true
	case 1:
		if r, ok := patterns[0].(bool); ok {
			return true, r
		}
	}
	return false, false
}

// Filter 验证并过滤值
// @param pattern  e.g. `[[:word:]]+` `\w+`
// Filter(pattern string, required bool)
// Filter(required bool)
// Filter(pattern string)
// Filter(default atype.Atype)
func (v *RawValue) Filter(patterns ...any) *ae.Error {
	// 处理默认值
	for _, pat := range patterns {
		switch defaultVal := pat.(type) {
		case *atype.Atype:
			if v.String() == "" {
				v.Reload(defaultVal.Raw())
				break
			}
		case atype.Atype:
			if v.String() == "" {
				v.Reload(defaultVal.Raw())
				break
			}
		}
	}

	var (
		required = true // 默认必需
		pattern  string
		re       *regexp.Regexp
	)

	// 解析验证规则
	for _, pat := range patterns {
		switch p := pat.(type) {
		case string:
			pattern = p
		case bool:
			required = p
		case *regexp.Regexp:
			re = p
		}
	}

	if v.String() == "" {
		if required {
			return ae.BadParam(v.name)
		}
		return nil
	}
	if pattern != "" && re == nil {
		var err error
		if re, err = regexp.Compile(pattern); err != nil {
			return ae.NewMsg("invalid request pattern `%s`", pattern)
		}
	}
	if re != nil && !re.MatchString(v.String()) {
		return ae.BadParam(v.name)
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

func (r *Request) findAny(programData map[string]any, userData userDataInterface, name string) any {
	// 读取程序设置的header
	if programData != nil {
		if v, exists := programData[name]; exists {
			return v
		}
	}

	if userData != nil {
		return userData.Get(name)
	}
	return ""
}
func (r *Request) findStringFast(programData map[string]any, userData userDataInterface, name string) string {
	v := r.findAny(programData, userData, name)
	if s, ok := v.(string); ok {
		return s
	}
	return atype.String(v)
}
func (r *Request) find(programData map[string]any, userData userDataInterface, name string, patterns ...any) (*RawValue, *ae.Error) {
	v := r.findAny(programData, userData, name)
	raw := newRawValue(name, v)
	return raw, raw.Filter(patterns...)
}

func (r *Request) findString(programData map[string]any, userData userDataInterface, name string, patterns ...any) (string, *ae.Error) {
	v := r.findStringFast(programData, userData, name)
	only, requiredWhenOnly := onlyRequired(patterns)
	if only {
		if requiredWhenOnly && v == "" {
			return "", ae.BadParam(name)
		}
		return v, nil
	}
	return v, newRawValue(name, v).Filter(patterns)
}
func (r *Request) queryString(name string, patterns ...any) (string, *ae.Error) {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.findString(r.partialQueries, userData, name, patterns...)
}
func (r *Request) query(name string, patterns ...any) (*RawValue, *ae.Error) {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.find(r.partialQueries, userData, name, patterns...)
}

// Headers 获取所有headers
// 这个读取少，直接每次独立解析即可
func (r *Request) queries(programData map[string]any, userData map[string][]string) map[string]any {
	data := make(map[string]any)
	// 优先级最低，读取用户data
	if programData != nil {
		for k, vs := range userData {
			if len(vs) != 0 && vs[0] != "" {
				data[k] = vs[0]
			}
		}
	}
	// 优先级高，读取程序设置的header
	if programData != nil {
		for k, v := range r.partialHeaders {
			data[k] = v
		}
	}
	if len(data) > 0 {
		return data
	}
	return nil
}

func (r *Request) Query(name string, patterns ...any) (*RawValue, *ae.Error) {
	return r.query(name, patterns...)
}

func (r *Request) Queries() map[string]any {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.queries(r.partialQueries, userData)
}
func (r *Request) setPartialBodyData(data map[string][]string) {
	if len(data) == 0 {
		return
	}
	if r.partialBodyData == nil {
		r.partialBodyData = make(map[string]any, len(data))
	}
	for k, v := range data {
		// 只返回首个元素
		// @see http.Request.FormValue()
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
		return ae.New(ae.CodeUnsupportedMedia, "invalid content media type")
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
		return ae.RequestEntityTooLarge
	}
	switch r.ContentType() {
	case CtForm.String():
		values, err := url.ParseQuery(string(b))
		if err != nil {
			return ae.New(ae.CodeUnsupportedMedia, "invalid form data")
		}
		r.setPartialBodyData(values)
	default:
		if err := json.Unmarshal(b, &r.partialBodyData); err != nil {
			return ae.New(ae.CodeUnsupportedMedia, "invalid json")
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
		return ae.New(ae.CodePreconditionFailed, "missing boundary")
	}

	form, err := multipart.NewReader(r.r.Body, boundary).ReadForm(maxMultiSize)
	if err != nil {
		return ae.New(ae.CodeUnsupportedMedia, "body should encode in multipart form")
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

func (r *Request) headerString(name string, patterns ...any) (string, *ae.Error) {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.findString(r.partialHeaders, userData, name, patterns...)
}
func (r *Request) header(name string, patterns ...any) (*RawValue, *ae.Error) {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.find(r.partialHeaders, userData, name, patterns...)
}

// @warn 尽量不要通过自定义header传参，因为可能某个web server会基于安全禁止某些无法识别的header
func (r *Request) Header(name string, patterns ...any) (*RawValue, *ae.Error) {
	return r.header(name, patterns...)
}
func (r *Request) HeaderFast(name string) string {
	// false 是必须的，表示 required=false。默认 required = true
	v, _ := r.headerString(name, false)
	return v
}

// HeaderString 查询更高效
func (r *Request) HeaderString(name string, patterns ...any) (string, *ae.Error) {
	// false 是必须的，表示
	return r.headerString(name, patterns...)
}

func (r *Request) HeaderValue(name string) *RawValue {
	value, _ := r.Header(name)
	return value
}

// HeaderWild 读取 HTTP HeaderValue（包括标准格式和 X- 前缀格式）
//  1. 原始格式 如 name, Name, user_agent
//  2. 标准格式  如 Referer, User-Agent,
//  3. X-前缀格式  如 X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
//
// @warn 尽量不要通过自定义header传参，因为可能某个web server会基于安全禁止某些无法识别的header
func (r *Request) HeaderWild(name string) *RawValue {
	// 1. 原始格式
	value := r.HeaderValue(name)
	if value.NotEmpty() {
		return value
	}

	// 	2. 标准格式
	key := cases.Title(language.English).String(strings.ReplaceAll(name, "_", "-"))
	if key != name {
		value = r.HeaderValue(key)
		if value.NotEmpty() {
			return value
		}
	}

	// 3. X-前缀格式
	return r.HeaderValue("X-" + key)
}

// HeaderWildFast 对 HeaderWild 的性能优化、类型简化
func (r *Request) HeaderWildFast(name string) string {
	// 1. 原始格式
	value := r.HeaderFast(name)
	if value != "" {
		return value
	}
	// 	2. 标准格式
	key := cases.Title(language.English).String(strings.ReplaceAll(name, "_", "-"))
	if key != name {
		value = r.HeaderFast(key)
		if value != "" {
			return value
		}
	}

	// 3. X-前缀格式
	return r.HeaderFast("X-" + key)
}

func (r *Request) UserAgent() string {
	ua := r.userAgent
	if ua == "" {
		ua = r.HeaderFast("User-Agent")
		r.userAgent = ua
	}
	return ua
}
func (r *Request) Accept() string {
	return r.HeaderFast("Accept")
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

// QueryWild 尝试从URL参数、HeaderValue（包括标准格式和X-前缀格式）、Cookie读取参数值
//  1. URL参数
//  2. HTTP头部 (支持标准和X-前缀格式)
//  3. Cookie
//
// e.g.  csrf_token: in url params? -> Csrf-Token: in header?  X-Csrf-Token: in header-> csrf_token: in cookie
func (r *Request) QueryWild(name string, patterns ...any) (*RawValue, *ae.Error) {
	// 1. URL参数直接查询模式
	v, e := r.Query(name, patterns...)
	if e == nil && v.NotEmpty() {
		return v, nil
	}
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	key := strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(name, "X-"), "-", "_"))
	if key != name {
		v, e = r.Query(key, patterns...)
		if e == nil && v.NotEmpty() {
			return v, nil
		}
	}

	// 2. HTTP头部（包括标准格式和 X- 前缀）
	v, e = r.Header(name, patterns...)
	if e == nil && v.NotEmpty() {
		return v, nil
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

func (r *Request) QueryWildString(name string, patterns ...any) (string, *ae.Error) {
	// 1. URL参数直接查询模式
	v, e := r.QueryString(name, patterns...)
	if e == nil && v != "" {
		return v, nil
	}
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	key := strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(name, "X-"), "-", "_"))
	if key != name {
		v, e = r.QueryString(key, patterns...)
		if e == nil && v != "" {
			return v, nil
		}
	}

	// 2. HTTP头部（包括标准格式和 X- 前缀）
	v, e = r.HeaderString(name, patterns...)
	if e == nil && v != "" {
		return v, nil
	}

	// 3. Cookie
	if cookie, err := r.Cookie(name); err == nil {
		return cookie.Value, newRawValue(cookie.Name, cookie.Value).Filter(patterns...)
	}
	if key != name {
		if cookie, err := r.Cookie(key); err == nil {
			return cookie.Value, newRawValue(cookie.Name, cookie.Value).Filter(patterns...)
		}
	}
	// 返回空值
	return v, newRawValue(name, v).Filter(patterns...)
}

// QueryWildFast 更高效地快速查询字符串
func (r *Request) QueryWildFast(name string) string {
	// 1. URL参数直接查询模式
	if v := r.QueryFast(name); v != "" {
		return v
	}
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	key := strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(name, "X-"), "-", "_"))
	if key != name {
		if v := r.QueryFast(key); v != "" {
			return v
		}
	}

	// 2. HTTP头部（包括标准格式和 X- 前缀）
	if v := r.HeaderFast(name); v != "" {
		return v
	}

	// 3. Cookie
	if cookie, err := r.Cookie(name); err == nil {
		return cookie.Value
	}
	if key != name {
		if cookie, err := r.Cookie(key); err == nil {
			return cookie.Value
		}
	}
	// 返回空值
	return ""
}
func (r *Request) QueryWildValue(name string) *RawValue {
	v, err := r.QueryWild(name)
	if err != nil {
		return newRawValue(name, "")
	}
	return v
}
