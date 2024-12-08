package afmt

import "strings"

// PadBoth 在字符串两端添加填充字符，使其达到指定长度
// 如果 padString 太长，无法适应 length，则会从末尾被截断。这个跟 JS padStart/padEnd 一致
// 返回的长度大于等于 length，除非 padString 为空
func PadBoth(msg string, length int, padString string, startFromEnd bool) string {
	msgLen := len(msg)
	padLen := length - msgLen
	
	// 如果不需要填充，直接返回原字符串
	if padLen <= 0 || padString == "" {
		return msg
	}
	
	// 计算两侧需要的填充长度
	leftPadLen := padLen / 2
	rightPadLen := padLen - leftPadLen
	
	// 如果长度不均匀且需要从开头多填充一个字符
	if padLen%2 != 0 && !startFromEnd {
		leftPadLen, rightPadLen = rightPadLen, leftPadLen
	}
	
	// 生成填充字符串
	leftPad := strings.Repeat(padString, (leftPadLen+len(padString)-1)/len(padString))
	rightPad := strings.Repeat(padString, (rightPadLen+len(padString)-1)/len(padString))
	
	// 截取需要的长度
	leftPad = leftPad[:leftPadLen]
	rightPad = rightPad[:rightPadLen]
	
	return leftPad + msg + rightPad
}
