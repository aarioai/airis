package httpc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	JsonUnmarshal = json.Unmarshal
)

func requestJsonBody(req *http.Request, target interface{}) (int, error) {
	code, body, err := requestBody(req)
	if err != nil {
		return code, err
	}
	return code, JsonUnmarshal(body, &target)
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

func RequestJsonBody(target interface{}, method, url string, params map[string]string, body []byte, headers ...map[string]string) (int, error) {
	var b *bytes.Reader
	if len(body) > 0 {
		b = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return 0, err
	}
	req.Close = true

	addJsonHeaders(req, headers)
	SetQueries(req, params)
	return requestJsonBody(req, target)
}

func GetJson(target interface{}, url string, params map[string]string, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "GET", url, params, nil, headers...)
}

func PostJson(target interface{}, url string, body []byte, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "POST", url, nil, body, headers...)
}

func DeleteJson(target interface{}, url string, params map[string]string, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "DELETE", url, params, nil, headers...)
}

func PutJson(target interface{}, url string, body []byte, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "PUT", url, nil, body, headers...)
}

func PatchJson(target interface{}, url string, body []byte, headers ...map[string]string) (int, error) {
	return RequestJsonBody(target, "PATCH", url, nil, body, headers...)
}
