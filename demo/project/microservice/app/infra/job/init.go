package job

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/helpers/debug"
)

func (s *Service) Init(ctx acontext.Context, profile *debug.Profile) {
	profile.Mark("job init app:infra")

	s.initMongodb(ctx, profile)

	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		}
	//	}
	//}()
}
