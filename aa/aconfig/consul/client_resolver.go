package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

const (
	defaultConsulScheme  = "consul"
	defaultWatchInterval = 15 * time.Second
)

var (
	defaultResolveNowOptions = resolver.ResolveNowOptions{}
)

// Resolver Interface: https://pkg.go.dev/google.golang.org/grpc/resolver#Resolver

type Resolver struct {
	cc           resolver.ClientConn
	client       *api.Client
	serviceName  string
	closeChannel chan struct{}
}

func NewResolver(cc resolver.ClientConn, client *api.Client, serviceName string, closeChannel chan struct{}) *Resolver {
	return &Resolver{
		cc:           cc,
		client:       client,
		serviceName:  serviceName,
		closeChannel: closeChannel,
	}
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {
	services, _, err := r.client.Health().Service(r.serviceName, "", true, nil)
	if err != nil {
		r.cc.ReportError(fmt.Errorf("consul resolve health service (%s) failed: %v", r.serviceName, err))
		return
	}
	if len(services) == 0 {
		r.cc.ReportError(fmt.Errorf("consul resolve no health service (%s)", r.serviceName))
		return
	}

	var addrs []resolver.Address
	for _, s := range services {
		addr := s.Service.Address
		if addr == "" {
			addr = s.Node.Address
		}
		attrs := attributes.New("consul.service.id", s.Service.ID).
			WithValue("consul.service.meta", &s.Service.Meta).
			WithValue("consul.node.id", s.Node.ID).
			WithValue("consul.node.name", s.Node.Node)

		addrs = append(addrs, resolver.Address{
			Addr:       fmt.Sprintf("%s:%d", addr, s.Service.Port),
			ServerName: s.Service.Service,
			Attributes: attrs,
		})

		fmt.Printf("consul resolve %s (%s:%d) %s\n", s.Service.Service, addr, s.Service.Port, attrs.String())
	}

	err = r.cc.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: r.cc.ParseServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
	})
	if err != nil {
		r.cc.ReportError(fmt.Errorf("consul update state failed: %v", err))
	}
}

// Close resolver
// Note: it is automatically invoked by grpc.ClientConn.Close()
func (r *Resolver) Close() {
	close(r.closeChannel)
}

func (r *Resolver) Watch() {
	ticker := time.NewTicker(defaultWatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.ResolveNow(defaultResolveNowOptions)
		case <-r.closeChannel:
			return
		}
	}
}
