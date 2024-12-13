package afmt_test

import (
	"bytes"
	"fmt"
	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/afmt"
	"math/rand/v2"
	"strings"
	"testing"
)

func testPadBoth(t *testing.T, title []byte, pad string, length int, startFromEnd bool) {
	titleClone := bytes.Clone(title)
	cmd := fmt.Sprintf("PadBoth(%s, %d, %s, %v)", string(title), length, pad, startFromEnd)
	s := afmt.PadBoth(title, pad, length, startFromEnd)
	if len(title) >= length && s != string(title) {
		t.Errorf("%s wrong: %s", cmd, s)
		return
	}
	if !bytes.Equal(title, titleClone) {
		t.Error("PadBoth " + afmt.ErrmsgSideEffect(title))
	}
	if len(s) != length {
		t.Errorf("%s len(%s)=%d; want %d, ", cmd, s, len(s), length)
	}

	if pad == " " || pad == "0" {
		leftLen := length / 2
		if length > leftLen+leftLen && !startFromEnd {
			leftLen++
		}
		ipad := strings.Trim(pad, " ")
		patternRight := "%-" + ipad + atype.FormatInt(length) + "s"
		leftPadded := fmt.Sprintf("%"+ipad+"*s", leftLen, s) // e.g. fmt.Sprintf("%0*s", 10,2)
		testS := fmt.Sprintf(patternRight, leftPadded)       // e.g. fmt.Sprintf("-010s", 2)
		if testS != s {
			t.Errorf("%s wrong: %s", cmd, testS)
		}
	}
}
func TestPadBoth(t *testing.T) {
	testPadBoth(t, []byte("Hello Aario"), " ", 20, false)
	testPadBoth(t, []byte("Hello Aario"), "0", 80, false)
	testPadBoth(t, []byte("Hello Aario"), "._-", 80, false)

	// 测试 trimEdge=false
	want := "|~_~||~_~ Hello Aario ~_~||~_~|"
	got := afmt.PadBoth(" Hello Aario ", "|~_~|", len(want))
	if got != want {
		t.Errorf("PadBoth got: %s; want: %s", got, want)
	}

	// 测试 trimEdge=true
	want = "~_~||~_~| Hello Aario |~_~||~_~"
	got = afmt.PadBoth(" Hello Aario ", "|~_~|", len(want), true)
	if got != want {
		t.Errorf("PadBoth trim edge got: %s; want: %s", got, want)
	}
}
func TestPadLeft(t *testing.T) {
	base := []byte("Hello Aario")
	baseClone := bytes.Clone(base)
	want := fmt.Sprintf("%020s", base)
	got := afmt.PadLeft(base, '0', 20)
	if len(got) != 20 || got != want {
		t.Errorf("PadLeft wrong: %s (len:%d), want %s", got, len(got), want)
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadLeft " + afmt.ErrmsgSideEffect(base))
	}

	// 测试 trimEdge=false
	want = "|~_~||~_~Hello Aario"
	got = afmt.PadLeft(base, "|~_~|", len(want))
	if got != want {
		t.Errorf("PadLeft wrong: %s (len:%d), want %s", got, len(got), want)
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadLeft " + afmt.ErrmsgSideEffect(base))
	}

	// 测试 trimEdge=true
	want = "~_~||~_~|Hello Aario"
	got = afmt.PadLeft(base, "|~_~|", len(want), true)
	if got != want {
		t.Errorf("PadLeft trim edge wrong: %s (len:%d), want %s", got, len(got), want)
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadLeft " + afmt.ErrmsgSideEffect(base))
	}
}
func TestPaRight(t *testing.T) {
	base := []byte("Hello Aario")
	baseClone := bytes.Clone(base)
	want := fmt.Sprintf("%-20s", base) // fmt.Sprintf 只能：左右填充空格，或左边填充0
	got := afmt.PadRight(base, ' ', 20)
	if len(got) != 20 || got != want {
		t.Errorf("PadRight wrong: %s (len:%d), want %s", got, len(got), want)
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadRight " + afmt.ErrmsgSideEffect(base))
	}

	// 测试 trimEdge=false
	want = "Hello Aario~_~||~_~|"
	got = afmt.PadRight(base, "|~_~|", len(want))
	if got != want {
		t.Errorf("PadRight wrong: %s (len:%d), want %s", got, len(got), want)
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadRight " + afmt.ErrmsgSideEffect(base))
	}

	// 测试 trimEdge=true
	want = "Hello Aario|~_~||~_~"
	got = afmt.PadRight(base, "|~_~|", len(want), true)
	if got != want {
		t.Errorf("PadRight trim edge wrong: %s (len:%d), want %s", got, len(got), want)
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadRight " + afmt.ErrmsgSideEffect(base))
	}
}

func TestPadBlock(t *testing.T) {
	base := make([]byte, 0)
	n := rand.IntN(256)
	for i := 0; i < n; i++ {
		char := rand.IntN('Z'-'A') + 'A'
		base = append(base, byte(char))
	}
	baseClone := bytes.Clone(base)
	got := afmt.PadBlock(base, '=', 8)
	if len(got)%8 != 0 {
		t.Errorf("PadBlock wrong: %s (len:%d), want %d", got, len(got), len(base))
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadBlock " + afmt.ErrmsgSideEffect(base))
	}

	// 测试分隔符
	got = afmt.PadBlock(base, '@', 8, '\n')
	if len(got)%9 != 0 {
		t.Errorf("PadBlock wrong: %s (len:%d), want %d", got, len(got), len(base))
	}
	if !bytes.Equal(baseClone, base) {
		t.Error("PadBlock " + afmt.ErrmsgSideEffect(base))
	}
}
