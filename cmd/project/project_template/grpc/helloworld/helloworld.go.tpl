package helloworld

import (
	"context"
	"log"
	"{{APP_BASE}}/grpc/helloworld/build"
)

type HelloWorld struct {
	build.UnimplementedHelloWorldServer
}

func (s *HelloWorld) SayHello(_ context.Context, in *build.HelloRequest) (*build.HelloReply, error) {
	log.Printf("Received message: %v", in.GetName())
	return &build.HelloReply{Message: "Hello " + in.GetName()}, nil
}
