package atype

import (
	"github.com/aarioai/airis/aa/aenum"
	"strings"
)

type AudioSrc struct {
	Provider int    `json:"provider"`
	Pattern  string `json:"pattern"` // e.g.  https://xxx/video.avi?quality=${QUALITY}
	Origin   string `json:"origin"`  // 不一定是真实的
	Path     string `json:"path"`
	//Filename  string `json:"filename"` // basename + extension   直接交path给服务端处理
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Duration int            `json:"duration"` // 时长，秒
	Jsonkey  string         `json:"jsonkey"`  // 特殊约定字段
}

func (s AudioSrc) Filename() Audio { return Audio(s.Path) }

func (s AudioSrc) Adjust(quality string) string {
	return strings.ReplaceAll(s.Pattern, "${QUALITY}", quality)
}
