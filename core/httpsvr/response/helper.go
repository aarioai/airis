package response

import (
	"github.com/aarioai/airis/core/ae"
)

func (w *Writer) FirstError(es ...*ae.Error) *ae.Error {
	e := ae.First(es...)
	if e != nil {
		w.WriteE(e)
		return e
	}
	return nil
}
