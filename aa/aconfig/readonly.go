package aconfig

import "sync"

type Readonly struct {
	data sync.Map
}

func NewReadonly() *Readonly {
	return &Readonly{
		data: sync.Map{},
	}
}

func (c *Readonly) Set(key string, value string) bool {
	_, loaded := c.data.LoadOrStore(key, value)
	return !loaded
}

func (c *Readonly) Get(key string) string {
	if value, ok := c.data.Load(key); ok {
		return value.(string)
	}
	return ""
}

func (c *Readonly) GetBytes(key string) []byte {
	s := c.Get(key)
	if s == "" {
		return nil
	}
	return []byte(s)
}
