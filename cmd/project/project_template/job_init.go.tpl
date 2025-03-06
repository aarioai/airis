package job

import (
	"context"
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/helpers/debug"
)

func (s *Service) Init(ctx acontext.Context, profile *debug.Profile) {
	profile.Mark("job init app:{{APP_NAME}}")

	s.initMongodb(ctx, profile)
}
