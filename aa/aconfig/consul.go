package aconfig

import "github.com/hashicorp/consul/api"

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
