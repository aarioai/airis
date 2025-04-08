package atype_test

import (
	"github.com/aarioai/airis/aa/atype"
	"testing"
)

func testToVersion(t *testing.T, s string, want atype.Version, semver string) {
	version := atype.ToVersion(s)
	if version != want {
		t.Errorf("ToVersion(%s) failed, want %d, got %d", s, want, version)
	}
	if version.Semver() != semver {
		t.Errorf("convert version %d to s failed  got %s", version, version.Semver())
	}
}

func testBeforeVersion(t *testing.T, a, b string) {
	x := atype.ToVersion(a)
	y := atype.ToVersion(b)
	if !x.Before(y) {
		t.Errorf("Before(%s, %s) failed", a, b)
	}
}

func testAfterVersion(t *testing.T, a, b string) {
	x := atype.ToVersion(a)
	y := atype.ToVersion(b)
	if !x.After(y) {
		t.Errorf("After(%s, %s) failed", a, b)
	}
}

func testVersionCompare(t *testing.T, a, b string, want int64) {
	v1 := atype.ToVersion(a)
	v2 := atype.ToVersion(b)
	got, _, c := v1.Compare(v2)
	if got != want {
		t.Errorf("Compare(%s => %d, %s => %d) failed, want %d, got %d", a, v1, b, v2, want, got)
	}
	if (got > 0 && c <= 0) || (got < 0 && c >= 0) {
		t.Errorf("Compare(%s, %s) failed, result %d", a, b, c)
	}
}

func testVersionTagAdd(t *testing.T, tag atype.VersionTag, add int8, want atype.VersionTag) {
	got := tag.Add(add)
	if got != want {
		t.Errorf("tag Add(%d) failed, want %d, got %d", add, want, got)
	}
}

func testVersionAdd(t *testing.T, version string, add int, tagAdd int8, want string) {
	v := atype.ToVersion(version)
	got := v.Add(add, tagAdd).Semver()
	if got != want {
		t.Errorf("version Add(%d) failed, want %s, got %s", add, want, got)
	}
}

func testVersionMinorAdd(t *testing.T, version string, minorAdd int, tagAdd int8, want string) {
	v := atype.ToVersion(version)
	got := v.AddMinor(minorAdd, tagAdd).MinorVersion()
	if got != want {
		t.Errorf("version AddMinor(%d) failed, want %s, got %s", minorAdd, want, got)
	}
}

func TestSemver(t *testing.T) {
	testToVersion(t, "v1.1.250", 1001250, "1.1.250")
	testToVersion(t, "v1.1.250a", 1001001250, "1.1.250-alpha")
	testToVersion(t, "v2.1a", 1002001000, "2.1.0-alpha")
	testToVersion(t, "v2.2b2", 2002002000, "2.2.0-beta")
	testToVersion(t, "v1.1.250.alpha", 1001001250, "1.1.250-alpha")
	testToVersion(t, "v1.1.250.alpha-250", 1001001250, "1.1.250-alpha")
	testToVersion(t, "v1.1.250.alpha-234.23", 1001001250, "1.1.250-alpha")
	testToVersion(t, "v1.1.250.241010_alpha", 1001001250, "1.1.250-alpha")

	// Test before
	testBeforeVersion(t, "v1.1.249", "v1.1.250")
	testBeforeVersion(t, "v0.0.666.alpha", "v0.0.666")
	testBeforeVersion(t, "v0.0.666.alpha", "v0.0.666.beta")
	testBeforeVersion(t, "v0.0.666.alpha-1.234", "v0.0.666.beta")
	testBeforeVersion(t, "v0.0.666.250101_alpha-1.234", "v0.0.666.beta")

	testAfterVersion(t, "v0.0.889", "v0.0.888")

	testVersionCompare(t, "v0.0.1", "v0.0.2", -1)
	testVersionCompare(t, "0.2.0", "v0.1.0", 1000)
	testVersionCompare(t, "0.1.0", "v0.1.1", -1)
	testVersionCompare(t, "v0.1.0", "v0.0.1", 999)
	testVersionCompare(t, "v0.0.1a", "v0.1.0", -999)

	testVersionTagAdd(t, atype.TagRelease, 1, atype.TagRelease)
	testVersionTagAdd(t, atype.TagAlpha, -1, atype.TagAlpha)
	testVersionTagAdd(t, atype.TagAlpha, 0, atype.TagAlpha)
	testVersionTagAdd(t, atype.TagAlpha, 1, atype.TagBeta)
	testVersionTagAdd(t, atype.TagRevision, -2, atype.TagBeta)

	testVersionAdd(t, "v0.0.1", 1, 0, "0.0.2")
	testVersionAdd(t, "v0.0.1a", 0, 1, "0.0.1-beta")
	testVersionAdd(t, "v0.0.999", 1, 1, "0.1.0")
	testVersionAdd(t, "v0.1.1-beta", 1, 1, "0.1.2-rc")

	testVersionMinorAdd(t, "v0.1b", 1, 1, "0.2rc")
	testVersionMinorAdd(t, "v0.999a", 1, 1, "1.0b")
	testVersionMinorAdd(t, "v1.1", 1, 1, "1.2")
}
