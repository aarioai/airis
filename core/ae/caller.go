package ae

import (
	"github.com/aarioai/airis/pkg/types"
	"runtime"
	"strings"
)

// CallerMsg 获取错误信息和调用栈
func CallerMsg(errmsg string, skip int) (string, string) {
	caller := Caller(skip)
	// 特殊处理 context canceled 错误，获取更完整的调用链
	if errmsg == "context canceled" {
		caller2 := Caller(skip + 1)
		if caller2 != caller {
			caller = caller2 + "->" + caller
		}
	}
	return errmsg, caller
}

// Caller 获取调用栈信息
// skip 表示要跳过的调用帧数
func Caller(skip int) string {
	// 跳过 Caller 自身
	skip++

	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	// 处理文件路径
	parts := strings.Split(file, "/")
	fileLen := len(parts)
	var filePath string
	switch {
	case fileLen <= 1:
		filePath = parts[0]
	case fileLen == 2:
		filePath = strings.Join(parts[0:2], "/")
	default:
		filePath = strings.Join(parts[fileLen-3:], "/")
	}

	// 检查是否为框架代码
	for _, part := range parts {
		if strings.HasPrefix(strings.ToLower(part), "!aa!go@") {
			return Caller(skip) // 递归调用获取业务代码位置
		}
	}

	// 构建调用位置信息
	var builder strings.Builder
	builder.WriteString(filePath)
	builder.WriteString(":")
	builder.WriteString(types.FormatInt(line))

	return builder.String()
}
