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

	// maybe /airis@v1.0.0 or /airis/
	if strings.Contains(file, "/airis") {
		return Caller(skip + 1)
	}

	parts := strings.Split(file, "/")
	fileLen := len(parts)
	if len(parts) > 2 {
		file = strings.Join(parts[fileLen-2:], "/")
	}
	return file + ":" + types.FormatInt(line)
}
