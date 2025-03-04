package atype_test

import (
	"github.com/aarioai/airis/aa/atype"
	"testing"
)

func TestString(t *testing.T) {
	if atype.String(false) != "false" {
		t.Errorf("bool(false) ==> string(%s)", atype.String(false))
	}
	if atype.String(true) != "true" {
		t.Errorf("bool(true) ==> string(%s)", atype.String(true))
	}

	// byte is a built-in alias of uint8, Name('A') returns "97"

	if atype.String('A') != "65" {
		t.Errorf("A ==> string(%s)", atype.String('A'))
	}

	if atype.String(byte('A')) != "65" {
		t.Errorf("byte(A) ==> string(%s)", atype.String(byte('A')))
	}

	if atype.String(atype.AByte('A')) != "A" {
		t.Errorf("atype.Abyte(A) ==> string(%s)", atype.String(atype.AByte('A')))
	}

	if atype.String([]byte{'A', 'a', 'r', 'i', 'o'}) != "Aario" {
		t.Errorf("[]byte(Aario) ==> string(%s)", atype.String([]byte{'I', 'w', 'i'}))
	}

	if atype.String("Aario") != "Aario" {
		t.Errorf("string(Aario) ==> string(%s)", atype.String("Aario"))
	}

	if atype.String(int8(100)) != "100" {
		t.Errorf("int8(100) ==> string(%s)", atype.String(int8(100)))
	}

	if atype.String(int16(100)) != "100" {
		t.Errorf("int16(100) ==> string(%s)", atype.String(int16(100)))
	}

	if atype.String(int32(100)) != "100" {
		t.Errorf("int32(100) ==> string(%s)", atype.String(int32(100)))
	}
	if atype.String(100) != "100" {
		t.Errorf("int(100) ==> string(%s)", atype.String(100))
	}

	if atype.String(int64(100)) != "100" {
		t.Errorf("int64(100) ==> string(%s)", atype.String(int64(100)))
	}

	if atype.String(float32(100.0)) != "100" {
		t.Errorf("float32(100.0) ==> string(%s)", atype.String(float32(100.0)))
	}
	if atype.String(100.0) != "100" {
		t.Errorf("float64(100.0) ==> string(%s)", atype.String(100.0))
	}

	b := 234242342342423.3
	if atype.String(b) != "234242342342423.3" {
		t.Errorf("float64(%f) ==> string(%s)", b, atype.String(b))
	}

}
