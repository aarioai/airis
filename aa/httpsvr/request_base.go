package httpsvr

import (
	json "encoding/json/v2"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// englishTitle is stateless and safe to reuse concurrently. Avoid building a caser on every header read.
var englishTitle = cases.Title(language.English)

func titleHeaderKey(key string) string {
	if strings.ContainsRune(key, '_') {
		key = strings.ReplaceAll(key, "_", "-")
	}
	return englishTitle.String(key)
}

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
	// Handle the default value.
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
	required := true // default is true
	var re *regexp.Regexp
	// Parse validation rules.
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

// Validate checks and filters the value.
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

func (c *Context) contentMediaType() (mediatype string, params map[string]string, err error) {
	// @see http.parsePostForm()
	ct := c.r.Header.Get("Content-Type")
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if ct == "" {
		ct = CtOctetStream.String()
	}
	return mime.ParseMediaType(ct)
}

// Origin e.g. https://luexu.com  Port 80 omits the port; any other port includes :$port.
// ${scheme}//${host}[:$port]
func (c *Context) Origin() string {
	scheme := c.r.URL.Scheme
	if scheme == "" {
		if c.r.TLS != nil {
			scheme = "https:"
		} else {
			scheme = "http:"
		}
	}
	// A redirect through this API is not enough: it only redirects this response, not the page. The client must handle that itself.
	host := c.r.Host
	h := scheme + "//" + host
	return h
}

func (c *Context) findAny(programData map[string]any, userData userDataInterface, key string) any {
	// 1. Prefer data set by the program.
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
func (c *Context) findStringFast(programData map[string]any, userData userDataInterface, key string) string {
	v := c.findAny(programData, userData, key)
	return atype.String(v)
}
func (c *Context) find(programData map[string]any, userData userDataInterface, key string, patterns ...any) (*RawValue, *ae.Error) {
	v := c.findAny(programData, userData, key)
	raw := newRawValue(key, v)
	e := raw.Validate(patterns)
	if e != nil {
		raw.Close()
		return nil, e
	}
	return raw, nil
}

func (c *Context) findString(programData map[string]any, userData userDataInterface, key string, patterns ...any) (string, *ae.Error) {
	v := c.findStringFast(programData, userData, key)
	only, requiredWhenOnly := onlyRequired(patterns)
	if only {
		if requiredWhenOnly && v == "" {
			return "", ae.NewBadParam(key)
		}
		return v, nil
	}
	return v, newRawValue(key, v).ReleaseValidate(patterns)
}
func (c *Context) queryString(name string, patterns ...any) (string, *ae.Error) {
	c.ensurePathParams()
	var userData url.Values
	if c.r != nil {
		userData = c.r.URL.Query()
	}
	return c.findString(c.injectedQueries, userData, name, patterns...)
}
func (c *Context) query(name string, patterns ...any) (*RawValue, *ae.Error) {
	c.ensurePathParams()
	var userData url.Values
	if c.r != nil {
		userData = c.r.URL.Query()
	}
	return c.find(c.injectedQueries, userData, name, patterns...)
}

func (c *Context) queries(programData map[string]any, userData map[string][]string) map[string]any {
	if len(userData) == 0 && len(programData) == 0 {
		return nil
	}
	data := make(map[string]any, len(userData)+len(programData))
	// Client values are stored first; a program-injected parameter of the same name overwrites them.
	for k, vs := range userData {
		if len(vs) > 0 && vs[0] != "" {
			data[k] = vs[0]
		}
	}
	for k, v := range programData {
		data[k] = v
	}
	if len(data) == 0 {
		return nil
	}
	return data
}

func (c *Context) Query(key string, patterns ...any) (*RawValue, *ae.Error) {
	return c.query(key, patterns...)
}

func (c *Context) Queries() map[string]any {
	c.ensurePathParams()
	var userData url.Values
	if c.r != nil {
		userData = c.r.URL.Query()
	}
	return c.queries(c.injectedQueries, userData)
}
func (c *Context) setPartialBodyData(data map[string][]string) {
	if len(data) == 0 {
		return
	}
	if c.injectedBodies == nil {
		c.injectedBodies = make(map[string]any, len(data))
	}
	for k, v := range data {
		// Only return first
		// @see http.Request.FormValue()
		if len(v) > 0 {
			c.injectedBodies[k] = v[0]
		}
	}
}

// parseBodyStream
// @see http.parsePostForm()
func (c *Context) parseBodyStream() *ae.Error {
	defer func() { c.bodyParsed = true }()
	// Body can be nil
	if c.r == nil || c.r.Body == nil {
		return nil
	}

	ct, params, err := c.contentMediaType()
	if err != nil {
		return ae.NewUnsupportedMedia()
	}
	switch ct {
	case "", CtJSON.String(), CtOctetStream.String(), CtForm.String():
		return c.parseSimpleBody()
	case CtFormData.String():
		return c.parseMultipartBody(params["boundary"])
	}
	return nil
}

func (c *Context) parseSimpleBody() *ae.Error {
	var reader io.Reader = c.r.Body
	if _, ok := reader.(*maxBytesReader); !ok {
		reader = io.LimitReader(c.r.Body, c.maxFormSize+1)
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return ae.NewErr(err)
	}
	if int64(len(b)) > c.maxFormSize {
		return ae.ErrorRequestEntityTooLarge
	}
	switch c.ReqContentType() {
	case CtForm.String():
		return c.parseFormBody(b)
	default:
		return c.parseJSONBody(b)
	}
	// @see http.ParseMultipartForm
	// .MultipartForm "multipart/form-data" ||  "multipart/mixed"
	// .PostFormValue  "application/x-www-form-urlencoded" url.ParseQuery(body) + .MultipartForm
	// .FormValue() calls .Form  =  url.ParseQuery(c.URL.RawQuery) + .PostFormValue
}

func (c *Context) parseFormBody(b []byte) *ae.Error {
	values, err := url.ParseQuery(string(b))
	if err != nil {
		return ae.NewUnsupportedMedia("form data")
	}
	if len(values) > 0 {
		c.setPartialBodyData(values)
	}
	return nil
}

func (c *Context) parseJSONBody(b []byte) *ae.Error {
	if err := json.Unmarshal(b, &c.injectedBodies); err != nil {
		return ae.NewUnsupportedMedia("json")
	}
	return nil
}

func (c *Context) parseMultipartBody(boundary string) *ae.Error {
	if boundary == "" {
		return ae.NewPreconditionFailed("missing boundary")
	}

	form, err := multipart.NewReader(c.r.Body, boundary).ReadForm(c.maxMultipartSize)
	if err != nil {
		return ae.NewUnsupportedMedia("multipart form")
	}
	if len(form.Value) > 0 {
		c.setPartialBodyData(form.Value)
	}
	if len(form.File) > 0 {
		c.injectedFiles = form.File
	}
	return nil
}

func (c *Context) Accept() string {
	return c.HeaderFast("Accept")
}

// ContentType
// Each request has its own Request, so concurrency is rare and locking is unnecessary.
func (c *Context) ReqContentType() string {
	if c.contentType != "" {
		return c.contentType
	}
	c.contentType, _, _ = c.contentMediaType()
	return c.contentType
}

func (c *Context) UserAgent() string {
	ua := c.userAgent
	if ua == "" {
		ua = c.HeaderFast("User-Agent")
		c.userAgent = ua
	}
	return ua
}
func (c *Context) Body(key string, patterns ...any) (*RawValue, *ae.Error) {
	raw := newRawValue(key, "")
	if !c.bodyParsed {
		e := c.parseBodyStream()
		if e != nil {
			raw.Close()
			return nil, e
		}
	}
	if c.injectedBodies != nil {
		if v, ok := c.injectedBodies[key]; ok {
			raw.Reload(v)
		}
	}
	return raw, raw.Validate(patterns)
}
func (c *Context) Cookie(key string) (*http.Cookie, error) {
	return c.r.Cookie(key)
}

func (c *Context) AddCookie(cookie *http.Cookie) {
	c.r.AddCookie(cookie)
}

func (c *Context) Cookies() []*http.Cookie {
	return c.r.Cookies()
}

func (c *Context) headerString(key string, patterns ...any) (string, *ae.Error) {
	var userData http.Header
	if c.r != nil {
		userData = c.r.Header
	}
	return c.findString(c.injectedHeaders, userData, key, patterns...)
}
func (c *Context) header(key string, patterns ...any) (*RawValue, *ae.Error) {
	var userData http.Header
	if c.r != nil {
		userData = c.r.Header
	}
	return c.find(c.injectedHeaders, userData, key, patterns...)
}

func (c *Context) HeaderParam(key string, patterns ...any) (*RawValue, *ae.Error) {
	return c.HeaderParam(key, patterns...)
}
func (c *Context) HeaderFast(key string) string {
	// false is required: it means required=false. The default required is true.
	v, _ := c.headerString(key, false)
	return v
}

// HeaderString looks up a header more efficiently.
func (c *Context) HeaderString(key string, patterns ...any) (string, *ae.Error) {
	// false is required, indicating
	return c.headerString(key, patterns...)
}

func (c *Context) HeaderValue(key string) *RawValue {
	value, _ := c.HeaderParam(key)
	return value
}

// HeaderWild read http header wildly (including standard or non-standard format)
//  1. origin format, e.g. name, Name, user_agent
//  2. standard format, e.g. Referer, User-Agent
//  3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
//
// Suggest any non-standard header should be allowed by the web server cors AllowedMethods
func (c *Context) HeaderWild(key string, patterns ...any) (*RawValue, *ae.Error) {
	// 1. origin format, e.g. key, Name, user_agent
	value, e := c.HeaderParam(key, patterns...)
	if e == nil && value.NotEmpty() {
		return value, nil
	}
	value.Close()

	// 2. standard format, e.g. Referer, User-Agent
	newKey := titleHeaderKey(key)
	if newKey != key {
		value, e = c.HeaderParam(newKey, patterns...)
		if e == nil && value.NotEmpty() {
			return value, nil
		}
		value.Close()
	}
	if strings.HasPrefix(newKey, "X-") {
		return nil, validateEmpty(key, patterns)
	}
	// 3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
	return c.HeaderParam("X-"+newKey, patterns...)
}

func (c *Context) HeaderWideString(key string, patterns ...any) (string, *ae.Error) {
	header, e := c.HeaderWild(key, patterns...)
	if e != nil {
		return "", e
	}
	return header.ReleaseString(), nil
}

// HeaderWildFast faster then HeaderWild
func (c *Context) HeaderWildFast(key string) string {
	// 1. origin format, e.g. key, Name, user_agent
	value := c.HeaderFast(key)
	if value != "" {
		return value
	}
	// 2. standard format, e.g. Referer, User-Agent
	newKey := titleHeaderKey(key)
	if newKey != key {
		if value = c.HeaderFast(newKey); value != "" {
			return value
		}
	}

	// 3. self-defined format, i.e. starts with X-, e.g. X-Csrf-Token, X-Request-Vuid, X-From, X-Inviter
	if !strings.HasPrefix(newKey, "X-") {
		return c.HeaderFast("X-" + newKey)
	}
	return ""
}

// Headers get all headers
func (c *Context) Headers() map[string]any {
	var userData http.Header
	if c.r != nil {
		userData = c.r.Header
	}
	return c.queries(c.injectedHeaders, userData)
}

// QueryWild query parameter from URL parameter, URL query, header and cookie
// Example  csrf_token: in url params? -> Csrf-Token: in header?  X-Csrf-Token: in header-> csrf_token: in cookie
func (c *Context) QueryWild(key string, patterns ...any) (*RawValue, *ae.Error) {
	// 1. Look up the URL parameter directly.
	v, e := c.Query(key, patterns...)
	if e == nil && v.NotEmpty() {
		return v, nil
	}
	utils.Close(v)
	guessName := key
	// 1.1. Alternate URL-parameter form. A Header-style name (leading capitals) is rewritten to lowercase with underscores.
	if strings.HasPrefix(key, "X-") {
		guessName = strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "X-"), "-", "_"))
		if guessName != key {
			v, e = c.Query(guessName, patterns...)
			if e == nil && v.NotEmpty() {
				return v, nil
			}
			utils.Close(v)
		}
	}

	// 2. HTTP headers, including the standard form and the X- prefix.
	v, e = c.HeaderWild(key, patterns...)
	if e == nil && v.NotEmpty() {
		return v, nil
	}
	utils.Close(v)

	// 3. Cookie
	if cookie, err := c.Cookie(key); err == nil && cookie.Value != "" {
		v = newRawValue(cookie.Name, cookie.Value)
	} else if guessName != key {
		if cookie, err = c.Cookie(guessName); err == nil && cookie.Value != "" {
			v = newRawValue(cookie.Name, cookie.Value)
		}
	}
	if v == nil {
		v = newRawValue(key, "")
	}
	e = v.Validate(patterns) // an empty value must still be checked against the pattern, for example required=false
	if e != nil {
		v.Close()
		return nil, e
	}
	return v, nil
}

func (c *Context) QueryWildString(key string, patterns ...any) (string, *ae.Error) {
	// 1. Look up the URL parameter directly.
	v, e := c.QueryString(key, patterns...)
	if e == nil && v != "" {
		return v, nil
	}
	// 1.1. Alternate URL-parameter form. A Header-style name (leading capitals) is rewritten to lowercase with underscores.
	guessName := key
	if strings.HasPrefix(key, "X-") {
		guessName = strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "X-"), "-", "_"))
		if guessName != key {
			v, e = c.QueryString(guessName, patterns...)
			if e == nil && v != "" {
				return v, nil
			}
		}
	}
	// 2. HTTP headers, including the standard form and the X- prefix.
	v, e = c.HeaderWideString(key, patterns...)
	if e == nil && v != "" {
		return v, nil
	}

	// 3. Cookie
	if cookie, err := c.Cookie(key); err == nil {
		return cookie.Value, newRawValue(cookie.Name, cookie.Value).ReleaseValidate(patterns)
	}
	if guessName != key {
		if cookie, err := c.Cookie(guessName); err == nil {
			return cookie.Value, newRawValue(cookie.Name, cookie.Value).ReleaseValidate(patterns)
		}
	}
	// Return the empty value.
	return v, newRawValue(key, v).ReleaseValidate(patterns)
}

// QueryWildFast looks up a string more efficiently.
func (c *Context) QueryWildFast(key string) string {
	// 1. Look up the URL parameter directly.
	if v := c.QueryFast(key); v != "" {
		return v
	}
	// 1.1. Alternate URL-parameter form. A Header-style name (leading capitals) is rewritten to lowercase with underscores.
	guessName := key
	if strings.HasPrefix(key, "X-") {
		guessName = strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "X-"), "-", "_"))
		if guessName != key {
			if v := c.QueryFast(guessName); v != "" {
				return v
			}
		}
	}

	// 2. HTTP headers, including the standard form and the X- prefix.
	if v := c.HeaderWildFast(key); v != "" {
		return v
	}

	// 3. Cookie
	if cookie, err := c.Cookie(key); err == nil {
		return cookie.Value
	}
	if guessName != key {
		if cookie, err := c.Cookie(guessName); err == nil {
			return cookie.Value
		}
	}
	// Return the empty value.
	return ""
}
func (c *Context) QueryWildValue(key string) *RawValue {
	v, err := c.QueryWild(key)
	if err != nil {
		return newRawValue(key, "")
	}
	return v
}

// isRequired
// 1. Omitting the argument is the same as passing true, so required = true.
// 2. Otherwise the first argument wins.
func isRequired(requireds []bool) bool {
	// 1. Omitting the argument is the same as passing true, so required = true.
	if len(requireds) == 0 {
		return true
	}
	return requireds[0]
}

// onlyRequired reports whether the arguments contain only required.
// 1. Omitting the argument is the same as passing true, so required = true.
// 2. If exactly one argument is passed and it is a bool, that value is required.
func onlyRequired(patterns []any) (only bool, required bool) {
	// 1. Omitting the argument is the same as passing true, so required = true.
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
