package acontext

// WithRemoteUser stores the authenticated user on the request.
func WithRemoteUser[T interface{ SetValue(string, any) }](ctx T, user string) T {
	ctx.SetValue(CtxRemoteUser, user)
	return ctx
}
