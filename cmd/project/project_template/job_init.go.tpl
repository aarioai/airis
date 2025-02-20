package job

import (
	"context"
	"github.com/aarioai/airis/core"
	"github.com/aarioai/airis/core/helpers/debug"
)

func Init(app *core.App, profile *debug.Profile) {
	ctx := context.Background()
	profile.Mark("init {{APP_NAME}}")

	initMongodb(ctx, app, profile)
}
