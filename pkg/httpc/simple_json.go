package httpc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

var (
	JsonUnmarshal = json.Unmarshal
)

func requestJsonBody(req *http.Request, target interface{}) (int, error) {
	code, data, err := requestBody(req)
	if err != nil {
		return code, err
	}
	return code, JsonUnmarshal(data, &target)
}

func addJsonHeaders(req *http.Request, headers []map[string]string) {
	SetHeaders(req, headers)
	SetNxHeader(req, "Content-Type", "application/json")
	SetNxHeader(req, "Accept", "application/json")
}

// RequestJson
// Example:
// - RequestJson("HEAD", url, nil)
// - RequestJson("TRACE", url, nil)
// - RequestJson("CONNECT", url, nil)
func RequestJson(ctx context.Context, method, url string, params map[string]string, headers ...map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	req.Close = true
	addJsonHeaders(req, headers)
	SetQueries(req, params)
	return request(req)
}

func RequestJsonBody(ctx context.Context, target interface{}, method, url string, params map[string]string, data io.Reader, headers ...map[string]string) (int, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return 0, err
	}
	req.WithContext(ctx)
	req.Close = true
	addJsonHeaders(req, headers)
	SetQueries(req, params)
	return requestJsonBody(req, target)
}

func DeleteJson(ctx context.Context, target interface{}, url string, params map[string]string, headers ...map[string]string) (int, error) {
	return RequestJsonBody(ctx, target, "DELETE", url, params, nil, headers...)
}

func GetJson(ctx context.Context, target interface{}, url string, params map[string]string, headers ...map[string]string) (int, error) {
	return RequestJsonBody(ctx, target, "GET", url, params, nil, headers...)
}

func PatchJson(ctx context.Context, target interface{}, url string, data io.Reader, headers ...map[string]string) (int, error) {
	return RequestJsonBody(ctx, target, "PATCH", url, nil, data, headers...)
}

func PostJson(ctx context.Context, target interface{}, url string, data io.Reader, headers ...map[string]string) (int, error) {
	return RequestJsonBody(ctx, target, "POST", url, nil, data, headers...)
}

func PutJson(ctx context.Context, target interface{}, url string, data io.Reader, headers ...map[string]string) (int, error) {
	return RequestJsonBody(ctx, target, "PUT", url, nil, data, headers...)
}
