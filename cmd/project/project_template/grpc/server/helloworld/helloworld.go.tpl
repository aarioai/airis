package helloworld

import (
	"context"
	"log"
	"{{APP_BASE}}/grpc/server/helloworld/pb"
)

type HelloWorld struct {
	pb.UnimplementedHelloWorldServer
}

func (s *HelloWorld) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received message: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
