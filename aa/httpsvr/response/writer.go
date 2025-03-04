package response

import (
	"errors"
	"fmt"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
	"slices"
	"sync"
)

// Writer
// @extend io.Closer
type Writer struct {
	SerializeTag      string
	serveContentTypes []string

	beforeFlush     []func(*Writer)
	beforeSerialize []func(ictx iris.Context, contentType string, d Body) Body
	serialize       func(contentType string, d Body) (bytes []byte, newContentType string, err error)
	errorHandler    func(ictx iris.Context, request *request.Request, contentType string, d Body) (int, error, bool)

	ictx    iris.Context
	request *request.Request

	code          int
	headers       map[string]string // 每个请求独立 Writer，不需要异步操作
	content       []byte
	contentStruct Body
}

var (
	// 对象池，减少内存分配
	// sync.Pool 通常不需要手动释放对象，当创建的对象，没有引用时会自动回收
	writerPool = sync.Pool{
		New: func() interface{} {
			return new(Writer)
		},
	}
)

func NewWriter(ictx iris.Context, request *request.Request) *Writer {
	var headers map[string]string
	w := Writer{
		SerializeTag:      "",
		serveContentTypes: nil,
		beforeFlush:       nil,
		beforeSerialize:   nil,
		serialize:         nil,
		errorHandler:      nil,
		ictx:              ictx,
		request:           request,
		code:              0,
		headers:           headers,
		content:           nil,
		contentStruct:     Body{},
	}
	return &w
}

func (w *Writer) WithHeader(key, value string) *Writer {
	if key == "Last-Modified" && value == "Thu, 01 Jan 1970 00:00:00 GMT" {
		return w
	}
	// 因为内容变更了，必须要把Content-Length设为空，不然客户端会读取错误
	// 这里设置Content-Length之后，iris Gzip 就会异常
	// Content-Length 可以在serialize里面设置
	if key == "Content-Length" {
		return w
	}
	if value == "" {
		return w.DeleteHeader(key)
	}
	if w.headers == nil {
		w.headers = make(map[string]string)
	}
	w.headers[key] = value
	return w
}
func (w *Writer) WithHeaders(headers map[string]string) *Writer {
	for k, v := range headers {
		w.WithHeader(k, v)
	}
	return w
}
func (w *Writer) WithServeContentTypes(contentTypes []string) *Writer {
	if len(contentTypes) == 0 {
		ae.Panic("must register at least one content type")
	}
	w.serveContentTypes = contentTypes
	return w
}
func (w *Writer) WithContentType(contentType string) *Writer {
	return w.WithHeader("Content-Type", contentType)
}
func (w *Writer) Header(key string) string {
	if w.headers == nil {
		return ""
	}
	value, ok := w.headers[key]
	if !ok {
		return ""
	}
	return value
}

func (w *Writer) ContentType() string {
	// ① 读取手动指定的
	ct := w.Header("Content-Type")
	if ct != "" {
		return ct
	}
	// ② 读取通过middleware统一对ictx设置的
	ct = w.ictx.Values().GetString(acontext.CtxContentType)
	if ct != "" {
		w.WithHeader("Content-Type", ct)
		return ct
	}
	// ③ 尝试解析 Accept 和 客户端ContentType
	serveTypes := w.serveContentTypes
	if serveTypes == nil {
		serveTypes = globalServeContentTypes
	}
	types := []string{w.request.Accept(), w.request.ContentType()}
	for _, t := range types {
		if t != "" {
			if ok := slices.Contains(serveTypes, t); ok {
				ct = t
				break
			}
		}
	}
	// ④ 使用注册的第一个ContentType。如果只注册了一个Content-Type，那么即表示只提供一种数据解析（如json）
	if ct != "" {
		ct = serveTypes[0]
	}
	w.WithHeader("Content-Type", ct)
	return ct
}
func (w *Writer) DeleteHeader(key string) *Writer {
	if w.headers != nil {
		delete(w.headers, key)
	}
	return w
}
func (w *Writer) WithErrorHandler(handler func(ictx iris.Context, request *request.Request, contentType string, d Body) (int, error, bool)) *Writer {
	w.errorHandler = handler
	return w
}
func (w *Writer) WithBeforeFlush(fn func(*Writer)) *Writer {
	if w.beforeFlush == nil {
		w.beforeFlush = make([]func(*Writer), 0)
	}
	w.beforeFlush = append(w.beforeFlush, fn)
	return w
}
func (w *Writer) WithBeforeSerialize(beforeSerialize func(ictx iris.Context, contentType string, d Body) Body) *Writer {
	if w.beforeSerialize == nil {
		w.beforeSerialize = make([]func(ictx iris.Context, contentType string, d Body) Body, 0)
	}
	w.beforeSerialize = append(w.beforeSerialize, beforeSerialize)
	return w
}
func (w *Writer) WithSerialize(f func(contentType string, d Body) (bytes []byte, newContentType string, err error)) *Writer {
	w.serialize = f
	return w
}
func (w *Writer) Context() iris.Context {
	return w.ictx
}
func (w *Writer) Request() *request.Request {
	return w.request
}
func (w *Writer) StatusCode(code int) {
	w.ictx.StatusCode(code)
}

func (w *Writer) write(bytes []byte) (int, error) {
	writer := w.ictx.ResponseWriter()
	if w.headers != nil {
		for k, v := range w.headers {
			writer.Header().Set(k, v)
		}
	}
	if w.request.Method() == "HEAD" || len(bytes) == 0 {
		return 0, nil
	}
	return writer.Write(bytes)
}
func (w *Writer) WriteRaw(contentType string, bytes []byte) (int, error) {
	return w.WithContentType(contentType).write(bytes)
}
func (w *Writer) WriteRawHTML(bytes []byte) (int, error) {
	return w.WriteRaw(request.CtHTML.String(), bytes)
}
func (w *Writer) WriteRawXML(bytes []byte) (int, error) {
	return w.WriteRaw(request.CtXML.String(), bytes)
}
func (w *Writer) WriteRawJSON(bytes []byte) (int, error) {
	return w.WriteRaw(request.CtJSON.String(), bytes)
}
func (w *Writer) WriteRawOctetStream(bytes []byte) (int, error) {
	return w.WriteRaw(request.CtOctetStream.String(), bytes)
}
func (w *Writer) WriteJSONP(v any, opts ...iris.JSONP) error {
	data, e := w.decorateData(v)
	if e != nil {
		w.StatusCode(ae.InternalServerError)
		return errors.New("handle json data error: " + e.String())
	}
	return w.ictx.JSONP(data, opts...)
}

func (w *Writer) writeDTO(d Body) (int, error) {
	ct := w.ContentType()
	if d.Code >= ae.BadRequest {
		// 避免重复调用，不再传 *Writer，而是直接操作 ictx
		if w.errorHandler != nil {
			n, err, next := w.errorHandler(w.ictx, w.request, ct, d)
			if !next {
				return n, err
			}
		}
		if globalErrorHandler != nil {
			n, err, next := globalErrorHandler(w.ictx, w.request, ct, d)
			if !next {
				return n, err
			}
		}
		n, err, next := defaultErrorHandler(w.ictx, w.request, ct, d)
		if !next {
			return n, err
		}
	}
	if len(globalBeforeSerialize) > 0 {
		for _, mw := range globalBeforeSerialize {
			d = mw(w.ictx, ct, d)
		}
	}
	if len(w.beforeSerialize) > 0 {
		for _, mw := range w.beforeSerialize {
			d = mw(w.ictx, ct, d)
		}
	}

	var (
		b              []byte
		newContentType string
		err            error
	)
	if w.serialize != nil {
		b, newContentType, err = w.serialize(ct, d)
	} else if globalSerialize != nil {
		b, newContentType, err = globalSerialize(ct, d)
	} else {
		b, newContentType, err = defaultSerialize(ct, d)
	}
	if err != nil {
		b = []byte(fmt.Sprintf(`{"code":500,"msg":"response serialize error: %s","data":null}`, err.Error()))
		return w.write(b)
	}
	if newContentType != "" && newContentType != ct {
		w.WithContentType(newContentType)
	}
	return w.write(b)
}

func (w *Writer) Write(a any) (int, error) {
	data, e := w.decorateData(a)
	if e != nil {
		return w.WriteE(e)
	}
	return w.writeDTO(Body{
		Code: ae.OK,
		Msg:  "OK",
		Data: data,
	})
}

func (w *Writer) WriteOK() (int, error) {
	return w.writeDTO(Body{
		Code: ae.OK,
		Msg:  "OK",
		Data: nil,
	})
}
func (w *Writer) WriteCode(code int) (int, error) {
	return w.writeDTO(Body{
		Code: code,
		Msg:  http.StatusText(code),
	})
}

func (w *Writer) WriteE(e *ae.Error) (int, error) {
	if e == nil {
		return w.WriteCode(ae.OK)
	}
	return w.writeDTO(Body{
		Code: e.Code,
		Msg:  e.Msg,
	})
}

func (w *Writer) WriteErr(err error) (int, error) {
	if err == nil {
		return w.WriteCode(ae.OK)
	}
	return w.writeDTO(Body{
		Code: ae.InternalServerError,
		Msg:  err.Error(),
	})
}

func (w *Writer) WriteError(code int, msg string, args ...any) (int, error) {
	return w.writeDTO(Body{
		Code: code,
		Msg:  afmt.Sprintf(msg, args...),
	})
}

// 返回插入数据的ID，ID 可能是联合主键，或者字段不为id，那么就会以对象形式返回
// 如： {"id":12314}   {"id":"ADREDD"}   {"id":{"k":"i_am_prinary_key"}}  {"id": {"k":"", "uid":""}}
func (w *Writer) WriteId(id string) (int, error) {
	return w.Write(map[string]string{"id": id})
}
func (w *Writer) TryWriteId(id string, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.WriteId(id)
}
func (w *Writer) WriteUint64Id(id uint64) (int, error) {
	return w.Write(map[string]uint64{"id": id})
}
func (w *Writer) TryWriteUint64Id(id uint64, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.WriteUint64Id(id)
}
func (w *Writer) WriteUintId(id uint) (int, error) {
	return w.Write(map[string]uint{"id": id})
}
func (w *Writer) TryWriteUintId(id uint, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.WriteUintId(id)
}
func (w *Writer) WriteAliasId(alias string, id string) (int, error) {
	return w.Write(map[string]string{alias: id})
}
func (w *Writer) TryWriteAliasId(alias string, id string, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.WriteAliasId(alias, id)
}
func (w *Writer) WriteUint64AliasId(alias string, id uint64) (int, error) {
	return w.Write(map[string]uint64{alias: id})
}
func (w *Writer) TryWriteUint64AliasId(alias string, id uint64, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.WriteUint64AliasId(alias, id)
}
func (w *Writer) WriteUintAliasId(alias string, id uint) (int, error) {
	return w.Write(map[string]uint{alias: id})
}
func (w *Writer) TryWriteUintAliasId(alias string, id uint, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.WriteUintAliasId(alias, id)
}

// k1,v1, k2, v2, k3,v3
func (w *Writer) WriteJointId(args ...any) (int, error) {
	l := len(args)
	if l < 2 || l%2 == 1 {
		w.StatusCode(ae.InternalServerError)
		return 0, fmt.Errorf("response no enough joint id args %+q", args)
	}
	id := make(map[string]any, l/2)
	for i := 0; i < l; i += 2 {
		id[args[i].(string)] = args[i+1]
	}
	return w.Write(id)
}

func (w *Writer) TryWrite(a any, e *ae.Error) (int, error) {
	if e != nil {
		return w.WriteE(e)
	}
	return w.Write(a)
}

// Close 释放实例到对象池
// 即使这个对象不是从对象池中获取的，也会放入对象池。不影响使用。
func (w *Writer) Close() error {
	writerPool.Put(w)
	return nil
}

func (w *Writer) CloseWith(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return w.Close()
}
