package atype

import (
	"errors"
	"github.com/aarioai/airis/pkg/types"
)

var (
	// Municipalities: Beijing, Shanghai, Tianjin, Chongqing
	Municipalities = []Distri{110000, 310000, 120000, 500000}
	// Special Economic Zones: Hainan, Shenzhen, Xiamen, Zhuhai, Shantou
	SEZs = []Distri{460000, 440300, 350200, 440400, 440500}
	// Autonomous regions: Xinjiang, Tibet, Ningxia, Inner Mongolia, Guangxi
	AutonomousRegions = []Distri{650000, 540000, 640000, 150000, 450000}
	// Special Administrative Regions: Hong Kong, Macau
	SARs = []Distri{810000, 820000}
)

func NewDistri(d Uint24) Distri {
	// Keep 6 digits
	if d == 0 {
		return 0
	}
	if d < 100 {
		return Distri(d * 10000)
	}
	if d < 10000 {
		return Distri(d * 100)
	}
	return Distri(d)
}
func ParseDistri(s string) (Distri, error) {
	n, err := ParseUint24(s)
	if err != nil {
		return 0, err
	}
	return NewDistri(n), nil
}
func ToDistri(d uint32) Distri {
	return NewDistri(Uint24(d))
}
func (d Distri) Uint24() Uint24 {
	return Uint24(d)
}
func (d Distri) Uint32() uint32 {
	return uint32(d)
}
func (d Distri) String() string {
	return types.FormatUint(d)
}
func (d Distri) Province() Province {
	return Province(d / 10000)
}
func (d Distri) Dist() Dist {
	return Dist(d / 100)
}
func (d Distri) AddrId() District {
	return NewAddrId(uint64(d) * 1000000)
}

func ParseDist(s string) (Dist, error) {
	n, err := types.ParseUint16(s)
	if err != nil {
		return 0, err
	}
	dist := Dist(n)
	if dist == 0 {
		return 0, errors.New("invalid distri string")
	}
	return dist, nil
}

// Whether one district is inside another
func (d Distri) Inside(p Distri) bool {
	if d == p {
		return true
	}
	b := d % 100
	if b != 0 && (d-b) == p {
		return true
	}
	b = d % 10000
	if b != 0 && (d-b) == p {
		return true
	}
	return false
}
func (d Distri) In(distris []Distri) bool {
	for _, distri := range distris {
		if d == distri {
			return true
		}
	}
	return false
}
func (d Distri) IsMunicipality() bool { return d.In(Municipalities) }
func (d Distri) IsSEZ() bool          { return d.In(SEZs) }
func (d Distri) IsAutonomous() bool   { return d.In(AutonomousRegions) }
func (d Distri) IsSAR() bool          { return d.In(SARs) }
func (d Distri) IsProvLevel() bool    { return d%10000 == 0 }
func (d Distri) IsProvince() bool     { return d.IsProvLevel() && !d.In(Municipalities) }

func ParseProvince(s string) (Province, error) {
	n, err := types.ParseUint8(s)
	if err != nil {
		return 0, err
	}
	return Province(n), nil
}
func (p Province) Distri() Distri {
	return Distri(uint32(p) * 10000)
}
func (p Province) String() string {
	return types.FormatUint(p)
}
func (d Dist) Distri() Distri {
	return Distri(uint32(d) * 100)
}
func (d Dist) String() string {
	return types.FormatUint(d)
}

func ParseAddrId(s string) (District, error) {
	n, err := types.ParseUint64(s)
	if err != nil {
		return 0, err
	}
	return District(n), nil
}
func NewAddrId(a uint64) District {
	return District(a)
}
func (a District) Uint64() uint64 {
	return uint64(a)
}
func (a District) String() string {
	return types.FormatUint(a)
}
