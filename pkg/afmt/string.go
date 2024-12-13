package afmt

import (
	"bytes"
	"github.com/aarioai/airis/pkg/types"
	"strings"
)

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

// PadBlock 按照 blockSize 对齐，不足部分填充。被广泛使用与 base64/DES 等编码
// 一般用于需要加密或编码的文本填充，因此返回bytes
//
// @note slice 入参安全
// @warn 出参可能会产生副作用，即有些情况会返回base原数组
func PadBlock(base []byte, pad byte, blockSize int, separator ...byte) []byte {
	sep := First(separator)
	paddingLength := blockSize - len(base)%blockSize
	if paddingLength == 0 && sep == 0 {
		return base
	}
	padding := bytes.Repeat([]byte{pad}, paddingLength)
	if sep == 0 {
		return append(base, padding...)
	}

	// 每个block尾部插入分隔符
	bn := len(base) / blockSize
	blockNum := bn + 1
	var result bytes.Buffer
	result.Grow(len(base) + blockNum + len(padding))
	for i := 0; i < bn; i++ {
		result.Write(base[i*blockSize : (i+1)*blockSize])
		result.WriteByte(sep)
	}
	result.Write(base[bn*blockSize:])
	result.Write(padding)
	result.WriteByte(sep)
	return result.Bytes()
}
