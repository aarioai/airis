package request

import "mime"

type ContentTypes string

const (
	CtJSON ContentTypes = "application/json"
	CtHTML ContentTypes = "text/html"
	// CtXML RFC 7303 最新的，已经将 application/xml 优先到9.1，优于 text/xml
	// @see https://www.rfc-editor.org/rfc/rfc7303#section-9.1
	CtXML         ContentTypes = "application/xml"
	CtJSONP       ContentTypes = "text/html"
	CtOctetStream ContentTypes = "application/octet-stream"
	CtForm        ContentTypes = "application/x-www-form-urlencoded"
	CtFormData    ContentTypes = "multipart/form-data"
	CtMixed       ContentTypes = "multipart/mixed"
)

const (
	ParamStringify = "_stringify"
	ParamField     = "_field"
	ParamPage      = "page"
	ParamPageSize  = "page_size"
	ParamPageEnd   = "page_end"
)

func (t ContentTypes) String() string {
	return string(t)
}
func (t ContentTypes) WithCharset(charset string) string {
	return t.String() + "; charset=" + charset
}
func (t ContentTypes) WithUtf8() string {
	return t.WithCharset("utf-8")
}
func IsHtml(contentType string) bool {
	ct, _, err := mime.ParseMediaType(contentType)
	return err == nil && ct == CtHTML.String()
}
