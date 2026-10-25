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
	return nameToLevel[strings.ToLower(name)] // fail-fast: an unknown name is a programmer error
}
func (lvl ErrorLevel) Name() string {
	return levelToName[lvl]
}

type LogInterface interface {
	// New adds a prefix and an optional suffix to the output.
	New(prefix string, f func(context.Context, string, ...any), suffix ...string) func(context.Context, string, ...any)

	// Debug carries detailed development information, usually only while debugging a program.
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, format string, args ...any)

	// Info is informational: normal system messages, such as reports gathered for monitoring and bandwidth data, that need no action.
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, args ...any)

	// Notice is not an error and does not need immediate handling.
	Notice(ctx context.Context, msg string)
	Noticef(ctx context.Context, format string, args ...any)

	// Warn is a warning, not an error, for example when a disk is 85% full.
	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, format string, args ...any)

	// Error is an error that is not urgent and can be fixed within some time.
	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, format string, args ...any)

	E(ctx context.Context, e *ae.Error, msg ...any)

	// Fatal is a critical condition, such as a disk error or loss of a backup link.
	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, format string, args ...any)

	// Alert is a problem that must be corrected immediately, such as a corrupted system database or a lost ISP link.
	Alert(ctx context.Context, msg string)
	Alertf(ctx context.Context, format string, args ...any)

	// LevelEmerg is an emergency that requires notifying technicians immediately.
	Emerg(ctx context.Context, msg string)
	Emergf(ctx context.Context, format string, args ...any)

	Print(ctx context.Context, msg string)
	Printf(ctx context.Context, format string, args ...any)

	// Trace follows a request chain, or monitors performance.
	Trace(ctx context.Context)
}
