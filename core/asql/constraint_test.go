package asql_test

import (
	"github.com/aarioai/airis/pkg/asql"
	"testing"
)

type stru struct {
	Name string `name:"name"`
	Age  int    `name:"age"`
}

func TestAnd(t *testing.T) {
	u := stru{
		Name: "Aario",
		Age:  18,
	}
	s := asql.And(u, "name", "age")
	s1 := "`age`=\"18\" AND `name`=\"Aario\""
	s2 := "`name`=\"Aario\" AND `age`=\"18\""
	if s != s1 && s != s2 {
		t.Errorf("asql.And(u, true, ...) == %s", s)
	}

	s = asql.And(u, "Name", "Age")
	if s != s1 && s != s2 {
		t.Errorf("asql.And(u, false, ...) == %s", s)
	}
}

func TestOr(t *testing.T) {
	u := stru{
		Name: "Aario",
		Age:  18,
	}
	s := asql.Or(u, "name", "age")
	s1 := "`age`=\"18\" OR `name`=\"Aario\""
	s2 := "`name`=\"Aario\" OR `age`=\"18\""
	if s != s1 && s != s2 {
		t.Errorf("asql.Or(u, true, ...) == %s", s)
	}

	s = asql.Or(u, "Name", "Age")
	if s != s1 && s != s2 {
		t.Errorf("asql.Or(u, false, ...) == %s", s)
	}
}

func TestAndWithWhere(t *testing.T) {
	u := stru{
		Name: "Aario",
		Age:  18,
	}
	s := asql.AndWithWhere(u, "name", "age")
	s1 := " WHERE `age`=\"18\" AND `name`=\"Aario\" "
	s2 := " WHERE `name`=\"Aario\" AND `age`=\"18\" "
	if s != s1 && s != s2 {
		t.Errorf("asql.Or(u, true, ...) == %s", s)
	}

	s = asql.AndWithWhere(u, "Name", "Age")
	if s != s1 && s != s2 {
		t.Errorf("asql.Or(u, false, ...) == %s", s)
	}
}

func TestOrWithWhere(t *testing.T) {
	u := stru{
		Name: "Aario",
		Age:  18,
	}
	s := asql.OrWithWhere(u, "name", "age")
	s1 := " WHERE `age`=\"18\" OR `name`=\"Aario\" "
	s2 := " WHERE `name`=\"Aario\" OR `age`=\"18\" "
	if s != s1 && s != s2 {
		t.Errorf("asql.Or(u, true, ...) == %s", s)
	}

	s = asql.OrWithWhere(u, "Name", "Age")
	if s != s1 && s != s2 {
		t.Errorf("asql.Or(u, false, ...) == %s", s)
	}
}
