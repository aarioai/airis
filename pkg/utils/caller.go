package utils

import (
	"github.com/aarioai/airis/pkg/types"
	"runtime"
	"strings"
)

func Caller(skip int) string {
	// skip Caller itself
	skip++

	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	
	parts := strings.Split(file, "/")
	fileLen := len(parts)
	var filePath string
	switch fileLen {
	case 0:
	case 1:
		filePath = parts[0]
	default:
		if parts[fileLen-2] == "alog" {
			return Caller(skip + 1)
		}
		filePath = strings.Join(parts[fileLen-2:], "/")
	}

	return filePath + ":" + types.FormatInt(line)
}
