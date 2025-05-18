package atype

import "github.com/aarioai/airis/pkg/afmt"

type Paging struct {
	Page     uint   `json:"page"`
	PageEnd  uint   `json:"page_end"`
	PageSize uint8  `json:"page_size"`
	Offset   uint   `json:"offset"`
	Limit    uint16 `json:"limit"`
	Prev     uint   `json:"prev"`
	Next     uint   `json:"next"`
}

var (
	DefaultPageSize = uint8(10)
	MaxPageRange    = uint8(5)
)

func NewPaging(pageStart, pageEnd uint, pageSizes ...uint8) Paging {
	if pageStart <= 1 {
		pageStart = 1
	}
	if pageEnd < pageStart {
		pageEnd = pageStart
	}
	pageSize := afmt.First(pageSizes)
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}

	var offset uint
	var prev uint

	next := pageEnd + 1
	pageRange := next - pageStart
	if pageRange > uint(MaxPageRange) {
		pageEnd = pageStart + 1
		next = pageEnd + 1
	}

	limit := uint16(pageSize) * uint16(next-pageStart)
	offset = (pageStart - 1) * uint(pageSize)
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
