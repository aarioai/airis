package atype

import "strings"

func (t Ext) Extension() Extension {
	return Extension("." + t)
}
func (t Ext) WithDot() string {
	return "." + string(t)
}

func (t Extension) Ext() Ext {
	return Ext(t[1:])
}

func NewStdFilename(s string) (StdFilename, bool) { return StdFilename(s), IsFilename(s) }
func (s StdFilename) Valid() bool {
	return IsFilename(string(s))
}

func NewFilename(s string) (Filename, bool) {
	return Filename(s), IsUnicodeFilename(s)
}
func (s Filename) Valid() bool {
	return IsUnicodeFilename(string(s))
}

func NewStdPath(s string) (StdPath, bool) { return StdPath(s), IsPath(s) }
func (s StdPath) Valid() bool {
	return IsPath(string(s))
}

func NewPath(s string) (Path, bool) {
	return Path(s), IsUnicodePath(s)
}
func (s Path) Valid() bool {
	return IsUnicodePath(string(s))
}

func NewEmail(s string) (Email, bool) {
	return Email(s), IsEmail(s)
}

func (p FilenamePattern) ReplaceAll(name string, paramType PathParamType, to string) FilenamePattern {
	if p == "" {
		return p
	}
	s := strings.ReplaceAll(string(p), "{"+name+"}", to)
	s = strings.ReplaceAll(s, "{"+name+string(paramType)+"}", to)
	return FilenamePattern(s)
}
func (p FilenamePattern) ReplaceMany(d map[string]string) FilenamePattern {
	if p == "" {
		return p
	}
	s := string(p)
	for old, to := range d {
		s = strings.ReplaceAll(s, old, to)
	}
	return FilenamePattern(s)
}
func (p FilenamePattern) Filename() Filename { return Filename(p) }

func (p PathPattern) ReplaceAll(name string, paramType PathParamType, to string) PathPattern {
	if p == "" {
		return p
	}
	s := strings.ReplaceAll(string(p), "{"+name+"}", to)
	s = strings.ReplaceAll(s, "{"+name+string(paramType)+"}", to)
	return PathPattern(s)
}
func (p PathPattern) ReplaceMany(d map[string]string) PathPattern {
	if p == "" {
		return p
	}
	s := string(p)
	for old, to := range d {
		s = strings.ReplaceAll(s, old, to)
	}
	return PathPattern(s)
}
func (p PathPattern) Path() Path { return Path(p) }

func (p UrlPattern) ReplaceAll(name string, paramType PathParamType, to string) UrlPattern {
	if p == "" {
		return p
	}
	s := strings.ReplaceAll(string(p), "{"+name+"}", to)
	s = strings.ReplaceAll(s, "{"+name+string(paramType)+"}", to)
	return UrlPattern(s)
}
func (p UrlPattern) ReplaceMany(d map[string]string) UrlPattern {
	if p == "" {
		return p
	}
	s := string(p)
	for old, to := range d {
		s = strings.ReplaceAll(s, old, to)
	}
	return UrlPattern(s)
}
func (p UrlPattern) URL() URL { return URL(p) }
