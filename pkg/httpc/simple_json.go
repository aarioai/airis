package httpc

import (
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
func RequestJson(method, url string, params map[string]string, headers ...map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	addJsonHeaders(req, headers)
	SetQueries(req, params)
	return request(req)
}

func RequestJsonBody(target interface{}, method, url string, params map[string]string, data io.Reader, headers ...map[string]string) (int, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return 0, err
	}
	req.Close = true

	addJsonHeaders(req, headers)
	SetQueries(req, params)
	return requestJsonBody(req, target)
}

func DeleteJson(target interface{}, url string, params map[string]string, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "DELETE", url, params, nil, headers...)
}

func GetJson(target interface{}, url string, params map[string]string, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "GET", url, params, nil, headers...)
}

func PatchJson(target interface{}, url string, data io.Reader, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "PATCH", url, nil, data, headers...)
}

func PostJson(target interface{}, url string, data io.Reader, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "POST", url, nil, data, headers...)
}

func PutJson(target interface{}, url string, data io.Reader, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "PUT", url, nil, data, headers...)
}
