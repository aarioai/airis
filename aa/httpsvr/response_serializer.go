package httpsvr

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
)

// EncodeJson marshals v as JSON without HTML escaping and without a trailing newline.
func EncodeJson(v any) ([]byte, error) {
	return json.Marshal(v, jsontext.EscapeForHTML(false))
}
