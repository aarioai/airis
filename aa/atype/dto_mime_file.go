package atype

import "github.com/aarioai/airis/aa/aenum"

// 存储在数据库里面，图片列表，为了节省空间，用数组来；具体见 atype.NullStrings or string
type FileSrc struct {
	Provider int    `json:"provider"` // 图片处理ID，如阿里云图片处理、网易云图片处理等
	Path     string `json:"path"`     // path 可能是 filename，也可能是 带文件夹的文件名
	/*
	   不要独立出来 filename，一方面太多内容了；另一方面增加业务侧复杂度
	*/
	//Filename  string `json:"filename"`  // basename + extension  直接交path给服务端处理
	Filetype aenum.FileType `json:"filetype"` // aenum.Filetype.Int8()
	Size     int            `json:"size"`     // atype.Uint24.Int8()
	Checksum string         `json:"checksum"` // 图片、视频、音频会被压缩，checksum 无意义；这类文件不能被压缩
	Info     string         `json:"info"`     // 冗余数据
	Jsonkey  string         `json:"jsonkey"`  // 特殊约定字段
}

func (s FileSrc) Filename() File { return File(s.Path) }
