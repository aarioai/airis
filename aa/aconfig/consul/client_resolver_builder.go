package consul

import (
	"github.com/aarioai/airis/pkg/basic"
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

// Build creates a new resolver for the given target.
// grpc.NewClient calls Build synchronously, and fails if the returned error is not nil.
func (b *Builder) Build(t resolver.Target, cc resolver.ClientConn, o resolver.BuildOptions) (resolver.Resolver, error) {
	serviceName := basic.Ter(t.Endpoint() != "", t.Endpoint(), t.URL.Host)
	r := NewResolver(cc, b.client, serviceName, make(chan struct{}))
	r.ResolveNow(defaultResolveNowOptions)
	go r.Watch()
	return r, nil
}

// Scheme returns the scheme supported by this resolver.  Scheme is defined
// at https://github.com/grpc/grpc/blob/master/doc/naming.md.  The returned
// string should not contain uppercase characters, as they will not match
// the parsed target's scheme as defined in RFC 3986.
func (b *Builder) Scheme() string {
	return Scheme
}
