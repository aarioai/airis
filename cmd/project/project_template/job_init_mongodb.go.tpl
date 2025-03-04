package job

import (
	"context"
	"github.com/aarioai/airis-driver/driver/index"
    "github.com/aarioai/airis-driver/driver/mongodbhelper"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/helpers/debug"
	"{{APP_BASE}}/conf"
)

var (
	mongoIndexEntities []index.Entity
)

func initMongodb(ctx context.Context, app *core.App, profile *debug.Profile) {
	profile.Mark("init mongodb: %s", conf.MongoDBConfigSection)
	mongoObj := mongodbhelper.NewDB(app, conf.MongoDBConfigSection)
	// create table and indexes
	if len(mongoIndexEntities) > 0 {
		for _, t := range mongoIndexEntities {
			ae.PanicOn(mongoObj.CreateIndexes(ctx, t))
		}
	}
}
