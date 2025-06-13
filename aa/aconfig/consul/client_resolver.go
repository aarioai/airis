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
	lastIndex    uint64
	closeChannel chan struct{}
}

func NewResolver(cc resolver.ClientConn, client *api.Client, serviceName string, lastIndex uint64, closeChannel chan struct{}) *Resolver {
	return &Resolver{
		cc:           cc,
		client:       client,
		serviceName:  serviceName,
		lastIndex:    lastIndex,
		closeChannel: closeChannel,
	}
}

func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {
	services, meta, err := r.client.Health().Service(r.serviceName, "", true, &api.QueryOptions{WaitIndex: r.lastIndex})
	if err != nil {
		r.cc.ReportError(fmt.Errorf("consul resolver: query failed: %v", err))
		return
	}
	if len(services) == 0 {
		r.cc.ReportError(fmt.Errorf("consul resolver: no addresses found for %s", r.serviceName))
		return
	}

	r.lastIndex = meta.LastIndex
	var addrs []resolver.Address
	for _, s := range services {
		addr := s.Service.Address
		if addr == "" {
			addr = s.Node.Address
		}
		attrs := attributes.New("consul.service.id", s.Service.ID).
			WithValue("consul.service.meta", s.Service.Meta).
			WithValue("consul.node.id", s.Node.ID).
			WithValue("consul.node.name", s.Node.Node)

		addrs = append(addrs, resolver.Address{
			Addr:       fmt.Sprintf("%s:%d", addr, s.Service.Port),
			ServerName: s.Service.Service,
			Attributes: attrs,
		})
		fmt.Printf("consul resolve %s (%s:%d) %s\n", s.Service.Service, addr, s.Service.Port, attrs.String())
	}

	r.cc.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: r.cc.ParseServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
	})
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
