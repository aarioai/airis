package atype

import "github.com/aarioai/airis/aa/ae"

// QuerySource is the basic Query surface httpsvr.Context provides.
type QuerySource interface {
	QueryBool(p string) (bool, *ae.Error)
	QueryString(p string, params ...any) (string, *ae.Error)
	QueryUint(p string, required ...bool) (uint, *ae.Error)
	QueryUint8(p string, required ...bool) (uint8, *ae.Error)
	QueryUint16(p string, required ...bool) (uint16, *ae.Error)
	QueryInt8(p string, required ...bool) (int8, *ae.Error)
	QueryInt64(p string, required ...bool) (int64, *ae.Error)
	QueryInt64s(p string, required, allowZero bool) ([]int64, *ae.Error)
	QueryUint64(p string, required ...bool) (uint64, *ae.Error)
	QueryUint64s(p string, required, allowZero bool) ([]uint64, *ae.Error)
}

// BodySource is the basic Body surface httpsvr.Context provides.
type BodySource interface {
	BodyPath(p string, required ...bool) (string, *ae.Error)
	BodyPaths(p string, required ...bool) ([]string, *ae.Error)
	BodyBool(p string) (bool, *ae.Error)
	BodyString(p string, patterns ...any) (string, *ae.Error)
	BodyUint8(p string, required ...bool) (uint8, *ae.Error)
	BodyUint16(p string, required ...bool) (uint16, *ae.Error)
	BodyInt8(p string, required ...bool) (int8, *ae.Error)
	BodyInt64(p string, required ...bool) (int64, *ae.Error)
	BodyInt64s(p string, required, allowZero bool) ([]int64, *ae.Error)
	BodyUint64(p string, required ...bool) (uint64, *ae.Error)
	BodyUint64s(p string, required, allowZero bool) ([]uint64, *ae.Error)
	BodyFloat64Map(p string, requireds ...bool) (map[string]float64, *ae.Error)
	BodyAnyMap(p string, requireds ...bool) (map[string]any, *ae.Error)
}
