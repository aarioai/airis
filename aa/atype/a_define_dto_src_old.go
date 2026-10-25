package atype

import (
	"github.com/aarioai/airis/aa/aenum"
)

type AudioSrcOld struct {
	Provider int            `json:"provider"`
	Pattern  UrlPattern     `json:"pattern"` // e.g. https://xxx/{user}/video.avi?quality={quality:int}
	Origin   URL            `json:"origin"`
	Path     Path           `json:"path"`
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Duration int            `json:"duration"` // duration, in seconds
	Jsonkey  string         `json:"jsonkey"`  // specially agreed field
}

type DocumentSrcOld struct {
	Provider int            `json:"provider"`
	Path     Path           `json:"path"`
	Url      URL            `json:"url"`
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Checksum string         `json:"checksum"` // images, video, and audio are compressed, so checksum is meaningless; such files must not be compressed
	Info     string         `json:"info"`     // redundant data
	Jsonkey  string         `json:"jsonkey"`  // specially agreed field
}

type FileSrcOld struct {
	Provider int  `json:"provider"` // image-processing ID, e.g. Alibaba Cloud or NetEase Cloud image processing
	Path     Path `json:"path"`     // path may be a filename, or a filename that includes directories
	/*
	   Do not split filename out. There is too much content, and it adds complexity on the business side.
	*/
	//StdFilename  string `json:"filename"`  // basename + extension; hand path to the server directly
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Checksum string         `json:"checksum"` // images, video, and audio are compressed, so checksum is meaningless; such files must not be compressed
	Info     string         `json:"info"`     // redundant data
	Jsonkey  string         `json:"jsonkey"`  // specially agreed field
}

type ImgSrcOld struct {
	Provider      int        `json:"provider"`       // image-processing ID, e.g. Alibaba Cloud or NetEase Cloud image processing
	CropPattern   UrlPattern `json:"crop_pattern"`   // e.g.  https://xxx/img.jpg?width={width:int}&height={height:int}
	ResizePattern UrlPattern `json:"resize_pattern"` // e.g. https://xxx/img.jpg?maxwidth={max_width:int}
	Origin        URL        `json:"origin"`         // not necessarily the real origin
	Path          Path       `json:"path"`           // path may be a filename, or a filename that includes directories
	/*
	   Do not split filename out. There is too much content, and it adds complexity on the business side.
	*/
	//StdFilename  string `json:"filename"`  // basename + extension; hand path to the server directly
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Allowed  [][2]int       `json:"allowed"` // allowed width,height
	Jsonkey  string         `json:"jsonkey"` // specially agreed field
}

type VideoSrcOld struct {
	Provider int        `json:"provider"`
	Pattern  UrlPattern `json:"pattern"` // e.g.  https://xxx/video.avi?quality={quality:int}
	Origin   URL        `json:"origin"`  // not necessarily the real origin
	Path     Path       `json:"path"`
	Preview  URL        `json:"preview"` // usually a GIF animation, so it cannot be scaled; use the URL directly
	//StdFilename  string `json:"filename"` // basename + extension; hand path to the server directly
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Duration int            `json:"duration"` // duration, in seconds
	Allowed  [][2]int       `json:"allowed"`  // restricted set of allowed width,height
	Jsonkey  string         `json:"jsonkey"`  // specially agreed field
}
