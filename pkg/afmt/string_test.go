package afmt_test

import (
	"fmt"
	"github.com/aarioai/airis/pkg/afmt"
	"testing"
)

func testPadBoth(t *testing.T, title string, length int, pad string, startFromEnd bool) {
	cmd := fmt.Sprintf("PadBoth(%s, %d, %s, %v)", title, length, pad, startFromEnd)
	s := afmt.PadBoth(title, length, pad, startFromEnd)
	if len(title) >= length && s != title {
		t.Errorf("%s wrong: %s", cmd, s)
		return
	}
	if len(s) != length {
		t.Errorf("%s len(%s)=%d; want %d, ", cmd, s, len(s), length)
	}
}
func TestPadBoth(t *testing.T) {
	testPadBoth(t, "Hello World", 80, "=", false)
	testPadBoth(t, "Hello World", 80, "._-", false)
}
