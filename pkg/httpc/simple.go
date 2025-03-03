package httpc

import (
	"io"
	"net/http"
)

func Request(method, url string, params map[string]string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	req, err := http.NewRequest(method, url, data)
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

func Patch(url string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	return Request("PATCH", url, nil, data, headers...)
}

func Post(url string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	return Request("POST", url, nil, data, headers...)
}

func Put(url string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	return Request("PUT", url, nil, data, headers...)
}
