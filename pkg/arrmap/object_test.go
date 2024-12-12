package arrmap_test

import (
	"slices"
	"testing"

	"github.com/aarioai/airis/pkg/arrmap"
)

func TestSortedKeys(t *testing.T) {
	m := map[string]any{
		"b": 2,
		"a": 1,
		"c": 3,
	}
	want := []string{"a", "b", "c"}
	got := arrmap.SortedKeys(m)
	// 顺序必须一致
	if !slices.Equal(got, want) {
		t.Errorf("SortedKeys(%v) = %v; want %v", m, got, want)
	}
}

func TestJoinKeysStringMap(t *testing.T) {
	m := map[string]any{
		"b": 2,
		"a": 1,
		"c": 3,
	}
	// 不对key排序，则无法判断正确的顺序
	want := "a,b,c"
	got := arrmap.JoinKeys(m, ",", true)
	if got != want {
		t.Errorf("JoinKeys() = %v, want %v", got, want)
	}
}
func TestJoinKeysRuneMap(t *testing.T) {
	m := map[rune]any{
		'b': 2,
		'a': 1,
		'c': 3,
	}
	want := "a,b,c"
	got := arrmap.JoinKeys(m, ",", true)
	if got != want {
		t.Errorf("JoinKeys() = %v, want %v", got, want)
	}
}

func TestJoinKeysByteMap(t *testing.T) {
	m := map[byte]any{
		'b': 2,
		'a': 1,
		'c': 3,
	}
	want := "a,b,c"
	got := arrmap.JoinKeys(m, ",", true)
	if got != want {
		t.Errorf("JoinKeys() = %v, want %v", got, want)
	}
}
