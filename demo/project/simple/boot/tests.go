package boot

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/acontext"
)

func SelfTest(app *aa.App) {
	acontext.ServFallbackReady()
}

// @TODO check mysql/posgresql/redis ... here, and then set to acontext ServiceHealth
