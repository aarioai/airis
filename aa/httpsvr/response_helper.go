package httpsvr

import (
	"github.com/aarioai/airis/aa/ae"
)

func (c *Context) FirstError(es ...*ae.Error) *ae.Error {
	e := ae.First(es...)
	if e != nil {
		c.WriteE(e)
		return e
	}
	return nil
}

func (c *Context) IsOK(catchResponded bool, e *ae.Error) bool {
	if catchResponded {
		return false
	}
	if e != nil {
		c.WriteE(e)
		return false
	}
	return true
}

func (c *Context) OK(catchResponded bool, e *ae.Error) {
	if ok := c.IsOK(catchResponded, e); ok {
		c.WriteOK()
	}
}
