package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"{{APP_BASE}}/grpc/server/helloworld"
	"{{APP_BASE}}/grpc/server/pb"
)

func (s *Service) registerServer() *grpc.Server {
	serve := grpc.NewServer()
	pb.RegisterHelloWorldServer(serve, helloworld.NewHelloWorld(s.s))

    // register GRPC health check
    // Can use `grpc_health_probe -addr=localhost:60000` to check GRPC health
    healthServer := health.NewServer()
    grpc_health_v1.RegisterHealthServer(serve, healthServer)
    healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

    return serve
}
