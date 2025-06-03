package alog

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"strings"
)

type ErrorLevel uint8

const (
	LevelErrAll ErrorLevel = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
	LevelAlert
	LevelEmerg
)

var (
	levelToName = map[ErrorLevel]string{
		LevelErrAll: "",
		LevelDebug:  "debug",
		LevelInfo:   "info",
		LevelNotice: "notice",
		LevelWarn:   "warn",
		LevelError:  "error",
		LevelFatal:  "fatal",
		LevelAlert:  "alert",
		LevelEmerg:  "emerg",
	}
	nameToLevel = map[string]ErrorLevel{
		"":         LevelErrAll,
		"all":      LevelErrAll,
		"debug":    LevelDebug,
		"info":     LevelInfo,
		"notice":   LevelNotice,
		"warn":     LevelWarn,
		"warning":  LevelWarn,
		"error":    LevelError,
		"fatal":    LevelFatal,
		"critical": LevelFatal,
		"alert":    LevelAlert,
		"emerg":    LevelEmerg,
	}
)

func NameToLevel(name string) ErrorLevel {
	return nameToLevel[strings.ToLower(name)] // 人为崩溃原则
}
func (lvl ErrorLevel) Name() string {
	return levelToName[lvl]
}

type LogInterface interface {
	// New 添加前缀、后缀到输出
	New(prefix string, f func(context.Context, string, ...any), suffix ...string) func(context.Context, string, ...any)

	// Debug 包含详细的开发情报的信息，通常只在调试一个程序时使用
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, format string, args ...any)

	// Info 情报信息，正常的系统消息，比如骚扰报告，带宽数据等，不需要处理。
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, args ...any)

	// Notice 不是错误情况，也不需要立即处理。
	Notice(ctx context.Context, msg string)
	Noticef(ctx context.Context, format string, args ...any)

	// Warn 警告信息，不是错误，比如系统磁盘使用了85%等。
	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, format string, args ...any)

	// Error 错误，不是非常紧急，在一定时间内修复即可。
	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, format string, args ...any)

	E(ctx context.Context, e *ae.Error, msg ...any)

	// Fatal 重要情况，如硬盘错误，备用连接丢失
	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, format string, args ...any)

	// Alert 应该被立即改正的问题，如系统数据库被破坏，ISP连接丢失。
	Alert(ctx context.Context, msg string)
	Alertf(ctx context.Context, format string, args ...any)

	// LevelEmerg 紧急情况，需要立即通知技术人员。
	Emerg(ctx context.Context, msg string)
	Emergf(ctx context.Context, format string, args ...any)

	Print(ctx context.Context, msg string)
	Printf(ctx context.Context, format string, args ...any)

	// Trace 跟踪请求链路，或性能监控
	Trace(ctx context.Context)
}
