package client

import (
	"github.com/aarioai/airis/aa/helpers/debug"
)

func (s *Service) Run(prof *debug.Profile) {
	//conn, _, err := helloworldsdk.New(s.app).Conn()
	//ae.PanicOnErrs(err)
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