package atype

import (
	"html/template"
	"regexp"

	"github.com/aarioai/airis/pkg/types"
)

var imagePatternTag = regexp.MustCompile(`([a-z]+)(\d+)`)

func (s AudioSrcOld) Filename() AudioPath { return AudioPath(s.Path) }

func (s AudioSrcOld) Adjust(quality string) template.URL {
	return s.Pattern.ReplaceAll("quality", "int", quality).URL()
}

// Stored in the database. Image lists use arrays to save space; see atype.NullStrings or string.

func (s FileSrcOld) Filename() FilePath { return FilePath(s.Path) }

// Stored in the database. Image lists use arrays to save space; see atype.NullStrings or string.

func (s ImgSrcOld) Filename() ImagePath { return ImagePath(s.Path) }

func (s ImgSrcOld) Crop(width, height int) template.URL {
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
				// First find a size whose scale is larger than required
				if aw >= w && ah >= h {
					w = aw
					h = ah
					matched = true
				}
			} else {
				// Compare the rest with the first match and pick the smallest match
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

func (s ImgSrcOld) Resize(maxWidth int) template.URL {
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
				// First find a size whose scale is larger than required
				if aw >= w {
					w = aw
					matched = true
				}
			} else {
				// Compare the rest with the first match and pick the smallest match
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

func (s VideoSrcOld) Filename() VideoPath { return VideoPath(s.Path) }
func (s VideoSrcOld) Adjust(quality string) template.URL {
	return s.Pattern.ReplaceAll("quality", ":int", quality).URL()
}

func ImageFill(width, height int) ImagePattern {
	return ImagePattern{Width: width, Height: height}
}
func ImageFitWidth(maxWidth int) ImagePattern {
	return ImagePattern{MaxWidth: maxWidth}
}
func ToImagePattern(tag string) ImagePattern {
	matches := imagePatternTag.FindAllStringSubmatch(tag, -1)
	var p ImagePattern
	for _, match := range matches {
		v, _ := types.Atoi(match[2])
		/**
		 * w width, h height, q quanlity, v max width, g max height
		 *    	img.width <= v ,   img.width = w  difference between the two
		 * xN is meaningful. For an allowlist of variable sizes, the automated scheme is: obtain the x1 size first, then xN, then crop the source
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
