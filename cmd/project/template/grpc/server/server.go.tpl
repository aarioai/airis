package server

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
    "github.com/aarioai/airis/pkg/basic"
    "github.com/aarioai/airis/pkg/types"
    "google.golang.org/grpc"
	"net"
)

func (s *Service) Serve(prof *debug.Profile) {
	prof.Fork("starting grpc server ({{APP_NAME}})")

	listener, serviceID, err := s.listen()
	if err != nil {
		return nil, err
	}
	server := s.registerServer()

	go func() {
		<-s.app.GlobalContext.Done()
		alog.Stopf("grpc server ({{APP_NAME}})")
		s.app.Config.DeregisterGRPCService(serviceID)
		server.GracefulStop()
	}()

	go func() {
		defer s.app.Config.DeregisterGRPCService(serviceID)
		ae.PanicOnErrs(server.Serve(listener))
	}()

	return server, nil
}

func (s *Service) listen() (net.Listener, string, error) {
	addr := s.app.Config.GetString("infra.grpc_addr", "127.0.0.1")
	checkAddr := s.app.Config.GetString("infra.grpc_check_addr", addr)

    port, _ := types.ParseInt(s.app.Config.GetString("{{APP_NAME}}.grpc_port"))
    if port <= 0 {
        return nil, "", errors.New("missing or invalid config {{APP_NAME}}.grpc_port")
    }

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
    if err != nil {
        return nil, "", err
    }

    serviceName := s.app.Config.GetString("{{APP_NAME}}.grpc_service_name", "{{APP_NAME}}")
    serviceID := s.app.Config.GetString("{{APP_NAME}}.grpc_service_id", "{{APP_NAME}}-grpc")
    err = s.app.Config.RegisterGRPCService(serviceName, serviceID, addr, checkAddr, port)

    return listener, serviceID, err
}