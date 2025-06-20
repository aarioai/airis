package service

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"project/simple/app/app_aajs/module/bs/dto"
)

func (s *Service) Login(ctx context.Context, account, password, state string) (dto.UserToken, *ae.Error) {
	if account != "12345" || password != "hello" {
		return dto.UserToken{}, ae.ErrorNotAcceptable
	}
	token := dto.UserToken{
		TokenType:    "Bearer",
		AccessToken:  "helloworld",
		ExpiresIn:    2 * atype.HourInSecond,
		RefreshToken: "refresh_helloworld",
		State:        state,
		Attach: map[string]any{
			"refresh_api":  "PUT /v1/auth/access_token",
			"refresh_ttl":  365 * atype.DayInSecond,
			"secure":       true,
			"validate_api": "HEAD /v1/auth/access_token",
		},
		Scope: nil,
	}
	return token, nil
}

func (s *Service) RefreshUserToken(ctx context.Context, refreshToken string) (dto.UserToken, *ae.Error) {
	if refreshToken != "refresh_helloworld" {
		return dto.UserToken{}, ae.ErrorNotAcceptable
	}
	token := dto.UserToken{
		TokenType:    "Bearer",
		AccessToken:  "helloworld",
		ExpiresIn:    2 * atype.HourInSecond,
		RefreshToken: "refresh_helloworld",
		State:        "",
		Attach: map[string]any{
			"refresh_api":  "PUT /v1/auth/access_token",
			"refresh_ttl":  365 * atype.DayInSecond,
			"secure":       true,
			"validate_api": "HEAD /v1/auth/access_token",
		},
		Scope: nil,
	}
	return token, nil
}
