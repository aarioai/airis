package afmt

import (
	"html/template"
	"net/url"
	"sort"
	"strings"
)

func Url(baseUrl string, params map[string]string) template.URL {
	if params == nil {
		return template.URL(strings.TrimSuffix(baseUrl, "?"))
	}
	var s strings.Builder
	s.WriteString(baseUrl)
	n := strings.IndexByte(baseUrl, '?')
	if n > 0 {
		g := baseUrl[len(baseUrl)-1]
		if g != '&' && g != '?' {
			s.WriteByte('&')
		}
	} else {
		s.WriteByte('?')
	}
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" && v != "0" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for i, k := range keys {
		s.WriteString(k)
		s.WriteByte('=')
		s.WriteString(url.QueryEscape(params[k]))
		if i < len(keys)-1 {
			s.WriteByte('&')
		}
	}
	return template.URL(s.String())
}
