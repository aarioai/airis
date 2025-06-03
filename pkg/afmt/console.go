package afmt

import (
	"fmt"
	"strings"
)

// font colors
const (
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

// background colors
const (
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// styles
const (
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Reverse   = "\033[7m"
	Reset     = "\033[0m"
)

// PrintBorder 命令行中，输出边界
func PrintBorder(msg string, styles ...string) {
	msg = PadBoth(strings.ToUpper(msg), "=", 80)
	Printc("\n"+msg+"\n", styles...)
}
func Console(msg string, styles ...string) {
	Printc(msg+"\n", styles...)
}

// Printc 命令行中，输出带颜色等样式的信息
func Printc(msg string, styles ...string) {
	fmt.Print(WithStyle(msg, styles...))
}

func WithStyle(msg string, styles ...string) string {
	if len(styles) > 0 {
		for _, style := range styles {
			msg = style + msg
		}
		msg += Reset
	}
	return msg
}

func PrintRed(msg string) {
	Printc(msg, Red)
}

func PrintfRed(format string, args ...any) {
	PrintRed(fmt.Sprintf(format, args...))
}

func PrintGreen(msg string) {
	Printc(msg, Green)
}
func PrintfGreen(format string, args ...any) {
	PrintGreen(fmt.Sprintf(format, args...))
}
func PrintYellow(msg string) {
	Printc(msg, Yellow)
}
func PrintfYellow(format string, args ...any) {
	PrintYellow(fmt.Sprintf(format, args...))
}
