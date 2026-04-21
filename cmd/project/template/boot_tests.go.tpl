package boot

import "github.com/aarioai/airis/aa"

func SelfTest(app *aa.App) bool {
	acontext.ServFallbackReady()
	return true
}


// @TODO check mysql/posgresql/redis ... here, and then set to acontext ServiceHealth