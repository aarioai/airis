package atype

import (
	"github.com/aarioai/airis/aa/aenum"
)

type VideoPattern struct {
}

// ImagePattern 统一使用 int，方便替换 pattern
type ImagePattern struct {
	Height    int    `json:"height"`
	Width     int    `json:"width"`
	Quality   int    `json:"quality"`
	MaxWidth  int    `json:"max_width"`
	MaxHeight int    `json:"max_height"`
	Watermark string `json:"watermark"`
}

type AudioSrc struct {
	Provider int            `json:"provider"`
	Pattern  UrlPattern     `json:"pattern"` // e.g. https://xxx/{user}/video.avi?quality={quality:int}
	Origin   URL            `json:"origin"`
	Path     Path           `json:"path"`
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Duration int            `json:"duration"` // 时长，秒
	Jsonkey  string         `json:"jsonkey"`  // 特殊约定字段
}

type DocumentSrc struct {
	Provider int            `json:"provider"`
	Path     Path           `json:"path"`
	Url      URL            `json:"url"`
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Checksum string         `json:"checksum"` // 图片、视频、音频会被压缩，checksum 无意义；这类文件不能被压缩
	Info     string         `json:"info"`     // 冗余数据
	Jsonkey  string         `json:"jsonkey"`  // 特殊约定字段
}

type FileSrc struct {
	Provider int  `json:"provider"` // 图片处理ID，如阿里云图片处理、网易云图片处理等
	Path     Path `json:"path"`     // path 可能是 filename，也可能是 带文件夹的文件名
	/*
	   不要独立出来 filename，一方面太多内容了；另一方面增加业务侧复杂度
	*/
	//StdFilename  string `json:"filename"`  // basename + extension  直接交path给服务端处理
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Checksum string         `json:"checksum"` // 图片、视频、音频会被压缩，checksum 无意义；这类文件不能被压缩
	Info     string         `json:"info"`     // 冗余数据
	Jsonkey  string         `json:"jsonkey"`  // 特殊约定字段
}

type ImgSrc struct {
	Provider      int        `json:"provider"`       // 图片处理ID，如阿里云图片处理、网易云图片处理等
	CropPattern   UrlPattern `json:"crop_pattern"`   // e.g.  https://xxx/img.jpg?width={width:int}&height={height:int}
	ResizePattern UrlPattern `json:"resize_pattern"` // e.g. https://xxx/img.jpg?maxwidth={max_width:int}
	Origin        URL        `json:"origin"`         // 不一定是真实的
	Path          Path       `json:"path"`           // path 可能是 filename，也可能是 带文件夹的文件名
	/*
	   不要独立出来 filename，一方面太多内容了；另一方面增加业务侧复杂度
	*/
	//StdFilename  string `json:"filename"`  // basename + extension  直接交path给服务端处理
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Allowed  [][2]int       `json:"allowed"` // 允许的width,height
	Jsonkey  string         `json:"jsonkey"` // 特殊约定字段
}

type VideoSrc struct {
	Provider int        `json:"provider"`
	Pattern  UrlPattern `json:"pattern"` // e.g.  https://xxx/video.avi?quality={quality:int}
	Origin   URL        `json:"origin"`  // 不一定是真实的
	Path     Path       `json:"path"`
	Preview  URL        `json:"preview"` // 一般是 gif 格式动图，所以不能缩放，直接url即可
	//StdFilename  string `json:"filename"` // basename + extension   直接交path给服务端处理
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Duration int            `json:"duration"` // 时长，秒
	Allowed  [][2]int       `json:"allowed"`  // 限定允许的width,height
	Jsonkey  string         `json:"jsonkey"`  // 特殊约定字段
}
