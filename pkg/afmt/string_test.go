package afmt_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/aarioai/airis/core/atype"
	"github.com/aarioai/airis/pkg/afmt"
)

func TestResize(t *testing.T) {
	text := []byte("Hello, Aario!")
	got := afmt.Resize(text, '-', -1, false)
	if !bytes.Equal(text, got) {
		t.Errorf("Resize wrong: got %s (len:%d)", string(got), len(got))
	}
	for i := 0; i < 50; i += 5 {
		got = afmt.Resize(text, '-', i, false)
		if len(got) != i {
			t.Errorf("Resize wrong: got %s (len:%d), want len:%d", string(got), len(got), i)
		}
		if i > len(text) {
			padding := bytes.Repeat([]byte("-"), i-len(text))
			want := append(bytes.Clone(text), padding...)
			if !bytes.Equal(got, want) {
				t.Errorf("Resize wrong: got %s (len:%d), want %s", string(got), len(got), string(want))
			}
		}
	}
	want := []byte("Hello")
	got = afmt.Resize(text, '-', 5, false)
	if !bytes.Equal(got, want) {
		t.Errorf("Resize wrong: got %s (len:%d), want %s", string(got), len(got), string(want))
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

func TestTrim(t *testing.T) {
	src := []rune("~~~~^_^我LOVE你^_^~~~~~")
	want := "^_^我LOVE你^_^"
	got := afmt.Trim(src, '~')
	if want != string(got) {
		t.Errorf("TrimLeft wrong: got %s (len:%d)", string(got), len(got))
	}

	srcBytes := []byte(string(src))
	gotBytes := afmt.Trim(srcBytes, '~')
	if want != string(gotBytes) {
		t.Errorf("TrimLeft wrong: got %s (len:%d)", string(got), len(got))
	}
}
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
