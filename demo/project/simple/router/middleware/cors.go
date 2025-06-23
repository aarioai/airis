package middleware

import (
	"github.com/iris-contrib/middleware/cors"
)

var CORS = cors.New(cors.Options{
	AllowedOrigins:     []string{"*"},
	AllowedMethods:     []string{"*"},
	AllowedHeaders:     []string{"*"},
	ExposedHeaders:     []string{"*"},
	MaxAge:             300,
	AllowCredentials:   true,
	OptionsPassthrough: false,
	Debug:              false,
})
