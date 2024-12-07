package atype_test

import (
	"github.com/aarioai/airis/core/atype"
	"testing"
)

func TestAtype(t *testing.T) {
	zero := atype.New("")
	if zero.String() != "" {
		t.Errorf("String() should return empty string")
	}
	b := 234242342342423.3
	s := atype.New(b).String()
	if s != "234242342342423.3" {
		t.Errorf("float64(%f) ==> string(%s)", b, s)
	}

}
func TestAtypeGet(t *testing.T) {
	arr := map[any]any{
		1:      100,
		"name": "Aario",
		"1":    "999",
		"test": map[string]any{
			"nation": "China",
			"city":   "Shenzhen",
		},
		2: map[string]string{
			"sex": "male",
		},
	}

	d := atype.New(arr)
	v, err := d.Get("name")
	if err != nil {
		t.Error("get name failed")
	}

	if v.String() != "Aario" {
		t.Errorf(`["name"] %s != Aario`, v.String())
	}

	v, err = d.Get(1)
	if err != nil {
		t.Error("get 1 failed")
	}

	i, err := v.Int()
	if i != 100 {
		t.Errorf("[1] %d != 100", i)
	}

	v, err = d.Get("1")
	if err != nil {
		t.Error(`get "1" failed`)
	}
	if v.String() != "999" {
		t.Errorf("[\"1\"] %s != 999", v.String())
	}

	v, err = d.Get("test.nation")
	if err != nil {
		t.Error("get test.nation failed")
	}
	if v.String() != "China" {
		t.Errorf("test.nation %s != China", v.String())
	}

	v, err = d.Get(2, "sex")
	if err != nil {
		t.Error("get 2.sex failed")
	}
	if v.String() != "male" {
		t.Errorf("2.sex %s != male", v.String())
	}
}
