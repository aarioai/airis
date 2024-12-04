package atype

import (
	"encoding/json"
	"html/template"
	"regexp"
	"strings"
)

type File string
type Document string
type Image string
type Video string
type Audio string
type Files struct{ NullStrings }
type Documents struct{ NullStrings }
type Images struct{ NullStrings }
type Videos struct{ NullStrings }
type Audios struct{ NullStrings }

// 文本不存在太长的；不要使用 null string，否则插入空字符串比较麻烦
// HTML 一律采用 template.HTML
type Text string // Text 65535 bytes

// 仅保留文件名，去掉目录
func trimDir(p string) string {
	if p == "" {
		return ""
	}
	i := strings.LastIndexByte(p, '/')
	if i == len(p) {
		return ""
	}
	return p[i+1:]
}

func (p File) String() string                                        { return string(p) }
func (p File) Src(filler func(string) *FileSrc) *FileSrc             { return filler(p.String()) }
func (p Document) String() string                                    { return string(p) }
func (p Document) Src(filler func(string) *DocumentSrc) *DocumentSrc { return filler(p.String()) }
func (p Image) String() string                                       { return string(p) }
func (p Image) Src(filler func(string) *ImgSrc) *ImgSrc              { return filler(p.String()) }
func (p Video) String() string                                       { return string(p) }
func (p Video) Src(filler func(string) *VideoSrc) *VideoSrc          { return filler(p.String()) }
func (p Audio) String() string                                       { return string(p) }
func (p Audio) Src(filler func(string) *AudioSrc) *AudioSrc          { return filler(p.String()) }

func NewFiles(s string) Files {
	var x Files
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToFiles(v []File) Files {
	if len(v) == 0 {
		return Files{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return Files{}
	}

	return NewFiles(string(s))
}

func (im Files) Srcs(filler func(path string) *FileSrc) []FileSrc {
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

func NewDocuments(s string) Documents {
	var x Documents
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToDocuments(v []Document) Documents {
	if len(v) == 0 {
		return Documents{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return Documents{}
	}

	return NewDocuments(string(s))
}

func (im Documents) Srcs(filler func(path string) *DocumentSrc) []DocumentSrc {
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

func NewImages(s string) Images {
	var x Images
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToImages(v []Image) Images {
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return Images{}
	}

	return NewImages(string(s))
}

func (im Images) Srcs(filler func(path string) *ImgSrc) []ImgSrc {
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

func NewVideos(s string) Videos {
	var x Videos
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToVideos(v []Video) Videos {
	if len(v) == 0 {
		return Videos{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return Videos{}
	}

	return NewVideos(string(s))
}

func NewAudios(s string) Audios {
	var x Audios
	if s != "" && strings.ToLower(s) != "null" {
		x.Scan(s)
	}
	return x
}
func ToAudios(v []Audio) Audios {
	if len(v) == 0 {
		return Audios{}
	}
	s, _ := json.Marshal(v)
	if len(s) == 0 {
		return Audios{}
	}

	return NewAudios(string(s))
}

func (im Files) Files() []File {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]File, len(imgs))
	for i, img := range imgs {
		ims[i] = File(img)
	}
	return ims
}
func (im Images) Images() []Image {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]Image, len(imgs))
	for i, img := range imgs {
		ims[i] = Image(img)
	}
	return ims
}
func (im Videos) Videos() []Video {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]Video, len(imgs))
	for i, img := range imgs {
		ims[i] = Video(img)
	}
	return ims
}
func (im Audios) Audios() []Audio {
	imgs := im.Strings()
	if len(imgs) == 0 {
		return nil
	}
	ims := make([]Audio, len(imgs))
	for i, img := range imgs {
		ims[i] = Audio(img)
	}
	return ims
}
func (t Text) String() string {
	return string(t)
}

func NewText(s string, trim bool) Text {
	if s == "" {
		return ""
	}
	if strings.Index(s, "<br") > 0 {
		s = strings.ReplaceAll(s, "<br>", "\r\n")
		s = strings.ReplaceAll(s, "<br/>", "\r\n")
	}
	if trim {
		re := regexp.MustCompile(`(^[\r\n\s\t]+)|([\r\n\s\t]$)`)
		s = re.ReplaceAllString(s, "")
		re = regexp.MustCompile(`[\s\t]*[\r\n]+[\s\t]*[\r\n]+[\s\t]*`)
		s = re.ReplaceAllString(s, `\r\n`)
	}
	return Text(s)
}

// 编码的时候
func (t Text) Html() template.HTML {
	if t == "" {
		return ""
	}
	s := t.String()
	if strings.IndexAny(s, `\r\n`) > 0 {
		s = strings.ReplaceAll(s, `\r\n`, "<br>")
		s = strings.ReplaceAll(s, `\r`, "<br>")
		s = strings.ReplaceAll(s, `\n`, "<br>")
	}
	return template.HTML(s)
}
