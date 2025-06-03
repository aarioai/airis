package aenum

import (
	"github.com/aarioai/airis/pkg/basic"
	"github.com/aarioai/airis/pkg/types"
	"strings"
)

type Sex uint8

const (
	NilSex   Sex = 0 // no, or invalid sex
	Male     Sex = 1
	Female   Sex = 2
	OtherSex Sex = 255
)

func (x Sex) Valid() bool {
	return x == Male || x == Female || x == OtherSex
}

func NewSex(sex uint8) Sex {
	x := Sex(sex)
	return basic.Ter(x.Valid(), x, NilSex)
}

func ToSex(x string) Sex {
	x = strings.ToUpper(x)
	switch x {
	case "1", "M", "MALE", "MAN", "男":
		return Male
	case "2", "F", "FEMALE", "WOMAN", "女":
		return Female
	case "255":
		return OtherSex
	default:
		return NilSex
	}
}
func (x Sex) Uint8() uint8   { return uint8(x) }
func (x Sex) String() string { return types.FormatUint8(x.Uint8()) }
