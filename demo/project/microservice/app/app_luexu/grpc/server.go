package grpc

import (
	"errors"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
    "github.com/aarioai/airis/pkg/types"
    "google.golang.org/grpc"
	"net"
)

func (s *Service) Serve(prof *debug.Profile) (*grpc.Server, string, error) {
	prof.Fork("starting gRPC server (app_luexu)")
	listener, serviceID, err := s.listen()
	if err != nil {
		return nil, "", err
	}
	server := s.registerServer()

	go func() {
		<-s.app.GlobalContext.Done()
		alog.Stopf("gRPC server (app_luexu)")
		s.app.Config.DeregisterGRPCService(serviceID)
		server.GracefulStop()
	}()

	go func() {
		defer s.app.Config.DeregisterGRPCService(serviceID)
		ae.PanicOnErrs(server.Serve(listener))
	}()

	prof.Forkf("register gRPC server (app_luexu: %s)", serviceID)

	return server, serviceID, nil
}

func (s *Service) listen() (net.Listener, string, error) {
	addr := s.app.Config.GetString("app_luexu.grpc_addr", "0.0.0.0")
	registerAddr := s.app.Config.GetString("app_luexu.grpc_register_addr", "127.0.0.1")
    checkAddr := s.app.Config.GetString("app_luexu.grpc_check_addr", registerAddr)

    port, _ := types.ParseInt(s.app.Config.GetString("app_luexu.grpc_port"))
    if port <= 0 {
        return nil, "", errors.New("missing or invalid config app_luexu.grpc_port")
    }
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
    if err != nil {
        return nil, "", err
    }

    serviceName := s.app.Config.GetString("app_luexu.grpc_service_name", "app_luexu")
    serviceID := s.app.Config.GetString("app_luexu.grpc_service_id", "app_luexu-grpc")
    err = s.app.Config.RegisterGRPCService(serviceName, serviceID, registerAddr, checkAddr, port)

    return listener, serviceID, err
}