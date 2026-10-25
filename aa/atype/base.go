package atype

import (
	"encoding/json/jsontext"
	"errors"
	"github.com/aarioai/airis/pkg/types"
	"strconv"
)

// MarshalUint8s converts a uint8 slice to a JSON array of numbers.
// encoding/json/v2 represents []byte as Base64; this keeps a numeric array.
func MarshalUint8s(bytes []uint8) (jsontext.Value, bool) {
	if len(bytes) == 0 {
		return nil, true
	}

	// Preallocate enough space
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

func UnmarshalUint8s(x jsontext.Value) ([]uint8, error) {
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
		u, _ := types.ParseUint8(c)
		v = append(v, u)
	}
	return v, nil
}
