package httpsvr

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
)

func (c *Context) ensureSettings() {
	if c.settings == nil {
		c.settings = SettingsFrom(c)
	}
}

func (c *Context) WithHeader(key, value string) *Context {
	if key == "Last-Modified" && value == "Thu, 01 Jan 1970 00:00:00 GMT" {
		return c
	}
	if key == "Content-Length" {
		return c
	}
	if value == "" {
		return c.DeleteOutHeader(key)
	}
	if c.outHeaders == nil {
		c.outHeaders = make(map[string]string)
	}
	c.outHeaders[key] = value
	return c
}

func (c *Context) WithHeaders(headers map[string]string) *Context {
	for k, v := range headers {
		c.WithHeader(k, v)
	}
	return c
}

func (c *Context) WithServeContentTypes(contentTypes []string) *Context {
	if len(contentTypes) == 0 {
		ae.Panic("must register at least one content type")
	}
	c.serveContentTypes = contentTypes
	return c
}

func (c *Context) WithContentType(contentType string) *Context {
	return c.WithHeader("Content-Type", contentType)
}

func (c *Context) OutHeader(key string) string {
	if c.outHeaders == nil {
		return ""
	}
	return c.outHeaders[key]
}

// ResolveContentType negotiates the response Content-Type.
func (c *Context) ResolveContentType() string {
	ct := c.OutHeader("Content-Type")
	if ct != "" {
		return ct
	}
	ct = c.Values().GetString(acontext.CtxContentType)
	if ct != "" {
		c.WithHeader("Content-Type", ct)
		return ct
	}
	c.ensureSettings()
	serveTypes := c.serveContentTypes
	if serveTypes == nil && c.settings != nil {
		serveTypes = c.settings.ServeContentTypes
	}
	if len(serveTypes) == 0 {
		serveTypes = []string{"application/json"}
	}
	if ct = negotiateContentType(c.Accept(), serveTypes); ct == "" {
		ct = negotiateContentType(c.ReqContentType(), serveTypes)
	}
	if ct == "" && len(serveTypes) > 0 {
		ct = serveTypes[0]
	}
	c.WithHeader("Content-Type", ct)
	return ct
}

func negotiateContentType(header string, serveTypes []string) string {
	if header == "" || len(serveTypes) == 0 {
		return ""
	}
	if slices.Contains(serveTypes, header) {
		return header
	}
	best := ""
	bestQ := -1.0
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		media := part
		q := 1.0
		if i := strings.IndexByte(part, ';'); i >= 0 {
			media = strings.TrimSpace(part[:i])
			for _, param := range strings.Split(part[i+1:], ";") {
				param = strings.TrimSpace(param)
				if len(param) > 2 && (param[0] == 'q' || param[0] == 'Q') && param[1] == '=' {
					if v, err := strconv.ParseFloat(param[2:], 64); err == nil {
						q = v
					}
				}
			}
		}
		if q <= bestQ {
			continue
		}
		for _, st := range serveTypes {
			if mediaMatches(media, st) {
				best = st
				bestQ = q
				break
			}
		}
	}
	return best
}

func mediaMatches(media, serve string) bool {
	if media == serve || media == "*/*" {
		return true
	}
	if strings.HasSuffix(media, "/*") {
		return strings.HasPrefix(serve, media[:len(media)-1])
	}
	return false
}

func (c *Context) DeleteOutHeader(key string) *Context {
	if c.outHeaders != nil {
		delete(c.outHeaders, key)
	}
	return c
}

func (c *Context) WithErrorHandler(handler ErrorHandler) *Context {
	c.errorHandler = handler
	return c
}

func (c *Context) WithBeforeFlush(fn func(*Context)) *Context {
	c.beforeFlush = append(c.beforeFlush, fn)
	return c
}

func (c *Context) WithBeforeSerialize(beforeSerialize BeforeSerialize) *Context {
	c.beforeSerialize = append(c.beforeSerialize, beforeSerialize)
	return c
}

func (c *Context) WithSerialize(f func(contentType string, d Body) (bytes []byte, newContentType string, err error)) *Context {
	c.serialize = f
	return c
}

func (c *Context) writeDTOBytes(bytes []byte) (int, error) {
	writer := c.ResponseWriter()
	if c.outHeaders != nil {
		for k, v := range c.outHeaders {
			writer.Header().Set(k, v)
		}
	}
	if c.Method() == "HEAD" || len(bytes) == 0 {
		c.FlushStatus()
		return 0, nil
	}
	return writer.Write(bytes)
}

func (c *Context) WriteRaw(contentType string, bytes []byte) (int, error) {
	return c.WithContentType(contentType).writeDTOBytes(bytes)
}

func (c *Context) WriteRawHTML(bytes []byte) (int, error) {
	return c.WriteRaw(CtHTML.String(), bytes)
}

func (c *Context) WriteRawXML(bytes []byte) (int, error) {
	return c.WriteRaw(CtXML.String(), bytes)
}

func (c *Context) WriteRawJSON(bytes []byte) (int, error) {
	return c.WriteRaw(CtJSON.String(), bytes)
}

func (c *Context) WriteRawOctetStream(bytes []byte) (int, error) {
	return c.WriteRaw(CtOctetStream.String(), bytes)
}

func (c *Context) WriteJSONP(v any) error {
	data, e := c.decorateData(v)
	if e != nil {
		c.StatusCode(ae.InternalServerError)
		return errors.New("handle json data error: " + e.String())
	}
	return c.JSONP(data)
}

func (c *Context) writeDTO(d Body) (int, error) {
	c.ensureSettings()
	ct := c.ResolveContentType()
	if d.Code >= ae.BadRequest {
		if c.errorHandler != nil {
			n, err, next := c.errorHandler(c, ct, d)
			if !next {
				return n, err
			}
		}
		if c.settings != nil && c.settings.ErrorHandler != nil {
			n, err, next := c.settings.ErrorHandler(c, ct, d)
			if !next {
				return n, err
			}
		}
		n, err, next := defaultErrorHandler(c, ct, d)
		if !next {
			return n, err
		}
	}
	if c.settings != nil {
		for _, mw := range c.settings.BeforeSerialize {
			d = mw(c, ct, d)
		}
	}
	for _, mw := range c.beforeSerialize {
		d = mw(c, ct, d)
	}

	var (
		b              []byte
		newContentType string
		err            error
	)
	if c.serialize != nil {
		b, newContentType, err = c.serialize(ct, d)
	} else if c.settings != nil && c.settings.Serialize != nil {
		b, newContentType, err = c.settings.Serialize(ct, d)
	} else {
		b, newContentType, err = defaultSerialize(ct, d)
	}
	if err != nil {
		b = []byte(fmt.Sprintf(`{"code":500,"msg":"response serialize error: %s","data":null}`, err.Error()))
		return c.writeDTOBytes(b)
	}
	if newContentType != "" && newContentType != ct {
		c.WithContentType(newContentType)
	}
	return c.writeDTOBytes(b)
}

// WriteData writes a success DTO body (formerly Writer.Write).
func (c *Context) WriteData(a any) (int, error) {
	data, e := c.decorateData(a)
	if e != nil {
		return c.WriteE(e)
	}
	return c.writeDTO(Body{
		Code: ae.OK,
		Msg:  "OK",
		Data: data,
	})
}

func (c *Context) WriteOK() (int, error) {
	return c.writeDTO(Body{Code: ae.OK, Msg: "OK", Data: nil})
}

func (c *Context) WriteCreated(object []byte) (int, error) {
	return c.writeDTO(Body{Code: ae.Created, Msg: "Created", Data: object})
}

func (c *Context) WriteAccepted() (int, error) {
	return c.writeDTO(Body{Code: ae.Accepted, Msg: "Accepted", Data: nil})
}

func (c *Context) WriteCode(code int) (int, error) {
	return c.writeDTO(Body{Code: code, Msg: ae.Text(code)})
}

func (c *Context) WriteE(e *ae.Error) (int, error) {
	if e == nil {
		return c.WriteCode(ae.OK)
	}
	return c.writeDTO(Body{Code: e.Code, Msg: e.Msg})
}

func (c *Context) WriteErr(err error) (int, error) {
	if err == nil {
		return c.WriteCode(ae.OK)
	}
	return c.writeDTO(Body{Code: ae.InternalServerError, Msg: err.Error()})
}

func (c *Context) WriteMsg(code int, msg string) (int, error) {
	return c.writeDTO(Body{Code: code, Msg: msg})
}

func (c *Context) WriteMsgf(code int, format string, args ...any) (int, error) {
	return c.WriteMsg(code, fmt.Sprintf(format, args...))
}

func (c *Context) WriteId(id string) (int, error) {
	return c.WriteData(map[string]string{"id": id})
}

func (c *Context) TryWriteId(id string, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteId(id)
}

func (c *Context) WriteUint64Id(id uint64) (int, error) {
	return c.WriteData(map[string]uint64{"id": id})
}

func (c *Context) TryWriteUint64Id(id uint64, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteUint64Id(id)
}

func (c *Context) WriteUintId(id uint) (int, error) {
	return c.WriteData(map[string]uint{"id": id})
}

func (c *Context) TryWriteUintId(id uint, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteUintId(id)
}

func (c *Context) WriteAliasId(alias string, id string) (int, error) {
	return c.WriteData(map[string]string{alias: id})
}

func (c *Context) TryWriteAliasId(alias string, id string, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteAliasId(alias, id)
}

func (c *Context) WriteUint64AliasId(alias string, id uint64) (int, error) {
	return c.WriteData(map[string]uint64{alias: id})
}

func (c *Context) TryWriteUint64AliasId(alias string, id uint64, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteUint64AliasId(alias, id)
}

func (c *Context) WriteUintAliasId(alias string, id uint) (int, error) {
	return c.WriteData(map[string]uint{alias: id})
}

func (c *Context) TryWriteUintAliasId(alias string, id uint, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteUintAliasId(alias, id)
}

func (c *Context) WriteJointId(args ...any) (int, error) {
	l := len(args)
	if l < 2 || l%2 == 1 {
		c.StatusCode(ae.InternalServerError)
		return 0, fmt.Errorf("response no enough joint id args %+q", args)
	}
	id := make(map[string]any, l/2)
	for i := 0; i < l; i += 2 {
		id[args[i].(string)] = args[i+1]
	}
	return c.WriteData(id)
}

func (c *Context) TryWrite(a any, e *ae.Error) (int, error) {
	if e != nil {
		return c.WriteE(e)
	}
	return c.WriteData(a)
}
