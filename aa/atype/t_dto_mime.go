package atype

import (
	"github.com/aarioai/airis/pkg/types"
	"regexp"
)

func (s AudioSrc) Filename() AudioPath { return AudioPath(s.Path) }

func (s AudioSrc) Adjust(quality string) URL {
	return s.Pattern.ReplaceAll("quality", "int", quality).URL()
}

// 存储在数据库里面，图片列表，为了节省空间，用数组来；具体见 atype.NullStrings or string

func (s FileSrc) Filename() FilePath { return FilePath(s.Path) }

// 存储在数据库里面，图片列表，为了节省空间，用数组来；具体见 atype.NullStrings or string

func (s ImgSrc) Filename() ImagePath { return ImagePath(s.Path) }

func (s ImgSrc) Crop(width, height int) URL {
	if s.Provider == 0 {
		return s.Origin
	}
	if width >= s.Width && height >= s.Height && s.Origin != "" {
		return s.Origin
	}
	if s.Allowed != nil {
		var matched, found bool
		var mw, mh int
		w := width
		h := height
		for _, a := range s.Allowed {
			aw := a[0]
			ah := a[1]
			if aw == width && ah == height {
				found = true
				break
			}
			if !matched {
				if aw > mw {
					mw = aw
					mh = ah
				}
				// 首先找到比缩放比例大过需求的
				if aw >= w && ah >= h {
					w = aw
					h = ah
					matched = true
				}
			} else {
				// 后面的都跟第一次匹配的比，找到最小匹配
				if aw >= width && aw <= w && ah >= height && ah <= h {
					w = aw
					h = ah
				}
			}
		}
		if !found {
			if !matched {
				width = mw
				height = mh
			} else {
				width = w
				height = h
			}
		}
	}

	sw := types.FormatInt(width)
	sh := types.FormatInt(height)
	return s.CropPattern.ReplaceMany(map[string]string{
		"{width}":      sw,
		"{width:int}":  sw,
		"{height}":     sh,
		"{height:int}": sh,
	}).URL()
}

func (s ImgSrc) Resize(maxWidth int) URL {
	if s.Provider == 0 {
		return s.Origin
	}
	if maxWidth >= s.Width && s.Origin != "" {
		return s.Origin
	}

	if s.Allowed != nil {
		var matched, found bool
		var mw int
		w := maxWidth
		for _, a := range s.Allowed {
			aw := a[0]
			if aw == maxWidth {
				found = true
				break
			}
			if !matched {
				if aw > mw {
					mw = aw
				}
				// 首先找到比缩放比例大过需求的
				if aw >= w {
					w = aw
					matched = true
				}
			} else {
				// 后面的都跟第一次匹配的比，找到最小匹配
				if aw >= maxWidth && aw <= w {
					w = aw
				}
			}
		}
		if !found {
			if !matched {
				maxWidth = mw
			} else {
				maxWidth = w
			}
		}
	}
	sw := types.FormatInt(maxWidth)
	return s.ResizePattern.ReplaceAll("max_width", ":int", sw).URL()
}

func (s VideoSrc) Filename() VideoPath { return VideoPath(s.Path) }
func (s VideoSrc) Adjust(quality string) URL {
	return s.Pattern.ReplaceAll("quality", ":int", quality).URL()
}

func ImageFill(width, height int) ImagePattern {
	return ImagePattern{Width: width, Height: height}
}
func ImageFitWidth(maxWidth int) ImagePattern {
	return ImagePattern{MaxWidth: maxWidth}
}
func ToImagePattern(tag string) ImagePattern {
	reg, _ := regexp.Compile(`([a-z]+)(\d+)`)
	matches := reg.FindAllStringSubmatch(tag, -1)
	var p ImagePattern
	for _, match := range matches {
		v, _ := types.Atoi(match[2])
		/**
		 * w width, h height, q quanlity, v max width, g max height
		 *    	img.width <= v ,   img.width = w  两者区别
		 * xN  有意义，对于不定尺寸的白名单，自动化方案是：先获取 x1 的尺寸，然后 xN ，之后把 source 裁剪
		 */
		t := match[1]
		switch t {
		case "h":
			p.Height = v
		case "w":
			p.Width = v
		case "g":
			p.MaxHeight = v
		case "v":
			p.MaxWidth = v
		case "q":
			p.Quality = v
		case "k":
			p.Watermark = match[2]
		}
	}
	return p
}

func ToVideoPattern(tag string) VideoPattern {
	return VideoPattern{}
}
