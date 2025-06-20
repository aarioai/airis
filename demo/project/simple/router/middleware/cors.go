package middleware

import (
	"github.com/aarioai/golib/sdk/irisz/middlewarez"
	"github.com/iris-contrib/middleware/cors"
)

var CORS = cors.New(cors.Options{
	AllowOriginFunc:    middlewarez.AllowDomainsFunc("luexu.com"),
	AllowedMethods:     []string{"*"},
	AllowedHeaders:     []string{"*"},
	ExposedHeaders:     []string{"*"},
	MaxAge:             300,
	AllowCredentials:   true,
	OptionsPassthrough: false,
	Debug:              false,
})
