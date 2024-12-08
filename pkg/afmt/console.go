package afmt

import (
	"fmt"
)

const (
	// 前景色
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// 背景色
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// 格式
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Reverse   = "\033[7m"
	Reset     = "\033[0m" // 重置
)

// PrintBorder 命令行中，输出边界
func PrintBorder(msg string, styles ...string) {
	msg = PadBoth(" "+msg+" ", 80, "=", false)
	Print("\n"+msg+"\n", styles...)
}
func Println(msg string, styles ...string) {
	Print(msg+"\n", styles...)
}

// Print 命令行中，输出带颜色等样式的信息
func Print(msg string, styles ...string) {
	if len(styles) > 0 {
		for _, style := range styles {
			msg = style + msg
		}
		msg += Reset
	}
	fmt.Print(msg)
}

// 前景色打印函数
func PrintRed(format string, args ...any) {
    Print(Sprintf(format, args...), Red)
}

func PrintGreen(format string, args ...any) {
    Print(Sprintf(format, args...), Green)
}

func PrintYellow(format string, args ...any) {
    Print(Sprintf(format, args...), Yellow)
}

func PrintBlue(format string, args ...any) {
    Print(Sprintf(format, args...), Blue)
}

func PrintMagenta(format string, args ...any) {
    Print(Sprintf(format, args...), Magenta)
}

func PrintCyan(format string, args ...any) {
    Print(Sprintf(format, args...), Cyan)
}

func PrintWhite(format string, args ...any) {
    Print(Sprintf(format, args...), White)
}

// 背景色打印函数
func PrintBgBlack(format string, args ...any) {
    Print(Sprintf(format, args...), BgBlack)
}

func PrintBgRed(format string, args ...any) {
    Print(Sprintf(format, args...), BgRed)
}

func PrintBgGreen(format string, args ...any) {
    Print(Sprintf(format, args...), BgGreen)
}

func PrintBgYellow(format string, args ...any) {
    Print(Sprintf(format, args...), BgYellow)
}

func PrintBgBlue(format string, args ...any) {
    Print(Sprintf(format, args...), BgBlue)
}

func PrintBgMagenta(format string, args ...any) {
    Print(Sprintf(format, args...), BgMagenta)
}

func PrintBgCyan(format string, args ...any) {
    Print(Sprintf(format, args...), BgCyan)
}

func PrintBgWhite(format string, args ...any) {
    Print(Sprintf(format, args...), BgWhite)
}

// 样式打印函数
func PrintBold(format string, args ...any) {
    Print(Sprintf(format, args...), Bold)
}

func PrintUnderline(format string, args ...any) {
    Print(Sprintf(format, args...), Underline)
}

func PrintReverse(format string, args ...any) {
    Print(Sprintf(format, args...), Reverse)
}
