package atype

import (
	"encoding/json"
	"strings"
)

// 文本不存在太长的；不要使用 null string，否则插入空字符串比较麻烦
// HTML 一律采用 template.HTML

func (p FilePath) String() string                                        { return string(p) }
func (p FilePath) Src(filler func(string) *FileSrc) *FileSrc             { return filler(p.String()) }
func (p DocumentPath) String() string                                    { return string(p) }
func (p DocumentPath) Src(filler func(string) *DocumentSrc) *DocumentSrc { return filler(p.String()) }
func (p ImagePath) String() string                                       { return string(p) }
func (p ImagePath) Src(filler func(string) *ImgSrc) *ImgSrc              { return filler(p.String()) }
func (p VideoPath) String() string                                       { return string(p) }
func (p VideoPath) Src(filler func(string) *VideoSrc) *VideoSrc          { return filler(p.String()) }
func (p AudioPath) String() string                                       { return string(p) }
func (p AudioPath) Src(filler func(string) *AudioSrc) *AudioSrc          { return filler(p.String()) }

func NewFiles(s string) FilePaths {
	var x FilePaths
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToFiles(v []FilePath) FilePaths {
	if len(v) == 0 {
		return FilePaths{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return FilePaths{}
	}

	return NewFiles(string(s))
}

func (im FilePaths) Srcs(filler func(path string) *FileSrc) []FileSrc {
	if !im.Valid || im.String == "" {
		return nil
	}
	ims := im.Strings()
	srcs := make([]FileSrc, 0, len(ims))
	for _, im := range ims {
		if im != "" {
			if fi := filler(im); fi != nil {
				srcs = append(srcs, *fi)
			}
		}
	}
	return srcs
}

func NewDocuments(s string) DocumentPaths {
	var x DocumentPaths
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToDocuments(v []DocumentPath) DocumentPaths {
	if len(v) == 0 {
		return DocumentPaths{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return DocumentPaths{}
	}

	return NewDocuments(string(s))
}

func (im DocumentPaths) Srcs(filler func(path string) *DocumentSrc) []DocumentSrc {
	if !im.Valid || im.String == "" {
		return nil
	}
	ims := im.Strings()
	srcs := make([]DocumentSrc, 0, len(ims))
	for _, im := range ims {
		if im != "" {
			if fi := filler(im); fi != nil {
				srcs = append(srcs, *fi)
			}
		}
	}
	return srcs
}

func NewImages(s string) ImagePaths {
	var x ImagePaths
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToImages(v []ImagePath) ImagePaths {
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return ImagePaths{}
	}

	return NewImages(string(s))
}

func (im ImagePaths) Srcs(filler func(path string) *ImgSrc) []ImgSrc {
	if !im.Valid || im.String == "" {
		return nil
	}
	ims := im.Strings()
	srcs := make([]ImgSrc, 0, len(ims))
	for _, im := range ims {
		if im != "" {
			if fi := filler(im); fi != nil {
				srcs = append(srcs, *fi)
			}
		}
	}
	return srcs
}

func NewVideos(s string) VideoPaths {
	var x VideoPaths
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToVideos(v []VideoPath) VideoPaths {
	if len(v) == 0 {
		return VideoPaths{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return VideoPaths{}
	}

	return NewVideos(string(s))
}

func NewAudios(s string) AudioPaths {
	var x AudioPaths
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToAudios(v []AudioPath) AudioPaths {
	if len(v) == 0 {
		return AudioPaths{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return AudioPaths{}
	}

	return NewAudios(string(s))
}

func (im FilePaths) Files() []FilePath {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]FilePath, len(imgs))
	for i, img := range imgs {
		ims[i] = FilePath(img)
	}
	return ims
}
func (im ImagePaths) Images() []ImagePath {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]ImagePath, len(imgs))
	for i, img := range imgs {
		ims[i] = ImagePath(img)
	}
	return ims
}
func (im VideoPaths) Videos() []VideoPath {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]VideoPath, len(imgs))
	for i, img := range imgs {
		ims[i] = VideoPath(img)
	}
	return ims
}
func (im AudioPaths) Audios() []AudioPath {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]AudioPath, len(imgs))
	for i, img := range imgs {
		ims[i] = AudioPath(img)
	}
	return ims
}
