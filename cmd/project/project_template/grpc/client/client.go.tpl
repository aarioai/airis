package client

import (
	"github.com/aarioai/airis/aa/helpers/debug"
)

func (s *Service) Run(prof *debug.Profile) {
	prof.Fork("staring grpc client ({{APP_NAME}})")
	//time.Sleep(5 * time.Second)
	//
	//addr := "localhost:8000"
	//conn, err := grpc.NewClient(addr,
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	//)
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//
	//go func() {
	//	<-s.app.GlobalContext.Done()
	//	alog.Stopf("grpc client (infra:%s)", addr)
	//	if conn != nil {
	//		conn.Close()
	//	}
	//}()
	//
	//c := pb.NewHelloWorldClient(conn)
	//
	//for {
	//	t := time.NewTicker(3 * time.Second)
	//	defer t.Stop()
	//	for range t.C {
	//		ctx, cancel := context.WithTimeout(s.app.GlobalContext, 5*time.Second)
	//		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	//		cancel()
	//		if err != nil {
	//			ae.PanicF("could not greet: %v", err)
	//		}
	//
	//		fmt.Printf("Greeting: %s\n", r.GetMessage())
	//
	//	}
	//}
}