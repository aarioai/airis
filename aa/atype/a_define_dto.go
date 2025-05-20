package atype

type Paging struct {
	Page     uint   `json:"page"`
	PageEnd  uint   `json:"page_end"`
	PageSize uint8  `json:"page_size"`
	Offset   uint   `json:"offset"`
	Limit    uint16 `json:"limit"`
	Prev     uint   `json:"prev"`
	Next     uint   `json:"next"`
}

type Location struct {
	Valid     bool    `json:"-"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"` // 保留
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"` // 保留
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
