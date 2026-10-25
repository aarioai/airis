package atype

import (
	"html/template"
	"regexp"
	"strings"
)

var (
	textEdgeSpace  = regexp.MustCompile(`(^[\r\n\s\t]+)|([\r\n\s\t]$)`)
	textBlankLines = regexp.MustCompile(`[\s\t]*[\r\n]+[\s\t]*[\r\n]+[\s\t]*`)
)

func (t Text) String() string {
	return string(t)
}

func NewText(s string, trim bool) Text {
	if s == "" {
		return ""
	}
	if strings.Contains(s, "<br") {
		s = strings.ReplaceAll(s, "<br>", "\r\n")
		s = strings.ReplaceAll(s, "<br/>", "\r\n")
	}
	if trim {
		s = textEdgeSpace.ReplaceAllString(s, "")
		s = textBlankLines.ReplaceAllString(s, `\r\n`)
	}
	return Text(s)
}


func (t Text) Html() template.HTML {
	if t == "" {
		return ""
	}
	s := t.String()
	if strings.Contains(s, `\r`) || strings.Contains(s, `\n`) {
		s = strings.ReplaceAll(s, `\r\n`, "<br>")
		s = strings.ReplaceAll(s, `\r`, "<br>")
		s = strings.ReplaceAll(s, `\n`, "<br>")
	}
	return template.HTML(s)
}
