package httpsvr

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/aarioai/airis/aa/atype"
)

type segment struct {
	static string
	name   string
	typ    string
}

type route struct {
	method   string
	segs     []segment
	paramN   int
	handlers []Handler
}

// Router is the Iris Application equivalent: root party + pre-router filters.
type Router struct {
	mu sync.Mutex

	root *Party

	// routerFilters run before route matching (iris UseRouter), e.g. CORS / recover.
	routerFilters []Handler

	byMethod map[string][]route

	errorHandlers map[int][]Handler

	hooks any
	views *template.Template
	srv   *http.Server
}

// Party groups routes under a shared prefix and middleware (iris Party).
type Party struct {
	router *Router
	prefix string
	begin  []Handler // iris Use
	done   []Handler // iris Done
}

// New is an alias of NewRouter (iris.New).
func New() *Router { return NewRouter() }

func NewRouter() *Router {
	rt := &Router{
		byMethod:      make(map[string][]route),
		errorHandlers: make(map[int][]Handler),
	}
	rt.root = &Party{router: rt, prefix: ""}
	return rt
}

// UseHooks attaches a value copied onto every request context.
func (rt *Router) UseHooks(hooks any) { rt.hooks = hooks }

// UseRouter registers handlers that run before route matching (iris Party.UseRouter).
// They always execute for every request, including 404/405.
func (rt *Router) UseRouter(handlers ...Handler) {
	rt.routerFilters = append(rt.routerFilters, handlers...)
}

// Use appends begin middleware on the root party (iris Application.Use).
func (rt *Router) Use(handlers ...Handler) { rt.root.Use(handlers...) }

// UseOnce inserts or replaces a begin handler by function identity (iris UseOnce).
func (rt *Router) UseOnce(handlers ...Handler) { rt.root.UseOnce(handlers...) }

// Done appends end middleware on the root party (iris Application.Done).
func (rt *Router) Done(handlers ...Handler) { rt.root.Done(handlers...) }

// Reset clears root begin/done middleware (iris Party.Reset).
func (rt *Router) Reset() *Party { return rt.root.Reset() }

// OnErrorCode registers handlers for an HTTP error status (iris OnErrorCode).
func (rt *Router) OnErrorCode(statusCode int, handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}
	rt.errorHandlers[statusCode] = append([]Handler{}, handlers...)
}

// OnAnyErrorCode registers handlers for common 4xx/5xx codes.
func (rt *Router) OnAnyErrorCode(handlers ...Handler) {
	for _, code := range []int{
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden,
		http.StatusNotFound, http.StatusMethodNotAllowed,
		http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable,
	} {
		rt.OnErrorCode(code, handlers...)
	}
}

// SetViews parses templates under dir.
func (rt *Router) SetViews(dir string) error {
	t, err := template.ParseGlob(strings.TrimRight(dir, `/\`) + "/*")
	if err != nil {
		return err
	}
	rt.views = t
	return nil
}

func (rt *Router) Party(relativePath string, handlers ...Handler) *Party {
	return rt.root.Party(relativePath, handlers...)
}

func (rt *Router) PartyFunc(relativePath string, fn func(*Party)) *Party {
	return rt.root.PartyFunc(relativePath, fn)
}

func (rt *Router) Get(path string, handlers ...Handler)    { rt.root.Get(path, handlers...) }
func (rt *Router) Post(path string, handlers ...Handler)   { rt.root.Post(path, handlers...) }
func (rt *Router) Put(path string, handlers ...Handler)    { rt.root.Put(path, handlers...) }
func (rt *Router) Patch(path string, handlers ...Handler)  { rt.root.Patch(path, handlers...) }
func (rt *Router) Head(path string, handlers ...Handler)   { rt.root.Head(path, handlers...) }
func (rt *Router) Delete(path string, handlers ...Handler) { rt.root.Delete(path, handlers...) }
func (rt *Router) Options(path string, handlers ...Handler) {
	rt.root.Options(path, handlers...)
}
func (rt *Router) Connect(path string, handlers ...Handler) {
	rt.root.Connect(path, handlers...)
}
func (rt *Router) Trace(path string, handlers ...Handler) { rt.root.Trace(path, handlers...) }
func (rt *Router) Any(path string, handlers ...Handler)   { rt.root.Any(path, handlers...) }
func (rt *Router) Handle(method, path string, handlers ...Handler) {
	rt.root.Handle(method, path, handlers...)
}

// GetRelPath returns "" for the root party.
func (rt *Router) GetRelPath() string { return rt.root.GetRelPath() }

func (p *Party) GetRelPath() string {
	if p.prefix == "" {
		return "/"
	}
	return p.prefix
}

func (p *Party) Party(relativePath string, handlers ...Handler) *Party {
	begin := append([]Handler{}, p.begin...)
	begin = append(begin, handlers...)
	done := append([]Handler{}, p.done...)
	return &Party{
		router: p.router,
		prefix: joinPath(p.prefix, relativePath),
		begin:  begin,
		done:   done,
	}
}

func (p *Party) PartyFunc(relativePath string, fn func(*Party)) *Party {
	child := p.Party(relativePath)
	if fn != nil {
		fn(child)
	}
	return child
}

// UseRouter registers pre-match filters on the parent router (iris Party.UseRouter).
func (p *Party) UseRouter(handlers ...Handler) {
	p.router.UseRouter(handlers...)
}

// OnErrorCode registers error handlers on the parent router.
func (p *Party) OnErrorCode(statusCode int, handlers ...Handler) {
	p.router.OnErrorCode(statusCode, handlers...)
}

// Use appends begin middleware for this party and future children (iris Party.Use).
func (p *Party) Use(handlers ...Handler) {
	p.begin = append(p.begin, handlers...)
}

// UseOnce inserts or replaces begin middleware by function pointer (iris UseOnce).
func (p *Party) UseOnce(handlers ...Handler) {
	for _, h := range handlers {
		if i := handlerIndex(p.begin, h); i >= 0 {
			p.begin[i] = h
			continue
		}
		p.begin = append(p.begin, h)
	}
}

// Done appends end middleware for this party and future children (iris Party.Done).
func (p *Party) Done(handlers ...Handler) {
	p.done = append(p.done, handlers...)
}

// Reset clears begin and done middleware inherited or registered on this party.
func (p *Party) Reset() *Party {
	p.begin = nil
	p.done = nil
	return p
}

func (p *Party) Get(path string, handlers ...Handler) {
	p.handle(http.MethodGet, path, handlers...)
}
func (p *Party) Post(path string, handlers ...Handler) {
	p.handle(http.MethodPost, path, handlers...)
}
func (p *Party) Put(path string, handlers ...Handler) {
	p.handle(http.MethodPut, path, handlers...)
}
func (p *Party) Patch(path string, handlers ...Handler) {
	p.handle(http.MethodPatch, path, handlers...)
}
func (p *Party) Head(path string, handlers ...Handler) {
	p.handle(http.MethodHead, path, handlers...)
}
func (p *Party) Delete(path string, handlers ...Handler) {
	p.handle(http.MethodDelete, path, handlers...)
}
func (p *Party) Options(path string, handlers ...Handler) {
	p.handle(http.MethodOptions, path, handlers...)
}
func (p *Party) Connect(path string, handlers ...Handler) {
	p.handle(http.MethodConnect, path, handlers...)
}
func (p *Party) Trace(path string, handlers ...Handler) {
	p.handle(http.MethodTrace, path, handlers...)
}

// Any registers the handlers for all standard HTTP methods (iris Party.Any).
func (p *Party) Any(path string, handlers ...Handler) {
	for _, m := range []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodConnect, http.MethodTrace,
	} {
		p.handle(m, path, handlers...)
	}
}

func (p *Party) Handle(method, path string, handlers ...Handler) {
	if method == "" {
		p.Any(path, handlers...)
		return
	}
	p.handle(method, path, handlers...)
}

func (p *Party) handle(method, path string, handlers ...Handler) {
	if len(handlers) == 0 {
		panic("httpsvr: route requires at least one handler")
	}
	full := joinPath(p.prefix, path)
	segs, err := parsePattern(full)
	if err != nil {
		panic(err)
	}
	paramN := 0
	for _, seg := range segs {
		if seg.name != "" {
			paramN++
		}
	}
	// iris order: begin (Use) → route handlers → done (Done)
	chain := make([]Handler, 0, len(p.begin)+len(handlers)+len(p.done))
	chain = append(chain, p.begin...)
	chain = append(chain, handlers...)
	chain = append(chain, p.done...)

	rt := p.router
	if rt.byMethod == nil {
		rt.byMethod = make(map[string][]route)
	}
	rt.byMethod[method] = append(rt.byMethod[method], route{
		method:   method,
		segs:     segs,
		paramN:   paramN,
		handlers: chain,
	})
}

func handlerIndex(list []Handler, h Handler) int {
	want := reflect.ValueOf(h).Pointer()
	for i, item := range list {
		if reflect.ValueOf(item).Pointer() == want {
			return i
		}
	}
	return -1
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(rt.routerFilters) > 0 {
		chain := make([]Handler, 0, len(rt.routerFilters)+1)
		chain = append(chain, rt.routerFilters...)
		chain = append(chain, rt.serveMatched)
		ctx := acquireContext(w, r, Params{}, chain, rt.hooks, rt.views)
		ctx.Next()
		if !ctx.wrote && ctx.status != 0 {
			ctx.writeHeader()
		}
		releaseContext(ctx)
		return
	}
	rt.dispatch(w, r)
}

// serveMatched is the terminal UseRouter handler that performs route matching.
func (rt *Router) serveMatched(c *Context) {
	rt.dispatchContext(c)
}

func (rt *Router) dispatch(w http.ResponseWriter, r *http.Request) {
	ctx := acquireContext(w, r, Params{}, nil, rt.hooks, rt.views)
	rt.dispatchContext(ctx)
	if !ctx.wrote && ctx.status != 0 {
		ctx.writeHeader()
	}
	releaseContext(ctx)
}

func (rt *Router) dispatchContext(c *Context) {
	r := c.Request()
	segs := pathSegments(r.URL.Path)
	if routes := rt.byMethod[r.Method]; len(routes) > 0 {
		for i := range routes {
			rtRoute := &routes[i]
			params, ok := matchPathSegs(rtRoute.segs, rtRoute.paramN, segs)
			if !ok {
				continue
			}
			c.params = params
			c.handlers = rtRoute.handlers
			c.index = -1
			c.stopped = false
			c.Next()
			return
		}
	}

	var allowed []string
	for method, routes := range rt.byMethod {
		if method == r.Method {
			continue
		}
		for i := range routes {
			if _, ok := matchPathSegs(routes[i].segs, routes[i].paramN, segs); ok {
				allowed = append(allowed, method)
				break
			}
		}
	}
	if len(allowed) > 0 {
		c.Header().Set("Allow", strings.Join(allowed, ", "))
		rt.writeError(c, http.StatusMethodNotAllowed)
		return
	}
	rt.writeError(c, http.StatusNotFound)
}

func (rt *Router) writeError(c *Context, status int) {
	c.StatusCode(status)
	if hs := rt.errorHandlers[status]; len(hs) > 0 {
		c.handlers = hs
		c.index = -1
		c.stopped = false
		c.Next()
		return
	}
	http.Error(c.ResponseWriter(), http.StatusText(status), status)
}

// Run / Listen start the HTTP server (iris Run / Listen).
func (rt *Router) Run(addr string) error { return rt.Listen(addr) }

func (rt *Router) Listen(addr string) error {
	srv := &http.Server{Addr: addr, Handler: rt}
	rt.mu.Lock()
	rt.srv = srv
	rt.mu.Unlock()
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (rt *Router) Shutdown(ctx context.Context) error {
	rt.mu.Lock()
	srv := rt.srv
	rt.mu.Unlock()
	if srv == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return srv.Shutdown(ctx)
}

func joinPath(prefix, path string) string {
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if prefix == "" {
		return path
	}
	if path == "/" {
		return prefix
	}
	return prefix + path
}

func parsePattern(path string) ([]segment, error) {
	parts := pathSegments(path)
	out := make([]segment, len(parts))
	for i, part := range parts {
		if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
			out[i] = segment{static: part}
			continue
		}
		body := part[1 : len(part)-1]
		name, typ, _ := strings.Cut(body, ":")
		if name == "" || strings.Contains(name, "}") {
			return nil, fmt.Errorf("httpsvr: invalid path parameter %q", part)
		}
		if typ == "" {
			typ = "string"
		}
		if !knownParamType(typ) {
			return nil, fmt.Errorf("httpsvr: unknown path parameter type %q", typ)
		}
		out[i] = segment{name: name, typ: typ}
	}
	return out, nil
}

func knownParamType(typ string) bool {
	switch typ {
	case "string", "bool", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"alphabetical", "uuid", "email", "mail", "weekday":
		return true
	default:
		return false
	}
}

func matchPathSegs(pattern []segment, paramN int, got []string) (Params, bool) {
	if len(pattern) != len(got) {
		return Params{}, false
	}
	var store []Param
	for i, seg := range pattern {
		if seg.name == "" {
			if seg.static != got[i] {
				return Params{}, false
			}
			continue
		}
		value, ok := parseParam(seg.typ, got[i])
		if !ok {
			return Params{}, false
		}
		if store == nil {
			store = make([]Param, 0, paramN)
		}
		store = append(store, Param{Key: seg.name, ValueRaw: value})
	}
	return Params{Store: store}, true
}

func parseParam(typ, raw string) (any, bool) {
	switch typ {
	case "", "string":
		return raw, raw != ""
	case "bool":
		switch raw {
		case "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes":
			return true, true
		case "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No":
			return false, true
		default:
			return nil, false
		}
	case "int":
		n, err := strconv.ParseInt(raw, 10, strconv.IntSize)
		if err != nil {
			return nil, false
		}
		return int(n), true
	case "int8":
		n, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			return nil, false
		}
		return int8(n), true
	case "int16":
		n, err := strconv.ParseInt(raw, 10, 16)
		if err != nil {
			return nil, false
		}
		return int16(n), true
	case "int32":
		n, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return nil, false
		}
		return int32(n), true
	case "int64":
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, false
		}
		return n, true
	case "uint":
		n, err := strconv.ParseUint(raw, 10, strconv.IntSize)
		if err != nil {
			return nil, false
		}
		return uint(n), true
	case "uint8":
		n, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			return nil, false
		}
		return uint8(n), true
	case "uint16":
		n, err := strconv.ParseUint(raw, 10, 16)
		if err != nil {
			return nil, false
		}
		return uint16(n), true
	case "uint32":
		n, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return nil, false
		}
		return uint32(n), true
	case "uint64":
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, false
		}
		return n, true
	case "alphabetical":
		if atype.IsAlphabetical(raw) {
			return raw, true
		}
		return nil, false
	case "uuid":
		if atype.IsUUID(raw) {
			return raw, true
		}
		return nil, false
	case "email", "mail":
		if atype.IsEmail(raw) {
			return raw, true
		}
		return nil, false
	case "weekday":
		w := atype.ToWeekday(raw)
		if w.Valid() {
			return w, true
		}
		return nil, false
	default:
		return nil, false
	}
}
