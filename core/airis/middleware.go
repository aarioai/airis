package airis

import (
	"github.com/kataras/iris/v12"
)

func WithRemoteUser(ictx iris.Context, user string) iris.Context {
	ictx.Values().Set(CtxRemoteUser, user)
	return ictx
}
