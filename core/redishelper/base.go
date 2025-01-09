package redishelper

import (
	"reflect"
	"strings"
	"time"
)

// TTL 常量定义
const (
	HourlyTTL  = 24 * time.Hour // 24小时 TTL，用于小时级缓存   要求每小时会自动清除之前表；为了避免宕机等影响，ttl设计长一点，24小时内宕机恢复，就能使用
	DailyTTL   = 72 * time.Hour // 3天 TTL，用于天级缓存  要求每天会自动清除之前表；为了避免宕机等影响，ttl设计长一点
	DefaultTTL = time.Hour      // 默认 TTL
)

func omitEmpty(opt string) bool {
	for opt != "" {
		var name string
		name, opt, _ = strings.Cut(opt, ",")
		if name == "omitempty" {
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

// appendStructField appends the field and value held by the structure v to dst, and returns the appended dst.
func appendStructField(dst []interface{}, v reflect.Value) []interface{} {
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("redis")
		if tag == "" || tag == "-" {
			continue
		}
		name, opt, _ := strings.Cut(tag, ",")
		if name == "" {
			continue
		}

		field := v.Field(i)

		// miss field
		if omitEmpty(opt) && isEmptyValue(field) {
			continue
		}

		if field.CanInterface() {
			dst = append(dst, name, field.Interface())
		}
	}

	return dst
}
