package helloworld

import (
	"context"
	"log"
	"project/microservice/app/app_luexu/service"
	"project/microservice/proto/app_luexu/pb"
)

type HelloWorld struct {
    s *service.Service
	pb.UnimplementedHelloWorldServer
}

func NewHelloWorld(s *service.Service) *HelloWorld {
	return &HelloWorld{
	    s: s,
	}
}

func (s *HelloWorld) SayHello(_ context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received message: %v", r.GetName())
	return &pb.HelloReply{Message: "Hello " + r.GetName()}, nil
}
