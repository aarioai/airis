package grpc

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"net"
)

func (s *Service) Run(prof *debug.Profile) {
	profile := prof.Fork("starting grpc server ({{APP_NAME}})")
	port, err := s.app.Config.MustGetString("svc_{{APP_NAME}}.grpc_port")
	ae.PanicOnErrs(err)
	listener, err := net.Listen("tcp", ":"+port)
	ae.PanicOnErrs(err)
	server := newServer()

	go func() {
		<-s.app.GlobalContext.Done()
		alog.Stopf("grpc server ({{APP_NAME}}:%s)", port)
		server.Stop()
	}()
	profile.Mark("grpc server ({{APP_NAME}}) listening: " + port)
	ae.PanicOnErrs(server.Serve(listener))
}
