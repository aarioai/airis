package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"{{APP_BASE}}/grpc/helloworld"
	"{{PROJECT_BASE}}/proto/{{APP_NAME}}/pb"
)

func (s *Service) registerServer() *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterHelloWorldServer(server, helloworld.NewHelloWorld(s.s))

    // register GRPC health check
    // Can use `grpc_health_probe -addr=localhost:60000` to check GRPC health
    healthServer := health.NewServer()
    grpc_health_v1.RegisterHealthServer(server, healthServer)
    healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

    // Enable `grpcurl -plaintext localhost:50000 list` to list all services and methods
    if s.app.Config.Env.BeforeDevelopment() {
        reflection.Register(server)
    }
    return server
}
