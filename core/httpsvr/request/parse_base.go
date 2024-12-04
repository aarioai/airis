package request

import (
	"encoding/json"
	"github.com/aarioai/airis/core/ae"
	"github.com/aarioai/airis/core/atype"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
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

/*
@param pattern  e.g. `[[:word:]]+` `\w+`
Filter(pattern string, required bool)
Filter(required bool)
Filter(pattern string)
Filter(default atype.Atype)
*/
func (v *RawValue) Filter(patterns ...any) *ae.Error {
	required := true
	pattern := ""

	for i := 0; i < len(patterns); i++ {
		pat := patterns[i]
		if s, ok := pat.(string); ok {
			pattern = s
		} else if b, ok := pat.(bool); ok {
			required = b
		} else if d, ok := pat.(*atype.Atype); ok && v.String() == "" {
			v.Reload(d.Raw())
		}
	}
	if v.String() == "" {
		if required {
			return ae.BadParam(v.name)
		}
	} else if pattern != "" {
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

// ContentType
// Request是每个请求独立内存，基本不存在并发情况，没有加锁的必要性
func (r *Request) ContentType() string {
	t := r.contentType
	if t != "" {
		return t
	}
	r.contentType, _, _ = r.contentMediaType()
	return t
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

func (r *Request) Query(param string, patterns ...any) (*RawValue, *ae.Error) {
	var userData url.Values
	if r.r != nil {
		userData = r.r.URL.Query()
	}
	return r.query(r.partialQueries, userData, param, patterns...)
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

// @see http.parsePostForm()
func (r *Request) parseBodyStream() *ae.Error {
	defer func() {
		r.bodyParsed = true
	}()
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
		var reader io.Reader = r.r.Body
		maxFormSize := int64(1<<63 - 1)
		if _, ok := reader.(*maxBytesReader); !ok {
			maxFormSize = int64(10 << 20) // 10 MB is a lot of json.
			reader = io.LimitReader(r.r.Body, maxFormSize+1)
		}
		b, err := io.ReadAll(reader)
		if err != nil {
			return ae.NewError(err)
		}
		if int64(len(b)) > maxFormSize {
			return ae.New(ae.CodeRequestEntityTooLarge, "body is too large")
		}
		if ct == CtForm.String() {
			var vs url.Values
			if vs, err = url.ParseQuery(string(b)); err != nil {
				return ae.New(ae.CodeBadRequest, "body should encode in application/x-www-form-urlencoded ")
			}
			r.setPartialBodyData(vs)
		} else {
			if err = json.Unmarshal(b, &r.partialBodyData); err != nil {
				return ae.New(ae.CodeBadRequest, "body should encode in json")
			}
		}
		// @see http.ParseMultipartForm
		// .MultipartForm "multipart/form-data" ||  "multipart/mixed"
	// .PostFormValue  "application/x-www-form-urlencoded" url.ParseQuery(body) + .MultipartForm
	// .FormValue() 调用 .Form  =  url.ParseQuery(r.URL.RawQuery) + .PostFormValue

	case CtFormData.String(), CtMixed.String():
		boundary, ok := params["boundary"]
		if !ok {
			return ae.New(ae.CodeBadRequest, "no multipart boundary param in Content-Type")
		}
		var f *multipart.Form
		f, err = multipart.NewReader(r.r.Body, boundary).ReadForm(32 << 20) // 32M
		if err != nil {
			return ae.New(ae.CodeBadRequest, "body should encode in multipart form")
		}
		r.setPartialBodyData(f.Value)
	}
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
func (r *Request) GetHeader(param string, patterns ...any) (*RawValue, *ae.Error) {
	var userData http.Header
	if r.r != nil {
		userData = r.r.Header
	}
	return r.query(r.partialHeaders, userData, param, patterns...)
}
func (r *Request) Header(param string) string {
	value, e := r.GetHeader(param)
	if e != nil {
		return ""
	}
	return value.String()
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

func (r *Request) UserAgent() string {
	ua := r.userAgent
	if ua == "" {
		ua = r.Header("User-Agent")
		r.userAgent = ua
	}
	return ua
}
