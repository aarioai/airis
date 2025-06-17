package helloworld

import (
	"context"
	"log"
	"project/simple/app/app_aajs/service"
	"project/simple/proto/app_aajs/pb"
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
