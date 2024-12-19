package atype_test

import (
	"github.com/aarioai/airis/core/atype"
	"math"
	"math/rand/v2"
	"strconv"
	"testing"
)

func TestParseInt(t *testing.T) {
	// 测试小数字
	for i := 0; i < math.MaxUint8; i++ {
		s := atype.FormatInt(i)
		got, err := atype.ParseInt(s)
		if err != nil {
			t.Errorf("FormatInt(%d) failed: %s", i, err.Error())
			continue
		}
		if got != i {
			t.Errorf("FormatInt(%d) got %s, want %d", i, s, i)
		}
	}

	// 测试FormatInt/ParseInt64
	for i := 0; i < 1000; i++ {
		n := rand.Int64N(math.MaxInt64)
		if i%2 == 0 {
			n = -n
		}
		want := strconv.FormatInt(n, 10)
		s := atype.FormatInt(n)
		if s != want {
			t.Errorf("FormatInt(%d) => %s, want %s", n, s, want)
		}
		got, err := atype.ParseInt64(s)
		if err != nil {
			t.Errorf("ParseInt64(%d) => %v, want nil", n, err)
			continue
		}
		if got != n {
			t.Errorf("ParseInt64(%d) => %d, want %d", n, got, n)
		}
	}

	for i := 0; i < 1000; i++ {
		n := rand.Uint64N(math.MaxInt64)
		want := strconv.FormatUint(n, 10)
		s := atype.FormatUint(n)
		if s != want {
			t.Errorf("FormatUint(%d) => %s, want %s", n, s, want)
		}
		got, err := atype.ParseUint64(s)
		if err != nil {
			t.Errorf("ParseUint64(%d) => %v, want nil", n, err)
			continue
		}
		if got != n {
			t.Errorf("ParseUint64(%d) => %d, want %d", n, got, n)
		}
	}
}
func TestConvertBase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fromBase int
		toBase   int
		want     string
		wantErr  bool
	}{
		// 基本测试用例
		//	{"decimal to hex", "255", 10, 16, "ff", false},
		{"hex to decimal", "ff", 16, 10, "255", false},
		//{"binary to decimal", "1010", 2, 10, "10", false},
		//{"decimal to binary", "10", 10, 2, "1010", false},
		//
		//// 负数测试
		//{"negative decimal to hex", "-255", 10, 16, "-ff", false},
		//{"negative hex to decimal", "-ff", 16, 10, "-255", false},
		//
		//// 边界情况
		//{"empty string", "", 10, 16, "", true},
		//{"zero value", "0", 10, 16, "0", false},
		//
		//// 大数测试
		//{"large number", "9223372036854775807", 10, 16, "7fffffffffffffff", false},
		//{"negative large number", "-9223372036854775808", 10, 16, "-8000000000000000", false},
		//
		//// 无效输入
		//{"invalid hex", "gg", 16, 10, "", true},
		//{"invalid decimal", "12a", 10, 16, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := atype.ConvertBase(tt.input, tt.fromBase, tt.toBase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertBase(%s, %d, %d) error = %v, wantErr %v", tt.input, tt.fromBase, tt.toBase, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertBase(%s, %d, %d) = %v, want %v", tt.input, tt.fromBase, tt.toBase, got, tt.want)
			}
		})
	}
}
