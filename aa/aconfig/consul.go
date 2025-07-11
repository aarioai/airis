package aconfig

import (
	"fmt"
	"github.com/aarioai/airis/pkg/nets"
	"github.com/hashicorp/consul/api"
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

func normalizeAddress(addr string) string {
	if addr == "" || addr == "localhost" || addr == "0.0.0.0" || addr == "::" || addr == "[::]" {
		if ip, err := nets.LanIP(); err == nil {
			return ip
		}
		return "127.0.0.1"
	}
	return addr
}

// RegisterGRPCService
// Note: sometimes in docker container, remote address may be different with check address
func (c *Config) RegisterGRPCService(serviceName, serviceID, address, checkAddr string, port int) error {
	address = normalizeAddress(address)
	checkAddr = normalizeAddress(checkAddr)
	reg := api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", checkAddr, port),
			Interval:                       "15s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "5m",
		},
	}
	client := c.DefaultConsul()

	// Register self to consul, each service instance should register once
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
