package client

import (
	"github.com/aarioai/airis/aa/helpers/debug"
)

func (s *Service) Run(prof *debug.Profile) {
	//profile := prof.Fork("staring grpc client ({{APP_NAME}})")
	//time.Sleep(10 * time.Second)
	//addr := "localhost:8000"
	//conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//fmt.Println(err)
	//defer conn.Close()
	//
	//go func() {
	//	<-s.app.GlobalContext.Done()
	//	alog.Stopf("grpc client ({{APP_NAME}}:%s)", addr)
	//	if conn != nil {
	//		conn.Close()
	//	}
	//}()
	//
	//c := pb.NewHelloWorldClient(conn)
	//
	//ctx, cancel := context.WithTimeout(s.app.GlobalContext, 10*time.Second)
	//defer cancel()
	//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	//if err != nil {
	//	ae.PanicF("could not greet: %v", err)
	//}
	//
	//fmt.Printf("Greeting: %s\n", r.GetMessage())
}