package consul

import (
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

// Builder Interface: https://pkg.go.dev/google.golang.org/grpc/resolver#Builder

type Builder struct {
	client *api.Client
}

func NewBuilder(client *api.Client) resolver.Builder {
	return &Builder{
		client: client,
	}
}

func (b *Builder) Build(t resolver.Target, cc resolver.ClientConn, o resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		cc:           cc,
		client:       b.client,
		serviceName:  t.Endpoint(),
		closeChannel: make(chan struct{}),
	}
	r.ResolveNow(defaultResolveNowOptions)
	go r.Watch()
	return r, nil
}

func (b *Builder) Scheme() string {
	return "consul"
}
