package atype

type Paging struct {
	Page     int `json:"page"`
	PageEnd  int `json:"page_end"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
	Prev     int `json:"prev"`
	Next     int `json:"next"`
}

const (
	DefaultPageSize = 10
)

// NewPaging
// @param pageStart page start
// @param [ends] [pageEnd[, pageSize]]
func NewPaging(pageStart int, ends ...int) Paging {
	if pageStart <= 1 {
		pageStart = 1
	}
	pageEnd := pageStart
	if len(ends) > 0 && ends[0] > 0 {
		pageEnd = ends[0]
	}
	pageSize := DefaultPageSize
	if len(ends) > 1 && ends[1] > 0 {
		pageSize = ends[1]
	}

	var offset int
	var prev int

	next := pageEnd + 1
	limit := pageSize * (next - pageStart)
	offset = (pageStart - 1) * pageSize
	prev = pageStart - 1

	return Paging{
		Page:     pageStart,
		PageEnd:  pageEnd,
		PageSize: pageSize,
		Offset:   offset,
		Limit:    limit,
		Prev:     prev,
		Next:     next,
	}
}
