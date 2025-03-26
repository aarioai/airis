package response

import (
	"bytes"
	"encoding/json"
)

// EncodeJson json.Marshal 默认 escapeHtml 为true,会转义 <、>、&
func EncodeJson(v any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	je := json.NewEncoder(buf)
	je.SetEscapeHTML(false) // json Marshal 不转译 HTML 字符
	if err := je.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
