package grpc

import (
	"google.golang.org/grpc"
	"{{APP_BASE}}/grpc/server/helloworld"
	"{{APP_BASE}}/grpc/server/helloworld/pb"
)

func newServer() *grpc.Server {
	serve := grpc.NewServer()
	pb.RegisterHelloWorldServer(serve, &helloworld.HelloWorld{})
	return serve
}
