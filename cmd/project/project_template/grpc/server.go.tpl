package grpc

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"net"
)

func (s *Service) Run(profile *debug.Profile) {
	profile.Fork("running grpc service: {{APP_NAME}}")
	port, err := s.app.Config.MustGetString("svc_{{APP_NAME}}.grpc_port")
	ae.PanicOnErrs(err)
	listener, err := net.Listen("tcp", ":"+port)
	ae.PanicOnErrs(err)
	server := newServer()

	go func() {
		<-s.app.GlobalContext.Done()
		alog.Stop("{{APP_NAME}}")
		server.Stop()
	}()

	ae.PanicOnErrs(server.Serve(listener))
}
