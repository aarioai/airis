package afmt

import (
	"strings"

	"github.com/aarioai/airis/pkg/types"
)

// Resize pads a []byte or []rune slice to the given length, and truncates it when it is longer.
func Resize[T ~byte | ~rune](s []T, pad T, length int, padHead bool) []T {
	paddingLen := length - len(s)
	if length < 0 || paddingLen == 0 {
		return s
	}
	if paddingLen > 0 {
		padding := make([]T, paddingLen)
		for i := 0; i < paddingLen; i++ {
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

// PadBoth pads both ends of a string so it reaches the given length. The result is at least length, unless padStr is empty.
// If padStr is too long to fit in length, it is truncated from the end. This matches JS padStart/padEnd.
// Usually used to pad short strings, so it returns a string directly.
//
// @param trimEdge trims extra padding from both outer edges, e.g. ~_~||~_~|AARIO|~_~||~_~; otherwise trims from the inner side, e.g. |~_~||~_~AARIO~_~||~_~|
//
// @note slice arguments are safe
func PadBoth[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padStr := string(pad)
	paddingLen := length - len(s)

	// If no padding is needed, return the original string.
	if paddingLen <= 0 || padStr == "" {
		return s
	}

	// Compute the padding length for each side.
	leftPadLen := paddingLen / 2
	rightPadLen := paddingLen - leftPadLen

	// Build the padding strings.
	leftRepeatCount := leftPadLen
	rightRepeatCount := rightPadLen
	psLen := len(padStr)
	if psLen > 1 {
		leftRepeatCount = (leftPadLen / psLen) + 1
		rightRepeatCount = (rightPadLen / psLen) + 1
	}
	leftPadding := strings.Repeat(padStr, leftRepeatCount)
	rightPadding := strings.Repeat(padStr, rightRepeatCount)

	// Trim extra padding from both outer edges, e.g. ~_~||~_~|AARIO|~_~||~_~.
	if First(trimEdge) {
		leftPadding = leftPadding[len(leftPadding)-leftPadLen:]
		rightPadding = rightPadding[:rightPadLen]
	} else {
		// Trim extra padding from the inner side, e.g. |~_~||~_~AARIO~_~||~_~|.
		leftPadding = leftPadding[:leftPadLen]
		rightPadding = rightPadding[len(rightPadding)-rightPadLen:]
	}

	return leftPadding + s + rightPadding
}

// PadLeft
// fmt.Sprintf("%04d", 10) and fmt.Sprintf("%-04d", 3) can only pad with spaces on either side, or with 0 on the left.
// Usually used to pad short strings, so it returns a string directly.
//
// @param trimEdge trims extra padding from the outer edge, e.g. ~_~||~_~|AARIO; otherwise trims from the inner side, e.g. |~_~||~_~AARIO
//
// @note slice arguments are safe
func PadLeft[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padStr := string(pad)
	paddingLen := length - len(s)
	if paddingLen <= 0 || padStr == "" {
		return s
	}
	repeatCount := paddingLen
	if len(padStr) > 1 {
		repeatCount = (paddingLen / len(padStr)) + 1
	}
	padding := strings.Repeat(padStr, repeatCount)
	if First(trimEdge) {
		padding = padding[len(padding)-paddingLen:]
	} else {
		padding = padding[:paddingLen]
	}
	return padding + s
}

// PadRight
// Usually used to pad short strings, so it returns a string directly.
//
// @param trimEdge trims extra padding from the outer edge, e.g. AARIO|~_~||~_~; otherwise trims from the inner side, e.g. AARIO~_~||~_~|
//
// @note slice arguments are safe
func PadRight[T1, T2 types.Stringable](base T1, pad T2, length int, trimEdge ...bool) string {
	s := string(base)
	padStr := string(pad)
	paddingLen := length - len(s)
	if paddingLen <= 0 || padStr == "" {
		return s
	}
	repeatCount := paddingLen
	if len(padStr) > 1 {
		repeatCount = (paddingLen / len(padStr)) + 1
	}
	padding := strings.Repeat(padStr, repeatCount)
	if First(trimEdge) {
		padding = padding[:paddingLen]
	} else {
		padding = padding[len(padding)-paddingLen:]
	}

	return s + padding
}

// TrimRight cuts every trailing cut element. To trim a string, use strings.Trim, bytes.Trim, or TrimFunc.
// This extends the byte == 0 case (DES padding alignment uses byte(0)).
//
// @note slice arguments are safe
// @warn the result may alias part of the input slice
func TrimRight[T ~byte | ~rune](s []T, cut T) []T {
	// var zero byte cannot use len(zero); that does not compile. Use len(string(zero)) == 1, or test whether the value is 0.
	var length int
	if cut == 0 {
		length = 1
	} else {
		length = len(s) // For []rune, length is the number of runes, not the number of bytes.
	}
	i := length - 1
	for ; i > -1; i-- {
		if s[i] != cut {
			break
		}
	}
	return s[:i+1]
}

func TrimLeft[T ~byte | ~rune](s []T, cut T) []T {
	// var zero byte cannot use len(zero).
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
func Trim[T ~byte | ~rune](s []T, cut T) []T {
	return TrimLeft(TrimRight(s, cut), cut)
}
