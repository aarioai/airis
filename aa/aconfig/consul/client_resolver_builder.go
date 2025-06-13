package consul

import (
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const Scheme = "consul"

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
	r := NewResolver(cc, b.client, t.Endpoint(), 0, make(chan struct{}))
	r.ResolveNow(defaultResolveNowOptions)
	go r.Watch()
	return r, nil
}

func (b *Builder) Scheme() string {
	return Scheme
}
