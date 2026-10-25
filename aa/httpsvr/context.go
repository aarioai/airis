package httpsvr

import (
	"context"
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/aarioai/airis/aa/acontext"
	ctxatype "github.com/aarioai/airis/aa/httpsvr/atype"
)

// RequestIDKey is the context value set by the RequestID middleware.
const RequestIDKey = acontext.CtxRequestID

// Handler is an HTTP handler in the router chain. Call Next to continue.
type Handler func(*Context)

// Param is one parsed path parameter.
type Param struct {
	Key      string
	ValueRaw any
}

// Params is the set of path parameters for one request.
type Params struct {
	Store []Param
}

func (p Params) Len() int { return len(p.Store) }

func (p Params) Get(key string) (any, bool) {
	for _, item := range p.Store {
		if item.Key == key {
			return item.ValueRaw, true
		}
	}
	return nil, false
}

func (p Params) GetString(key string) string {
	v, ok := p.Get(key)
	if !ok || v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	case int:
		return strconv.Itoa(s)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int32:
		return strconv.FormatInt(int64(s), 10)
	case int64:
		return strconv.FormatInt(s, 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case bool:
		if s {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprint(s)
	}
}

func (p Params) GetEntry(key string) (Param, bool) {
	for _, item := range p.Store {
		if item.Key == key {
			return item, true
		}
	}
	return Param{}, false
}

func (p Params) Visit(fn func(key string, value any)) {
	for _, item := range p.Store {
		fn(item.Key, item.ValueRaw)
	}
}

func (p Params) GetInt64(key string) (int64, bool) {
	v, ok := p.Get(key)
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case int8:
		return int64(n), true
	case int16:
		return int64(n), true
	case int32:
		return int64(n), true
	case uint64:
		return int64(n), true
	case uint:
		return int64(n), true
	case uint8:
		return int64(n), true
	case uint16:
		return int64(n), true
	case uint32:
		return int64(n), true
	case string:
		x, err := strconv.ParseInt(n, 10, 64)
		return x, err == nil
	default:
		return 0, false
	}
}

func (p Params) GetInt(key string) (int, bool) {
	n, ok := p.GetInt64(key)
	return int(n), ok
}

func (p Params) GetUint64(key string) (uint64, bool) {
	v, ok := p.Get(key)
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case uint64:
		return n, true
	case uint:
		return uint64(n), true
	case uint8:
		return uint64(n), true
	case uint16:
		return uint64(n), true
	case uint32:
		return uint64(n), true
	case int64:
		if n < 0 {
			return 0, false
		}
		return uint64(n), true
	case int:
		if n < 0 {
			return 0, false
		}
		return uint64(n), true
	case string:
		x, err := strconv.ParseUint(n, 10, 64)
		return x, err == nil
	default:
		return 0, false
	}
}

func (p Params) GetBool(key string) (bool, bool) {
	v, ok := p.Get(key)
	if !ok || v == nil {
		return false, false
	}
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		switch b {
		case "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes":
			return true, true
		case "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No":
			return false, true
		default:
			return false, false
		}
	default:
		return false, false
	}
}

// Values is the per-request bag previously stored on the web framework context.
type Values struct {
	m map[string]any
}

func (v *Values) Get(key string) (any, bool) {
	if v == nil || v.m == nil {
		return nil, false
	}
	val, ok := v.m[key]
	return val, ok
}

func (v *Values) GetString(key string) string {
	val, ok := v.Get(key)
	if !ok || val == nil {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}

func (v *Values) GetBoolDefault(key string, def bool) bool {
	val, ok := v.Get(key)
	if !ok || val == nil {
		return def
	}
	b, ok := val.(bool)
	if !ok {
		return def
	}
	return b
}

func (v *Values) Set(key string, value any) {
	if v.m == nil {
		v.m = make(map[string]any, 4)
	}
	v.m[key] = value
}

// Context is one HTTP request. It carries the response writer, path parameters,
// a value bag, middleware chain, request parsing state, and DTO response state.
type Context struct {
	w        http.ResponseWriter
	r        *http.Request
	params   Params
	values   Values
	handlers []Handler
	index    int
	status   int
	wrote    bool
	stopped  bool
	viewData map[string]any
	views    *template.Template
	hooks    any
	sw       statusWriter

	// std context (lazy via StdContext)
	stdCtx context.Context

	// request parsing (former httpsvr/request.Request)
	contentType      string
	userAgent        string
	bodyParsed       bool
	maxMultipartSize int64
	maxFormSize      int64
	injectedHeaders  map[string]any
	injectedQueries  map[string]any
	injectedBodies   map[string]any
	injectedFiles    map[string][]*multipart.FileHeader
	paramsInjected   bool

	// DTO response (former httpsvr/response.Writer)
	SerializeTag      string
	serveContentTypes []string
	beforeFlush       []func(*Context)
	beforeSerialize   []BeforeSerialize
	serialize         func(contentType string, d Body) ([]byte, string, error)
	errorHandler      ErrorHandler
	settings          *Settings
	outHeaders        map[string]string

	// atype / aenum parsers: ctx.QueryAtype.Sex(...) / ctx.BodyAtype.AudioPaths(...)
	// Implementation lives in package aa/httpsvr/atype.
	QueryAtype ctxatype.Query
	BodyAtype  ctxatype.Body
}

var contextPool = sync.Pool{
	New: func() any { return new(Context) },
}

func acquireContext(w http.ResponseWriter, r *http.Request, params Params, handlers []Handler, hooks any, views *template.Template) *Context {
	c := contextPool.Get().(*Context)
	c.w = w
	c.r = r
	c.params = params
	c.handlers = handlers
	c.index = -1
	c.status = 0
	c.wrote = false
	c.stopped = false
	c.viewData = nil
	c.views = views
	c.hooks = hooks
	c.sw.ctx = c
	c.stdCtx = nil
	c.contentType = ""
	c.userAgent = ""
	c.bodyParsed = false
	c.maxMultipartSize = 5 << 20
	c.maxFormSize = 1 << 20
	c.injectedHeaders = nil
	c.injectedQueries = nil
	c.injectedBodies = nil
	c.injectedFiles = nil
	c.paramsInjected = false
	c.SerializeTag = ""
	c.serveContentTypes = nil
	c.beforeFlush = nil
	c.beforeSerialize = nil
	c.serialize = nil
	c.errorHandler = nil
	c.settings = nil
	c.outHeaders = nil
	c.QueryAtype = ctxatype.Query{Src: c}
	c.BodyAtype = ctxatype.Body{Src: c}
	return c
}

func releaseContext(c *Context) {
	if c == nil {
		return
	}
	if c.values.m != nil {
		clear(c.values.m)
	}
	m := c.values.m
	*c = Context{values: Values{m: m}}
	contextPool.Put(c)
}

// StdContext returns the stdlib context bound to this request (lazy).
func (c *Context) StdContext() context.Context {
	if c.stdCtx == nil {
		c.stdCtx = acontext.Bind(c)
	}
	return c.stdCtx
}

func (c *Context) SetMaxMultipartSize(size int64) { c.maxMultipartSize = size }
func (c *Context) SetMaxFormSize(size int64)      { c.maxFormSize = size }

func (c *Context) InjectHeader(name string, value any) {
	if c.injectedHeaders == nil {
		c.injectedHeaders = make(map[string]any)
	}
	c.injectedHeaders[name] = value
}

func (c *Context) InjectQuery(name string, value any) {
	c.ensurePathParams()
	if c.injectedQueries == nil {
		c.injectedQueries = make(map[string]any)
	}
	c.injectedQueries[name] = value
}

func (c *Context) InjectBody(name string, value any) {
	if c.injectedBodies == nil {
		c.injectedBodies = make(map[string]any)
	}
	c.injectedBodies[name] = value
}

func (c *Context) ensurePathParams() {
	if c.paramsInjected {
		return
	}
	c.paramsInjected = true
	n := c.params.Len()
	if n == 0 {
		return
	}
	if c.injectedQueries == nil {
		c.injectedQueries = make(map[string]any, n)
	}
	for _, v := range c.params.Store {
		if _, exists := c.injectedQueries[v.Key]; !exists {
			c.injectedQueries[v.Key] = v.ValueRaw
		}
	}
}

func (c *Context) Request() *http.Request { return c.r }

func (c *Context) ResetRequest(r *http.Request) {
	if r != nil {
		c.r = r
	}
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	c.sw.ctx = c
	return &c.sw
}

func (c *Context) Method() string {
	if c.r == nil {
		return ""
	}
	return c.r.Method
}

func (c *Context) Path() string {
	if c.r == nil || c.r.URL == nil {
		return ""
	}
	return c.r.URL.Path
}

func (c *Context) Host() string {
	if c.r == nil {
		return ""
	}
	return c.r.Host
}

func (c *Context) FullRequestURI() string {
	if c.r == nil || c.r.URL == nil {
		return ""
	}
	return c.r.URL.RequestURI()
}

func (c *Context) URLParam(key string) string {
	if c.r == nil || c.r.URL == nil {
		return ""
	}
	return c.r.URL.Query().Get(key)
}

func (c *Context) FormValue(key string) string {
	if c.r == nil {
		return ""
	}
	return c.r.FormValue(key)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	http.SetCookie(c.w, cookie)
}

func (c *Context) Redirect(url string, statusCode ...int) {
	code := http.StatusFound
	if len(statusCode) > 0 && statusCode[0] > 0 {
		code = statusCode[0]
	}
	http.Redirect(c.ResponseWriter(), c.r, url, code)
	c.StopExecution()
}

func (c *Context) SetContentType(ct string) {
	if ct != "" {
		c.Header().Set("Content-Type", ct)
	}
}

func (c *Context) GetHeader(key string) string {
	if c.r == nil {
		return ""
	}
	return c.r.Header.Get(key)
}

func (c *Context) GetCookie(name string) string {
	if c.r == nil {
		return ""
	}
	ck, err := c.r.Cookie(name)
	if err != nil {
		return ""
	}
	return ck.Value
}

func (c *Context) RemoteAddr() string {
	if c.r == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(c.r.RemoteAddr)
	if err != nil {
		return c.r.RemoteAddr
	}
	return host
}

func (c *Context) ValueString(key string) string { return c.Values().GetString(key) }

func (c *Context) SetValue(key string, value any) { c.Values().Set(key, value) }

func (c *Context) Params() Params { return c.params }

func (c *Context) Values() *Values { return &c.values }

// Hooks returns the value attached by Router.UseHooks, usually *Settings.
func (c *Context) Hooks() any { return c.hooks }

func (c *Context) Status() int {
	if c.status == 0 {
		return http.StatusOK
	}
	return c.status
}

func (c *Context) StatusCode(code int) { c.status = code }

func (c *Context) writeHeader() {
	if c.wrote {
		return
	}
	c.wrote = true
	if c.status == 0 {
		c.status = http.StatusOK
	}
	c.w.WriteHeader(c.status)
}

// Header returns the outbound response headers (Iris-style).
// Validated request header reads use HeaderParam / HeaderString / HeaderWild.
func (c *Context) Header() http.Header { return c.w.Header() }

func (c *Context) Write(b []byte) (int, error) {
	c.writeHeader()
	if c.r != nil && c.r.Method == http.MethodHead {
		return 0, nil
	}
	return c.w.Write(b)
}

func (c *Context) WriteString(s string) (int, error) {
	c.writeHeader()
	if c.r != nil && c.r.Method == http.MethodHead {
		return 0, nil
	}
	return io.WriteString(c.w, s)
}

func (c *Context) JSON(v any) error {
	c.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.writeHeader()
	return json.MarshalWrite(c.w, v, jsontext.EscapeForHTML(false))
}

func (c *Context) JSONP(v any) error {
	callback := "callback"
	if c.r != nil {
		if q := c.r.URL.Query().Get("callback"); q != "" && validCallback(q) {
			callback = q
		}
	}
	c.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	c.writeHeader()
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err = io.WriteString(c.w, callback); err != nil {
		return err
	}
	if _, err = c.w.Write([]byte{'('}); err != nil {
		return err
	}
	if _, err = c.w.Write(body); err != nil {
		return err
	}
	_, err = c.w.Write([]byte{')', ';'})
	return err
}

func validCallback(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		r := s[i]
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '.', r == '$':
		default:
			return false
		}
	}
	return true
}

func (c *Context) ViewData(name string, v any) {
	if c.viewData == nil {
		c.viewData = make(map[string]any)
	}
	c.viewData[name] = v
}

func (c *Context) View(name string) error {
	if c.views == nil {
		return fmt.Errorf("view %s: no templates registered", name)
	}
	c.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.writeHeader()
	return c.views.ExecuteTemplate(c.w, name, c.viewData)
}

// Next runs the next handler in the chain (iris Context.Next).
// Handlers must call Next to continue; Done middleware only runs if the chain reaches them.
func (c *Context) Next() {
	if c.stopped {
		return
	}
	c.index++
	if c.index >= 0 && c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

// StopExecution skips the rest of the handler chain (iris Context.StopExecution).
func (c *Context) StopExecution() { c.stopped = true }

// IsStopped reports whether StopExecution was called.
func (c *Context) IsStopped() bool { return c.stopped }

// Proceed runs h and reports whether the chain may continue (iris Context.Proceed).
func (c *Context) Proceed(h Handler) bool {
	if c.stopped || h == nil {
		return false
	}
	h(c)
	return !c.stopped
}

type statusWriter struct {
	ctx *Context
}

func (s *statusWriter) Header() http.Header { return s.ctx.w.Header() }

func (s *statusWriter) WriteHeader(code int) {
	s.ctx.status = code
	s.ctx.writeHeader()
}

func (s *statusWriter) Write(b []byte) (int, error) {
	s.ctx.writeHeader()
	if len(b) == 0 {
		return 0, nil
	}
	return s.ctx.w.Write(b)
}

func (s *statusWriter) Unwrap() http.ResponseWriter { return s.ctx.w }

// FlushStatus sends the stored status when the body is empty.
func (c *Context) FlushStatus() { c.writeHeader() }

func pathSegments(p string) []string {
	for len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	for len(p) > 0 && p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}
	if p == "" {
		return nil
	}
	n := 1
	for i := 0; i < len(p); i++ {
		if p[i] == '/' {
			n++
		}
	}
	out := make([]string, 0, n)
	start := 0
	for i := 0; i <= len(p); i++ {
		if i == len(p) || p[i] == '/' {
			out = append(out, p[start:i])
			start = i + 1
		}
	}
	return out
}
