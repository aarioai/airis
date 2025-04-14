package response

import (
	"github.com/aarioai/airis/aa/ae"
)

func (w *Writer) FirstError(es ...*ae.Error) *ae.Error {
	e := ae.First(es...)
	if e != nil {
		w.WriteE(e)
		return e
	}
	return nil
}

func (w *Writer) IsOK(catchResponded bool, e *ae.Error) bool {
	if catchResponded {
		return false
	}
	if e != nil {
		w.WriteE(e)
		return false
	}
	return true
}

func (w *Writer) OK(catchResponded bool, e *ae.Error) {
	if ok := w.IsOK(catchResponded, e); ok {
		w.WriteOK()
	}
}
