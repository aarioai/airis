package httpsvr

func (c *Context) QueryToViewData(names ...string) {
	for _, name := range names {
		v, e := c.QueryString(name, false)
		if e == nil {
			c.ViewData(name, v)
		}
	}
}
