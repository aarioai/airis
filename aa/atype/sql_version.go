package atype

import (
	"github.com/aarioai/airis/pkg/types"
	"strings"
)

// Version 版本，nginx 等比较版本，就是转为每节3位数字的整数进行比较
// Semantic Versioning https://semver.org/lang/zh-CN/
// [tag 0 release, 1 alpha, 2 beta,3 RC, 4 Revision][major 000-999][minor 000-999][build/patch 000-999]
type Version uint
type VersionTag uint8
type VersionStruct struct {
	Main  uint // Major*1000000 + Minor*1000 + Patch
	Major uint
	Minor uint
	Patch uint
	Tag   VersionTag
}

const (
	TagRelease  VersionTag = 0
	TagAlpha    VersionTag = 1
	TagBeta     VersionTag = 2
	TagRC       VersionTag = 3
	TagRevision VersionTag = 4
)
const (
	// Version multipliers for packing/unpacking
	majorMultiplier = 1000000
	minorMultiplier = 1000
	tagMultiplier   = 1000000000
)

var exceedVersionTag = int(TagRevision) + 1

// tagMap maps version tags to their numeric values
var tagMap = map[string]VersionTag{
	"a":         TagAlpha,
	"alpha":     TagAlpha,
	"b":         TagBeta,
	"beta":      TagBeta,
	"rc":        TagRC,
	"candidate": TagRC,
	"revision":  TagRevision,
}
var tagUnmap = map[VersionTag]string{
	TagAlpha:    "alpha",
	TagBeta:     "beta",
	TagRC:       "rc",
	TagRevision: "revision",
}

func NewVersionTag(t uint8) VersionTag {
	tag := VersionTag(t)
	if _, ok := tagUnmap[tag]; ok {
		return tag
	}
	return 0
}

func (t VersionTag) Name() string {
	if name, ok := tagUnmap[t]; ok {
		return name
	}
	return ""
}
func (t VersionTag) Short() string {
	switch t {
	case TagAlpha:
		return "a"
	case TagBeta:
		return "b"
	case TagRC:
		return "rc"
	case TagRevision:
		return "rev"
	}
	return ""
}

// Compare version tag to another one
// Return <0 on less than, 0 on equal to, >0 on greater than
func (t VersionTag) Compare(t2 VersionTag) int {
	if t == t2 {
		return 0
	}
	a := int(t)
	b := int(t2)
	if t == TagRelease {
		a = exceedVersionTag
	}
	if t2 == TagRelease {
		b = exceedVersionTag
	}
	return a - b
}

func (t VersionTag) Before(t2 VersionTag) bool {
	return t.Compare(t2) < 0
}

func (t VersionTag) After(t2 VersionTag) bool {
	return t.Compare(t2) > 0
}

func (t VersionTag) PackValue() uint {
	return uint(t) * tagMultiplier
}

func (t VersionTag) Add(n int8) VersionTag {
	if n == 0 {
		return t
	}
	a := int(t)
	if t == TagRelease {
		a = exceedVersionTag
	}
	a += int(n)
	if a <= 0 {
		return TagAlpha
	}
	if a >= exceedVersionTag {
		return TagRelease
	}
	return VersionTag(a)
}

// Example 2.0b2 2.0a
func splitMinorVersion(a, b string) (major, minor, patch, tag string) {
	major = a
	patch = "0"
	tagStart := false
	for _, c := range b {
		if !tagStart {
			if c >= '0' && c <= '9' {
				minor += string(c)
				continue
			}
			tagStart = true
		}
		tag += string(c)
	}
	return
}

func splitVersion(s string) (major, minor, patch, tag string) {
	if s == "" {
		return
	}
	// Remove 'v' or 'V' prefix if present, Example v0.0.xxx --> 0.0.xxx
	if s[0] == 'v' || s[0] == 'V' {
		s = s[1:]
	}
	parts := strings.Split(s, ".")
	if len(parts) == 2 {
		return splitMinorVersion(parts[0], parts[1])
	}
	if len(parts) < 3 {
		return
	}
	major = parts[0]
	minor = parts[1]
	patch = parts[2]

	// 1alpha  1a
	index := -1
	for i, c := range patch {
		if c < '0' || c > '9' {
			index = i
			break
		}
	}
	if index > -1 {
		tag = patch[index:]
		patch = patch[:index]
		return
	}

	if len(parts) > 3 {
		patch = parts[2]
		tag = parts[3]
		return
	}
	pr := strings.SplitN(parts[2], "-", 2)
	patch = pr[0]
	if len(pr) > 1 {
		tag = pr[1]
	}
	return
}

// ToVersion convert a semantic versioning number to a Version type
// Example v0.0.1-alpha.250 0.0.0.1 0.0.0.alpha 0.32a
func ToVersion(version string, withoutTag ...bool) Version {
	major, minor, patch, tag := splitVersion(version)
	if len(withoutTag) > 0 && withoutTag[0] {
		tag = ""
	}
	return PackVersion(types.ToUint(major), types.ToUint(minor), types.ToUint(patch), ParseVersionTag(tag))
}

// ParseVersionTag converts a tag string to its numeric value
// Example alpha alpha-1 alpha.250   20220830_alpha a b a1 rc.30
func ParseVersionTag(s string) VersionTag {
	if s == "" {
		return TagRelease
	}
	if tag, ok := tagMap[s]; ok {
		return tag
	}
	var cleanTag strings.Builder
	cleanTag.Grow(10)
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			cleanTag.WriteByte(byte(c))
			continue
		}
		if c >= 'A' && c <= 'Z' {
			cleanTag.WriteByte(byte(c + 32))
			continue
		}
		if cleanTag.Len() > 0 {
			break
		}
	}
	if tag, ok := tagMap[cleanTag.String()]; ok {
		return tag
	}
	return TagRelease
}

func PackVersion(major, minor, patch uint, tag VersionTag) Version {
	return Version(tag.PackValue() + major*majorMultiplier + minor*minorMultiplier + patch)
}

func (v Version) Struct() VersionStruct {
	w := uint(v)
	main := w % tagMultiplier
	tagN := w / tagMultiplier
	w = main
	patch := w % 1000
	w /= 1000
	minor := w % 1000
	w /= 1000
	major := w % 1000
	return VersionStruct{
		Main:  main,
		Major: major,
		Minor: minor,
		Patch: patch,
		Tag:   NewVersionTag(uint8(tagN)),
	}
}

// Semver returns the version in semantic versioning format
// Return <major>.<minor>.<patch>[-alpha|beta|rc|revision]
func (t VersionStruct) Semver(withoutTag ...bool) string {
	s := types.FormatUint(t.Major) + "." + types.FormatUint(t.Minor) + "." + types.FormatUint(t.Patch)
	if t.Tag != TagRelease && (len(withoutTag) == 0 || !withoutTag[0]) {
		s += "-" + t.Tag.Name()
	}
	return s
}

// MinorVersion returns the minor versioning format
// Return <major>.<minor>[a|b|rc|rev]
func (t VersionStruct) MinorVersion(withoutTag ...bool) string {
	minor := types.FormatUint(t.Minor)
	if len(withoutTag) == 0 || !withoutTag[0] {
		minor += t.Tag.Short()
	}
	return types.FormatUint(t.Major) + "." + minor
}

// Semver returns the version in semantic versioning format
// Return <major>.<minor>.<patch>[-alpha|beta|rc|revision]
func (v Version) Semver() string {
	return v.Struct().Semver()
}

func (v Version) MinorVersion() string {
	return v.Struct().MinorVersion()
}

// Compare version to another one
// Return  <0 on less than, 0 on equal to, >0 on greater than
func (v Version) Compare(v2 Version) (versionDiff int64, tagDiff int, result int) {
	if v == v2 {
		return 0, 0, 0
	}
	a := int64(v % tagMultiplier)
	tagA := NewVersionTag(uint8(v / tagMultiplier))
	b := int64(v2 % tagMultiplier)
	tagB := NewVersionTag(uint8(v2 / tagMultiplier))
	versionDiff = a - b
	tagDiff = tagA.Compare(tagB)
	result = -1
	if versionDiff > 0 || (versionDiff == 0 && tagDiff > 0) {
		result = 1
	}
	return versionDiff, tagDiff, result
}

func (v Version) Before(v2 Version) bool {
	a := v < tagMultiplier
	b := v2 < tagMultiplier
	if a && b {
		return v < v2
	}
	if a {
		return false
	}
	if b {
		return true
	}
	return v/tagMultiplier < v2/tagMultiplier
}

func (v Version) After(v2 Version) bool {
	return v2.Before(v)
}

func (v Version) Add(n int, tagAdd int8) Version {
	w := int64(v)
	main := w % tagMultiplier
	tag := NewVersionTag(uint8(w / tagMultiplier)).Add(tagAdd)

	main += int64(n)
	if main < 0 {
		main = 0
	} else if main >= tagMultiplier {
		main = tagMultiplier - 1
	}
	return Version(tag.PackValue() + uint(main))
}

func (v Version) AddMinor(minorAdd int, tagAdd int8) Version {
	w := int64(v)
	main := w % tagMultiplier
	tag := NewVersionTag(uint8(w / tagMultiplier)).Add(tagAdd)
	main += int64(minorAdd) * minorMultiplier
	if main < 0 {
		main = 0
	} else if main >= tagMultiplier {
		main = tagMultiplier - 1
	}
	return Version(tag.PackValue() + uint(main))
}
