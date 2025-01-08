package alog

import "context"

type ErrorLevel uint8

const (
	ErrAll ErrorLevel = iota
	DEBUG
	INFO
	NOTICE
	WARN
	Err
	CRIT
	ALERT
	EMERG
)

const (
	ErrorLevelKey = "aa_error_level"
)

func (lvl ErrorLevel) Name() string {
	switch lvl {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case NOTICE:
		return "notice"
	case WARN:
		return "warn"
	case Err:
		return "err"
	case CRIT:
		return "crit"
	case ALERT:
		return "alert"
	case EMERG:
		return "emerg"
	}
	return ""
}

type LogInterface interface {
	// 添加前缀、后缀到输出
	New(prefix string, f func(context.Context, string, ...any), suffix ...string) func(context.Context, string, ...any)

	Assert(ctx context.Context, condition bool, msg string, args ...any)

	// AuthDebug 包含详细的开发情报的信息，通常只在调试一个程序时使用
	Debug(ctx context.Context, msg string, args ...any)

	// INFO 情报信息，正常的系统消息，比如骚扰报告，带宽数据等，不需要处理。
	Info(ctx context.Context, msg string, args ...any)

	// NOTICE 不是错误情况，也不需要立即处理。
	Notice(ctx context.Context, msg string, args ...any)

	// Warn 警告信息，不是错误，比如系统磁盘使用了85%等。
	Warn(ctx context.Context, msg string, args ...any)

	// Error 错误，不是非常紧急，在一定时间内修复即可。
	Error(ctx context.Context, msg string, args ...any)

	// CRIT 重要情况，如硬盘错误，备用连接丢失
	Crit(ctx context.Context, msg string, args ...any)

	// ALERT 应该被立即改正的问题，如系统数据库被破坏，ISP连接丢失。
	Alert(ctx context.Context, msg string, args ...any)

	// EMERG 紧急情况，需要立即通知技术人员。
	Emerg(ctx context.Context, msg string, args ...any)

	Println(ctx context.Context, msg ...any)

	// Trace 跟踪请求链路，用于性能监控
	Trace(ctx context.Context)
}

func errorlevel(ctx context.Context) ErrorLevel {
	level, _ := ctx.Value(ErrorLevelKey).(ErrorLevel)
	return level
}
