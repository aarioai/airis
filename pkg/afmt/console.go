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

// PrintcBorder 命令行中，输出边界
func PrintcBorder(msg string, styles ...string) {
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

// 前景色打印函数
func PrintcRed(format string, args ...any) {
	Printc(Sprintf(format, args...), Red)
}

func PrintcGreen(format string, args ...any) {
	Printc(Sprintf(format, args...), Green)
}

func PrintcYellow(format string, args ...any) {
	Printc(Sprintf(format, args...), Yellow)
}

func PrintcBlue(format string, args ...any) {
	Printc(Sprintf(format, args...), Blue)
}

func PrintcMagenta(format string, args ...any) {
	Printc(Sprintf(format, args...), Magenta)
}

func PrintcCyan(format string, args ...any) {
	Printc(Sprintf(format, args...), Cyan)
}

func PrintcWhite(format string, args ...any) {
	Printc(Sprintf(format, args...), White)
}

// 背景色打印函数
func PrintcBgBlack(format string, args ...any) {
	Printc(Sprintf(format, args...), BgBlack)
}

func PrintcBgRed(format string, args ...any) {
	Printc(Sprintf(format, args...), BgRed)
}

func PrintcBgGreen(format string, args ...any) {
	Printc(Sprintf(format, args...), BgGreen)
}

func PrintcBgYellow(format string, args ...any) {
	Printc(Sprintf(format, args...), BgYellow)
}

func PrintcBgBlue(format string, args ...any) {
	Printc(Sprintf(format, args...), BgBlue)
}

func PrintcBgMagenta(format string, args ...any) {
	Printc(Sprintf(format, args...), BgMagenta)
}

func PrintcBgCyan(format string, args ...any) {
	Printc(Sprintf(format, args...), BgCyan)
}

func PrintcBgWhite(format string, args ...any) {
	Printc(Sprintf(format, args...), BgWhite)
}

// 样式打印函数
func PrintcBold(format string, args ...any) {
	Printc(Sprintf(format, args...), Bold)
}

func PrintcUnderline(format string, args ...any) {
	Printc(Sprintf(format, args...), Underline)
}

func PrintcReverse(format string, args ...any) {
	Printc(Sprintf(format, args...), Reverse)
}
