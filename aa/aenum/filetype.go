package aenum

import "github.com/aarioai/airis/pkg/types"

// Deprecated

type FileType uint16

func ParseFileType(mime string, types map[FileType][]string) (FileType, bool) {
	if mime == "" {
		return 0, false
	}
	for ft, mimes := range types {
		for _, m := range mimes {
			if m == mime {
				return ft, true
			}
		}
	}
	return 0, false
}

func (t FileType) Uint16() uint16    { return uint16(t) }
func (t FileType) String() string    { return types.FormatUint(t) }
func (t FileType) Is(t2 uint16) bool { return t.Uint16() == t2 }
func (t FileType) In(ts ...FileType) bool {
	for _, ty := range ts {
		if ty == t {
			return true
		}
	}
	return false
}
