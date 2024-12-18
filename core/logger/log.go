package logger

import "context"

type ErrorLevel uint8

const (
	AllError ErrorLevel = iota
	Debug
	Info
	Notice
	Warn
	Err
	Crit
	Alert
	Emerg
)

const (
	ErrorLevelKey = "aa_error_level"
)

func (lvl ErrorLevel) Name() string {
	switch lvl {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Notice:
		return "notice"
	case Warn:
		return "warn"
	case Err:
		return "err"
	case Crit:
		return "crit"
	case Alert:
		return "alert"
	case Emerg:
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

	// Info 情报信息，正常的系统消息，比如骚扰报告，带宽数据等，不需要处理。
	Info(ctx context.Context, msg string, args ...any)

	// Notice 不是错误情况，也不需要立即处理。
	Notice(ctx context.Context, msg string, args ...any)

	// Warn 警告信息，不是错误，比如系统磁盘使用了85%等。
	Warn(ctx context.Context, msg string, args ...any)

	// Error 错误，不是非常紧急，在一定时间内修复即可。
	Error(ctx context.Context, msg string, args ...any)

	// Crit 重要情况，如硬盘错误，备用连接丢失
	Crit(ctx context.Context, msg string, args ...any)

	// Alert 应该被立即改正的问题，如系统数据库被破坏，ISP连接丢失。
	Alert(ctx context.Context, msg string, args ...any)

	// Emerg 紧急情况，需要立即通知技术人员。
	Emerg(ctx context.Context, msg string, args ...any)

	Println(ctx context.Context, msg ...any)

	// Trace 跟踪请求链路，用于性能监控
	Trace(ctx context.Context)
}

func errorlevel(ctx context.Context) ErrorLevel {
	level, _ := ctx.Value(ErrorLevelKey).(ErrorLevel)
	return level
}
