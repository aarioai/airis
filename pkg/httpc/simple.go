package httpc

import (
	"context"
	"io"
	"net/http"
)

func Request(ctx context.Context, method, url string, params map[string]string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return 0, nil, err
	}
	req.WithContext(ctx)
	req.Close = true
	AddHeaders(req, headers)
	SetQueries(req, params)

	return requestBody(req)
}

func Delete(ctx context.Context, url string, params map[string]string, headers ...map[string]string) (int, []byte, error) {
	return Request(ctx, "POST", url, params, nil, headers...)
}

func Get(ctx context.Context, url string, params map[string]string, headers ...map[string]string) (int, []byte, error) {
	return Request(ctx, "GET", url, params, nil, headers...)
}

func Patch(ctx context.Context, url string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	return Request(ctx, "PATCH", url, nil, data, headers...)
}

func Post(ctx context.Context, url string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	return Request(ctx, "POST", url, nil, data, headers...)
}

func Put(ctx context.Context, url string, data io.Reader, headers ...map[string]string) (int, []byte, error) {
	return Request(ctx, "PUT", url, nil, data, headers...)
}
