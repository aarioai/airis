package dto

import "github.com/aarioai/airis/aa/atype"

// UserToken
// @doc https://www.rfc-editor.org/rfc/rfc6749#section-4.2.2
// standard: access_token, expires_in, scope, state, token_type
type UserToken struct {
	TokenType   string       `json:"token_type"` // Bearer  --> 客户端上传header: Authorization: Bearer $access_token
	AccessToken string       `json:"access_token"`
	ExpiresIn   atype.Second `json:"expires_in"`

	RefreshToken string `json:"refresh_token"` // optional
	State        string `json:"state"`         // 透传回客户端

	Attach map[string]any `json:"attach"`
	Scope  map[string]any `json:"scope"`
}
