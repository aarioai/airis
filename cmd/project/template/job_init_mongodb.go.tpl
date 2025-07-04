package job

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/helpers/debug"
	"{{APP_BASE}}/conf"
)

var (
	mongoIndexEntities []index.Entity
)

func (s *Service) initMongodb(ctx acontext.Context, profile *debug.Profile) {
	profile.Markf("init mongodb: %s", conf.MongoCfgSection)
	mongoObj := mongodb.NewDB(s.app, conf.MongoCfgSection)
	// create table and indexes
	if len(mongoIndexEntities) > 0 {
		for _, t := range mongoIndexEntities {
			ae.PanicOn(mongoObj.ORM(t).CreateIndexes(ctx))
		}
	}
}
