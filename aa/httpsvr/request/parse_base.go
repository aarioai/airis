package request

import (
	"encoding/json"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/utils"
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

type userDataInterface interface {
	Get(key string) string
}

// RawValue
// @extend io.Closer
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

func (v *RawValue) Close() error {
	if v.Atype != nil {
		return v.Atype.Close()
	}
	return nil
}
func errorOnEmpty(p string, required bool) *ae.Error {
	if required {
		return ae.NewBadParam(p)
	}
	return nil
}

func findDefaultValue(patterns []any) any {
	// 处理默认值
	for _, pat := range patterns {
		switch defaultVal := pat.(type) {
		case *atype.Atype:
			return defaultVal.Raw()
		case atype.Atype:
			return defaultVal.Raw()
		}
	}
	return nil
}
func parseValidationRules(name string, patterns []any) (bool, *regexp.Regexp, *ae.Error) {
	required := true // 默认为true
	var re *regexp.Regexp
	// 解析验证规则
	for _, pat := range patterns {
		switch p := pat.(type) {
		case string:
			if p != "" && re == nil {
				var err error
				if re, err = regexp.Compile(p); err != nil {
					return false, nil, ae.NewBadParam(name, fmt.Sprintf("invalid request string pattern `%s`", p))
				}
			}
		case bool:
			required = p
		case *regexp.Regexp:
			re = p
		case *atype.Atype:
		case atype.Atype:
		default:
			return false, nil, ae.NewBadParam(name, fmt.Sprintf("invalid request pattern `%s`", p))
		}
	}
	return required, re, nil
}
func (v *RawValue) ReleaseValidate(patterns []any) *ae.Error {
	defer v.Close()
	return v.Validate(patterns)
}

// Validate 验证并过滤值
// @param pattern  e.g. `[[:word:]]+` `\w+`
// Filter(pattern string, required bool)
// Filter(required bool)
// Filter(pattern string)
// Filter(default atype.Atype)
func (v *RawValue) Validate(patterns []any) *ae.Error {
	if v.IsEmpty() {
		if defaultVal := findDefaultValue(patterns); defaultVal != nil {
			v.Reload(defaultVal)
		}
	}
	required, re, e := parseValidationRules(v.name, patterns)
	if e != nil {
		return e
	}

	if required && v.IsEmpty() {
		return ae.NewBadParam(v.name)
	}

	if re != nil && !re.MatchString(v.String()) {
		return ae.NewBadParam(v.name)
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

func (r *Request) findAny(programData map[string]any, userData userDataInterface, key string) any {
	// 1. 优先读取程序设置的数据
	if programData != nil {
		if v, exists := programData[key]; exists {
			return v
		}
	}

	if userData != nil {
		return userData.Get(key)
	}
	return nil
}
func (r *Request) findStringFast(programData map[string]any, userData userDataInterface, key string) string {
	v := r.findAny(programData, userData, key)
	return atype.String(v)
}
func (r *Request) find(programData map[string]any, userData userDataInterface, key string, patterns ...any) (*RawValue, *ae.Error) {
	v := r.findAny(programData, userData, key)
	raw := newRawValue(key, v)
	e := raw.Validate(patterns)
	if e != nil {
		raw.Close()
		return nil, e
	}
	return raw, nil
}

func (r *Request) findString(programData map[string]any, userData userDataInterface, key string, patterns ...any) (string, *ae.Error) {
	v := r.findStringFast(programData, userData, key)
	only, requiredWhenOnly := onlyRequired(patterns)
	if only {
		if requiredWhenOnly && v == "" {
			return "", ae.NewBadParam(key)
		}
		return v, nil
	}
	return v, newRawValue(key, v).ReleaseValidate(patterns)
}
func (r *Request) queryString(name string, patterns ...any) (string, *ae.Error) {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.findString(r.injectedQueries, userData, name, patterns...)
}
func (r *Request) query(name string, patterns ...any) (*RawValue, *ae.Error) {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.find(r.injectedQueries, userData, name, patterns...)
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
		for k, v := range r.injectedHeaders {
			data[k] = v
		}
	}
	if len(data) > 0 {
		return data
	}
	return nil
}

func (r *Request) Query(key string, patterns ...any) (*RawValue, *ae.Error) {
	return r.query(key, patterns...)
}

func (r *Request) Queries() map[string]any {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.queries(r.injectedQueries, userData)
}
func (r *Request) setPartialBodyData(data map[string][]string) {
	if len(data) == 0 {
		return
	}
	if r.injectedBodies == nil {
		r.injectedBodies = make(map[string]any, len(data))
	}
	for k, v := range data {
		// 只返回首个元素
		// @see http.Request.FormValue()
		if len(v) > 0 {
			r.injectedBodies[k] = v[0]
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
		return ae.NewUnsupportedMedia()
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
	if _, ok := reader.(*maxBytesReader); !ok {
		reader = io.LimitReader(r.r.Body, r.maxFormSize+1)
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return ae.NewError(err)
	}
	if int64(len(b)) > r.maxFormSize {
		return ae.ErrorRequestEntityTooLarge
	}
	switch r.ContentType() {
	case CtForm.String():
		return r.parseFormBody(b)
	default:
		return r.parseJSONBody(b)
	}
	// @see http.ParseMultipartForm
	// .MultipartForm "multipart/form-data" ||  "multipart/mixed"
	// .PostFormValue  "application/x-www-form-urlencoded" url.ParseQuery(body) + .MultipartForm
	// .FormValue() 调用 .Form  =  url.ParseQuery(r.URL.RawQuery) + .PostFormValue
}

func (r *Request) parseFormBody(b []byte) *ae.Error {
	values, err := url.ParseQuery(string(b))
	if err != nil {
		return ae.NewUnsupportedMedia("form data")
	}
	if len(values) > 0 {
		r.setPartialBodyData(values)
	}
	return nil
}

func (r *Request) parseJSONBody(b []byte) *ae.Error {
	if err := json.Unmarshal(b, &r.injectedBodies); err != nil {
		return ae.NewUnsupportedMedia("json")
	}
	return nil
}

// parseMultipartBody 解析multipart请求体
func (r *Request) parseMultipartBody(boundary string) *ae.Error {
	if boundary == "" {
		return ae.NewPreconditionFailed("missing boundary")
	}

	form, err := multipart.NewReader(r.r.Body, boundary).ReadForm(r.maxMultipartSize)
	if err != nil {
		return ae.NewUnsupportedMedia("multipart form")
	}
	if len(form.Value) > 0 {
		r.setPartialBodyData(form.Value)
	}
	if len(form.File) > 0 {
		r.injectedFiles = form.File
	}
	return nil
}

func (r *Request) Accept() string {
	return r.HeaderFast("Accept")
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
func (r *Request) UserAgent() string {
	ua := r.userAgent
	if ua == "" {
		ua = r.HeaderFast("User-Agent")
		r.userAgent = ua
	}
	return ua
}
func (r *Request) Body(key string, patterns ...any) (*RawValue, *ae.Error) {
	raw := newRawValue(key, "")
	if !r.bodyParsed {
		e := r.parseBodyStream()
		if e != nil {
			raw.Close()
			return nil, e
		}
	}
	if r.injectedBodies != nil {
		if v, ok := r.injectedBodies[key]; ok {
			raw.Reload(v)
		}
	}
	return raw, raw.Validate(patterns)
}
func (r *Request) Cookie(key string) (*http.Cookie, error) {
	return r.r.Cookie(key)
}

func (r *Request) AddCookie(c *http.Cookie) {
	r.r.AddCookie(c)
}

func (r *Request) Cookies() []*http.Cookie {
	return r.r.Cookies()
}

func (r *Request) headerString(key string, patterns ...any) (string, *ae.Error) {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.findString(r.injectedHeaders, userData, key, patterns...)
}
func (r *Request) header(key string, patterns ...any) (*RawValue, *ae.Error) {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.find(r.injectedHeaders, userData, key, patterns...)
}

func (r *Request) Header(key string, patterns ...any) (*RawValue, *ae.Error) {
	return r.header(key, patterns...)
}
func (r *Request) HeaderFast(key string) string {
	// false 是必须的，表示 required=false。默认 required = true
	v, _ := r.headerString(key, false)
	return v
}

// HeaderString 查询更高效
func (r *Request) HeaderString(key string, patterns ...any) (string, *ae.Error) {
	// false 是必须的，表示
	return r.headerString(key, patterns...)
}

func (r *Request) HeaderValue(key string) *RawValue {
	value, _ := r.Header(key)
	return value
}

// HeaderWild read http header wildly (including standard or non-standard format)
//  1. origin format, e.g. name, Name, user_agent
//  2. standard format, e.g. Referer, User-Agent
//  3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
//
// Suggest any non-standard header should be allowed by the web server cors AllowedMethods
func (r *Request) HeaderWild(key string, patterns ...any) (*RawValue, *ae.Error) {
	// 1. origin format, e.g. key, Name, user_agent
	value, e := r.Header(key, patterns...)
	if e == nil && value.NotEmpty() {
		return value, nil
	}
	value.Close()

	// 2. standard format, e.g. Referer, User-Agent
	newKey := cases.Title(language.English).String(strings.ReplaceAll(key, "_", "-"))
	if newKey != key {
		value, e = r.Header(newKey, patterns...)
		if e == nil && value.NotEmpty() {
			return value, nil
		}
		value.Close()
	}
	if strings.HasPrefix(newKey, "X-") {
		return nil, validateEmpty(key, patterns)
	}
	// 3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
	return r.Header("X-"+newKey, patterns...)
}

func (r *Request) HeaderWideString(key string, patterns ...any) (string, *ae.Error) {
	header, e := r.HeaderWild(key, patterns...)
	if e != nil {
		return "", e
	}
	return header.ReleaseString(), nil
}

// HeaderWildFast faster then HeaderWild
func (r *Request) HeaderWildFast(key string) string {
	// 1. origin format, e.g. key, Name, user_agent
	value := r.HeaderFast(key)
	if value != "" {
		return value
	}
	// 2. standard format, e.g. Referer, User-Agent
	newKey := cases.Title(language.English).String(strings.ReplaceAll(key, "_", "-"))
	if newKey != key {
		if value = r.HeaderFast(newKey); value != "" {
			return value
		}
	}

	// 3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
	if !strings.HasPrefix(newKey, "X-") {
		return r.HeaderFast("X-" + newKey)
	}
	return ""
}

// Headers get all headers
func (r *Request) Headers() map[string]any {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.queries(r.injectedHeaders, userData)
}

// QueryWild query parameter from URL parameter, URL query, header and cookie
// Example  csrf_token: in url params? -> Csrf-Token: in header?  X-Csrf-Token: in header-> csrf_token: in cookie
func (r *Request) QueryWild(key string, patterns ...any) (*RawValue, *ae.Error) {
	// 1. URL参数直接查询模式
	v, e := r.Query(key, patterns...)
	if e == nil && v.NotEmpty() {
		return v, nil
	}
	utils.Close(v)
	guessName := key
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	if strings.HasPrefix(key, "X-") {
		guessName = strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "X-"), "-", "_"))
		if guessName != key {
			v, e = r.Query(guessName, patterns...)
			if e == nil && v.NotEmpty() {
				return v, nil
			}
			utils.Close(v)
		}
	}

	// 2. HTTP头部（包括标准格式和 X- 前缀）
	v, e = r.HeaderWild(key, patterns...)
	if e == nil && v.NotEmpty() {
		return v, nil
	}
	utils.Close(v)

	// 3. Cookie
	if cookie, err := r.Cookie(key); err == nil && cookie.Value != "" {
		v = newRawValue(cookie.Name, cookie.Value)
	} else if guessName != key {
		if cookie, err = r.Cookie(guessName); err == nil && cookie.Value != "" {
			v = newRawValue(cookie.Name, cookie.Value)
		}
	}
	if v == nil {
		v = newRawValue(key, "")
	}
	e = v.Validate(patterns) // 空值也需要判断是否符合pattern，如 required=false
	if e != nil {
		v.Close()
		return nil, e
	}
	return v, nil
}

func (r *Request) QueryWildString(key string, patterns ...any) (string, *ae.Error) {
	// 1. URL参数直接查询模式
	v, e := r.QueryString(key, patterns...)
	if e == nil && v != "" {
		return v, nil
	}
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	guessName := key
	if strings.HasPrefix(key, "X-") {
		guessName = strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "X-"), "-", "_"))
		if guessName != key {
			v, e = r.QueryString(guessName, patterns...)
			if e == nil && v != "" {
				return v, nil
			}
		}
	}
	// 2. HTTP头部（包括标准格式和 X- 前缀）
	v, e = r.HeaderWideString(key, patterns...)
	if e == nil && v != "" {
		return v, nil
	}

	// 3. Cookie
	if cookie, err := r.Cookie(key); err == nil {
		return cookie.Value, newRawValue(cookie.Name, cookie.Value).ReleaseValidate(patterns)
	}
	if guessName != key {
		if cookie, err := r.Cookie(guessName); err == nil {
			return cookie.Value, newRawValue(cookie.Name, cookie.Value).ReleaseValidate(patterns)
		}
	}
	// 返回空值
	return v, newRawValue(key, v).ReleaseValidate(patterns)
}

// QueryWildFast 更高效地快速查询字符串
func (r *Request) QueryWildFast(key string) string {
	// 1. URL参数直接查询模式
	if v := r.QueryFast(key); v != "" {
		return v
	}
	// 1.1. URL参数替换格式查询，可能使用的是Header（大写开头）参数名，改为小写下划线模式
	guessName := key
	if strings.HasPrefix(key, "X-") {
		guessName = strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "X-"), "-", "_"))
		if guessName != key {
			if v := r.QueryFast(guessName); v != "" {
				return v
			}
		}
	}

	// 2. HTTP头部（包括标准格式和 X- 前缀）
	if v := r.HeaderWildFast(key); v != "" {
		return v
	}

	// 3. Cookie
	if cookie, err := r.Cookie(key); err == nil {
		return cookie.Value
	}
	if guessName != key {
		if cookie, err := r.Cookie(guessName); err == nil {
			return cookie.Value
		}
	}
	// 返回空值
	return ""
}
func (r *Request) QueryWildValue(key string) *RawValue {
	v, err := r.QueryWild(key)
	if err != nil {
		return newRawValue(key, "")
	}
	return v
}

// isRequired
// 1. 不传参数等同于传了 true，即 required = true
// 2. 其他的以第一个为准
func isRequired(requireds []bool) bool {
	// 1. 不传参数等同于传了 true，即 required = true
	if len(requireds) == 0 {
		return true
	}
	return requireds[0]
}

// onlyRequired 参数是否只包含 required
// 1. 不传参数等同于传了 true，即 required = true
// 2. 只传了一个参数，且该参数为bool类型，则该参数值为required值
func onlyRequired(patterns []any) (only bool, required bool) {
	// 1. 不传参数等同于传了 true，即 required = true
	if len(patterns) == 0 {
		return true, true
	}
	if len(patterns) == 1 {
		if required, ok := patterns[0].(bool); ok {
			return true, required
		}
	}
	return false, false
}
func validateEmpty(key string, patterns []any) *ae.Error {
	return newRawValue(key, "").ReleaseValidate(patterns)
}
