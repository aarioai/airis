package afmt

import (
	"strings"

	"github.com/aarioai/airis/pkg/types"
)

// Resize 将[]byte/[]rune slice 填充到指定长度，超过长度需要截断
func Resize[T byte | rune](s []T, pad T, length int, padHead bool) []T {
	paddingLength := length - len(s)
	if length < 0 || paddingLength == 0 {
		return s
	}
	if paddingLength > 0 {
		padding := make([]T, paddingLength)
		for i := 0; i < paddingLength; i++ {
			padding[i] = pad
		}
		if len(s) == 0 {
			return padding
		}
		if padHead {
			s = append(padding, s...)
		} else {
			s = append(s, padding...)
		}
	}
	return s[:length]
}

// PadBoth 在字符串两端添加填充字符，使其达到指定长度。返回的长度大于等于 length，除非 padString 为空
// 如果 padString 太长，无法适应 length，则会从末尾被截断。这个跟 JS padStart/padEnd 一致
// 通常使用于短字符填充，因此直接返回字符串
//
// @param trimEdge 从两边截取多余填充，如 ~_~||~_~|AARIO|~_~||~_~；否则就是从中间截取 如 |~_~||~_~AARIO~_~||~_~|
//
// @note slice 入参安全
func PadBoth[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padString := string(pad)
	paddingLength := length - len(s)

	// 如果不需要填充，直接返回原字符串
	if paddingLength <= 0 || padString == "" {
		return s
	}

	// 计算两侧需要的填充长度
	leftPadLen := paddingLength / 2
	rightPadLen := paddingLength - leftPadLen

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
// 通常使用于短字符填充，因此直接返回字符串
//
// @param trimEdge 从边缘截取多余填充，如 ~_~||~_~|AARIO；否则就是从中间截取 如 |~_~||~_~AARIO
//
// @note slice 入参安全
func PadLeft[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padString := string(pad)
	paddingLength := length - len(s)
	if paddingLength <= 0 || padString == "" {
		return s
	}
	repeatCount := paddingLength
	if len(padString) > 1 {
		repeatCount = (paddingLength / len(padString)) + 1
	}
	padding := strings.Repeat(padString, repeatCount)
	if First(trimEdge) {
		padding = padding[len(padding)-paddingLength:]
	} else {
		padding = padding[:paddingLength]
	}
	return padding + s
}

// PadRight
// 通常使用于短字符填充，因此直接返回字符串
//
// @param trimEdge 从边缘截取多余填充，如 AARIO|~_~||~_~；否则就是从中间截取 如 AARIO~_~||~_~|
//
// @note slice 入参安全
func PadRight[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padString := string(pad)
	paddingLength := length - len(s)
	if paddingLength <= 0 || padString == "" {
		return s
	}
	repeatCount := paddingLength
	if len(padString) > 1 {
		repeatCount = (paddingLength / len(padString)) + 1
	}
	padding := strings.Repeat(padString, repeatCount)
	if First(trimEdge) {
		padding = padding[:paddingLength]
	} else {
		padding = padding[len(padding)-paddingLength:]
	}

	return s + padding
}

// TrimRight 截断右侧所有cut字符。若想cut字符串，应该使用 strings/bytes.Trim/TrimFunc等
// 这里是对 byte = 0 情况的扩展（DES加密对齐使用byte(0))
//
// @note slice 入参安全
// @warn 出参可能会产生副作用，即有些情况会返回入参slice部分
func TrimRight[T byte | rune](s []T, cut T) []T {
	// var zero byte 不能使用 len(zero) 会报错，需要转为 len(string(zero)) = 1；或者直接判断是否为0
	var length int
	if cut == 0 {
		length = 1
	} else {
		length = len(s) // []rune length 是 rune 个数，而不是字节数
	}
	i := length - 1
	for ; i > -1; i-- {
		if s[i] != cut {
			break
		}
	}
	return s[:i+1]
}

func TrimLeft[T byte | rune](s []T, cut T) []T {
	// var zero byte 不能使用 len(zero)
	var length int
	if cut == 0 {
		length = 1
	} else {
		length = len(s)
	}
	i := 0
	for ; i < length; i++ {
		if s[i] != cut {
			break
		}
	}
	return s[i:]
}
func Trim[T byte | rune](s []T, cut T) []T {
	return TrimLeft(TrimRight(s, cut), cut)
}
