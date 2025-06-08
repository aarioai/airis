package grpc

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"net"
)

func (s *Service) Run(ctx acontext.Context, profile *debug.Profile) {
	profile.Fork("running grpc server: account")
	port, err := s.app.Config.MustGetString("svc_account.grpc_port")
	ae.PanicOnErrs(err)
	listener, err := net.Listen("tcp", ":"+port)
	ae.PanicOnErrs(err)
	server := newServer()

	go func() {
		<-ctx.Done()
		alog.Stop("account")
		server.Stop()
	}()

	ae.PanicOnErrs(server.Serve(listener))
}
