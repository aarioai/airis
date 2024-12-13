package afmt

import (
	"github.com/aarioai/airis/pkg/types"
	"strings"
)

// PadBoth 在字符串两端添加填充字符，使其达到指定长度
// 如果 padString 太长，无法适应 length，则会从末尾被截断。这个跟 JS padStart/padEnd 一致
// 返回的长度大于等于 length，除非 padString 为空
// @param trimEdge 从两边截取多余填充，如 ~_~||~_~|AARIO|~_~||~_~；否则就是从中间截取 如 |~_~||~_~AARIO~_~||~_~|
func PadBoth[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padString := string(pad)
	padLen := length - len(s)

	// 如果不需要填充，直接返回原字符串
	if padLen <= 0 || padString == "" {
		return s
	}

	// 计算两侧需要的填充长度
	leftPadLen := padLen / 2
	rightPadLen := padLen - leftPadLen

	// 生成填充字符串
	leftRepeatCount := leftPadLen
	rightRepeatCount := rightPadLen
	psLen := len(padString)
	if psLen > 1 {
		leftRepeatCount = (leftPadLen / psLen) + 1
		rightRepeatCount = (rightPadLen / psLen) + 1
	}
	leftPadding := strings.Repeat(padString, leftRepeatCount)
	rightPadding := strings.Repeat(padString, rightRepeatCount)

	// 从两边截取多余填充，如 ~_~||~_~|AARIO|~_~||~_~
	if First(trimEdge) {
		leftPadding = leftPadding[len(leftPadding)-leftPadLen:]
		rightPadding = rightPadding[:rightPadLen]
	} else {
		// 从中间截取多余填充 如 |~_~||~_~AARIO~_~||~_~|
		leftPadding = leftPadding[:leftPadLen]
		rightPadding = rightPadding[len(rightPadding)-rightPadLen:]
	}

	return leftPadding + s + rightPadding
}

// PadLeft
// fmt.Sprintf("%04d", 10)    fmt.Sprintf("%-04d", 3)  只能：左右填充空格，或左边填充0
// @param trimEdge 从边缘截取多余填充，如 ~_~||~_~|AARIO；否则就是从中间截取 如 |~_~||~_~AARIO
func PadLeft[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padString := string(pad)
	padLen := length - len(s)
	if padLen <= 0 || padString == "" {
		return s
	}
	repeatCount := padLen
	if len(padString) > 1 {
		repeatCount = (padLen / len(padString)) + 1
	}
	padding := strings.Repeat(padString, repeatCount)
	if First(trimEdge) {
		padding = padding[len(padding)-padLen:]
	} else {
		padding = padding[:padLen]
	}
	return padding + s
}

// PadRight
// @param trimEdge 从边缘截取多余填充，如 AARIO|~_~||~_~；否则就是从中间截取 如 AARIO~_~||~_~|
func PadRight[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padString := string(pad)
	padLen := length - len(s)
	if padLen <= 0 || padString == "" {
		return s
	}
	repeatCount := padLen
	if len(padString) > 1 {
		repeatCount = (padLen / len(padString)) + 1
	}
	padding := strings.Repeat(padString, repeatCount)
	if First(trimEdge) {
		padding = padding[:padLen]
	} else {
		padding = padding[len(padding)-padLen:]
	}

	return s + padding
}
