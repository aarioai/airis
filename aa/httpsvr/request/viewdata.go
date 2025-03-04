package request

func (r *Request) QueryToViewData(names ...string) {
	for _, name := range names {
		v, e := r.QueryString(name, false)
		if e == nil {
			r.ictx.ViewData(name, v)
		}
	}
}
