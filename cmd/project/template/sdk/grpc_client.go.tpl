package {{APP_NAME}}

import (
	"context"
	"github.com/aarioai/airis/aa/aconfig/consul"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/alog"
	"github.com/aarioai/airis/aa/helpers/debug"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const (
	initConnectTimeout        = 10 * time.Second
	connectTimeout            = 5 * time.Second
	connectStateChangeTimeout = 500 * time.Millisecond
)

func (s *Service) Init(prof *debug.Profile) {
	prof.Forkf("initial grpc client ({{APP_NAME}}: %s)", s.target)
	ae.PanicOn(s.initGRPCClient())

    if ok := s.waitForConnectReady(s.conn); !ok {
        ae.PanicF("initial grpc client ({{APP_NAME}}: %s) did not become ready", s.target)
    }

	go s.watchTerminate()
}

func (s *Service) Conn() (*grpc.ClientConn, string, *ae.Error) {
	s.mtx.RLock()
	conn, target := s.conn, s.target
	s.mtx.RUnlock()

	if s.conn != nil {
		switch conn.GetState() {
		case connectivity.Ready:
			return conn, target, nil
		case connectivity.Idle:
			conn.Connect()
			return conn, target, nil
		case connectivity.Connecting, connectivity.TransientFailure:
			// Wait briefly for connection to recover
			ctx, cancel := context.WithTimeout(context.Background(), connectStateChangeTimeout)
			defer cancel()
			if conn.WaitForStateChange(ctx, conn.GetState()) {
				if conn.GetState() == connectivity.Ready {
					return conn, target, nil
				}
			}
		}
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	// Double-check in case another goroutine already recreated the connection
	if s.conn != nil && s.conn.GetState() == connectivity.Ready {
		return s.conn, s.target, nil
	}

    if e := s.initGRPCClient(); e != nil {
		return nil, "", e
	}

	return s.conn, s.target, nil
}

func (s *Service) initGRPCClient() *ae.Error {
	serviceName := s.app.Config.GetString("{{CONFIG_SECTION}}.grpc_service_name", "{{APP_NAME}}")
	addr := consul.Scheme + ":///" + serviceName

	// gRPC is highly efficient and reusable, no need to write a connection pool
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingPolicy":"round_robin",
			"healthCheckConfig": {
				"serviceName": ""
			}
		}`),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: connectTimeout,
		}),
	)
	if err != nil {
		return ae.NewF(ae.GatewayTimeout, "failed to create gRPC client for %s: %v", addr, err.Error())
	}
	s.conn = conn
	s.target = addr
	return nil
}

func (s *Service) waitForConnectReady(conn *grpc.ClientConn) bool {
	ctx, cancel := context.WithTimeout(s.app.GlobalContext, initConnectTimeout)
	defer cancel()

	for {
		state := conn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			return true
		}
		if !conn.WaitForStateChange(ctx, state) {
			return false // Timeout
		}
	}
}

func (s *Service) watchTerminate() {
	// Wait for application shutdown, including SIGINT, SIGTERM
	<-s.app.GlobalContext.Done()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	alog.Stopf("grpc client ({{APP_NAME}}: %s)", s.target)
	if s.conn != nil {
		s.conn.Close()
	}
}
