package aconfig

type Readonly struct {
	data map[string]string
}

func NewReadonly() *Readonly {
	return &Readonly{
		data: make(map[string]string),
	}
}
func (c *Readonly) Init(key string, value string) bool {
	if _, ok := c.data[key]; ok {
		return false
	}
	c.data[key] = value
	return true
}

func (c *Readonly) Get(key string) string {
	if v, ok := c.data[key]; ok {
		return v
	}
	return ""
}

func (c *Readonly) GetBytes(key string) []byte {
	if v, ok := c.data[key]; ok && v != "" {
		return []byte(v)
	}
	return nil
}
