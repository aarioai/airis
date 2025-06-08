package grpc

import (
	"google.golang.org/grpc"
	"{{APP_BASE}}/grpc/helloworld"
	"{{APP_BASE}}/grpc/helloworld/build"
)

func newServer() *grpc.Server {
	serve := grpc.NewServer()
	build.RegisterHelloWorldServer(serve, &helloworld.HelloWorld{})
	return serve
}
