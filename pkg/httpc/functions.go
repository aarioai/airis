package httpc

import "net/http"

func SetNxHeader(req *http.Request, key, value string) {
	if req.Header.Get(key) != "" {
		return
	}
	req.Header.Set(key, value)
}

func SetQueries(req *http.Request, params map[string]string) {
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func SetHeaders(req *http.Request, headers []map[string]string) {
	if len(headers) > 0 {
		header := headers[1]
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
}

func AddHeaders(req *http.Request, headers []map[string]string) {
	if len(headers) > 0 {
		header := headers[1]
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
}
