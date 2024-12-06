package atype

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

// AByte String('A') will returns "97". So you have to use String(AByte('A')) to return "A"
type AByte byte

func IsEmpty(d any) bool {
	return !NotEmpty(d)
}

/*
在Go语言中，一个any类型的变量包含了2个指针，一个指针指向值的在编译时确定的类型，另外一个指针指向实际的值。
*/
func IsNil(x any) bool {
	// @warn 断言和反射性能不是特别好，如果不得已再使用，控制使用有助于提升程序性能。
	return x == nil || (reflect.ValueOf(x).Kind() == reflect.Ptr && reflect.ValueOf(x).IsNil())
}

// NotEmpty check if a value is empty
// @warn NotEmpty(byte(0)) == false,  NotEmpty(byte('0')) == true
//
//	NotEmpty(0) == false, NotEmpty(-1) == true, NotEmpty(1) == true
func NotEmpty(d any) bool {
	if d == nil {
		return false
	}
	switch v := d.(type) {
	case AByte:
		return v != 0
	case bool:
		return v
	case string:
		return v != ""
	case int8:
		return v != 0 // 复数，不算 empty，只有0 才算empty
	case int16:
		return v != 0
	case rune:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	case byte:
		return v > 0
	case Booln:
		return v > 0
	case uint16:
		return v > 0
	case Uint24:
		return v > 0
	case uint32:
		return v > 0
	case uint:
		return v > 0
	case uint64:
		return v > 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	}

	return String(d) != ""
}

// MarshalUint8s 将 uint8 切片转换为 JSON 数组
// json.Marshal() 不能正常转换 []byte 及 []uint8
func MarshalUint8s(bytes []uint8) (json.RawMessage, bool) {
	if len(bytes) == 0 {
		return nil, true
	}

	// 预分配足够的空间
	result := make([]byte, 0, len(bytes)*4+2)
	result = append(result, '[')

	for i, v := range bytes {
		if i > 0 {
			result = append(result, ',')
		}
		result = strconv.AppendUint(result, uint64(v), 10)
	}

	result = append(result, ']')
	return result, true
}
func UnmarshalUint8s(x json.RawMessage) ([]uint8, error) {
	if x == nil || len(x) < 2 {
		return nil, nil
	}

	n := len(x) - 1 // remove last ']'
	v := make([]uint8, 0)
	for i := 1; i < n; {
		for x[i] == ' ' || x[i] == ',' {
			i++
		}
		if x[i] < '0' || x[i] > '9' {
			return nil, errors.New("invalid uint8 json: " + string(x))
		}
		var c string
		for x[i] >= '0' && x[i] <= '9' {
			c += string(x[i])
			i++
		}
		u, _ := strconv.ParseUint(c, 10, 8)
		v = append(v, uint8(u))
	}
	return v, nil
}
