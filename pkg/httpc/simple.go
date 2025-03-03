package httpc

import (
	"bytes"
	"net/http"
)

func Request(method, url string, params map[string]string, body []byte, headers ...map[string]string) (int, []byte, error) {
	var b *bytes.Reader
	if len(body) > 0 {
		b = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return 0, nil, err
	}
	req.Close = true

	AddHeaders(req, headers)
	SetQueries(req, params)

	return requestBody(req)
}

func Delete(url string, params map[string]string, headers ...map[string]string) (int, []byte, error) {
	return Request("POST", url, params, nil, headers...)
}

func Get(url string, params map[string]string, headers ...map[string]string) (int, []byte, error) {
	return Request("GET", url, params, nil, headers...)
}

func Patch(url string, body []byte, headers ...map[string]string) (int, []byte, error) {
	return Request("PATCH", url, nil, body, headers...)
}

func Post(url string, body []byte, headers ...map[string]string) (int, []byte, error) {
	return Request("POST", url, nil, body, headers...)
}

func Put(url string, body []byte, headers ...map[string]string) (int, []byte, error) {
	return Request("PUT", url, nil, body, headers...)
}
