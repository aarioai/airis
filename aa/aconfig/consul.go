package aconfig

import (
	"fmt"
	"github.com/aarioai/airis/aa/aconfig/consul"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

func (c *Config) SetConsul(name string, client *api.Client) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	c.consulMap[name] = client
}

func (c *Config) SetDefaultConsul(client *api.Client) {
	cfgMtx.Lock()
	defer cfgMtx.Unlock()
	c.consulMap["DEFAULT"] = client
}

func (c *Config) Consul(name string) *api.Client {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	return c.consulMap[name] // panic on doesn't exist
}

func (c *Config) DefaultConsul() *api.Client {
	cfgMtx.RLock()
	defer cfgMtx.RUnlock()
	return c.consulMap["DEFAULT"] // panic on doesn't exist
}

func (c *Config) RegisterGRPCService(serviceName, serviceID, address string, port int) error {
	reg := api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", address, port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "5m",
		},
	}
	client := c.DefaultConsul()
	resolver.Register(consul.NewBuilder(client))
	return client.Agent().ServiceRegister(&reg)
}
func (c *Config) DeregisterGRPCService(serviceID string) error {
	
	return c.DefaultConsul().Agent().ServiceDeregister(serviceID)
}

// DiscoverGRPCServices
// Better use grpc/resolver
func (c *Config) DiscoverGRPCServices(serviceName string) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	services, meta, err := c.DefaultConsul().Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, nil, err
	}
	if len(services) == 0 {
		return nil, nil, fmt.Errorf("no grpc service instance %s available", serviceName)
	}
	return services, meta, nil
}
