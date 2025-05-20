package atype

func IsNumber[T byte | rune](char T) bool {
	return char >= '0' && char <= '9'
}

func IsLower[T byte | rune](char T) bool {
	return char >= 'a' && char <= 'z'
}

func IsUpper[T byte | rune](char T) bool {
	return char >= 'A' && char <= 'Z'
}

func IsAlphabeticalChar[T byte | rune](char T) bool {
	return IsLower[T](char) || IsUpper[T](char)
}

func IsAlphaDigitChar[T byte | rune](char T) bool {
	return IsNumber[T](char) || IsAlphabeticalChar[T](char)
}

func IsWordChar[T byte | rune](char T) bool {
	return IsAlphabeticalChar[T](char) || char == '_'
}
func IsFilenameChar[T byte | rune](char T) bool {
	return IsWordChar(char) || char == '-'
}

func IsPathChar[T byte | rune](char T) bool {
	return IsFilenameChar(char) || char == '/'
}

func IsUnicodeFilenameChar[T byte | rune](char T) bool {
	return IsFilenameChar(char) ||
		char == '!' ||
		char == '@' ||
		char == '#' ||
		char == '$' ||
		char == '%' ||
		char == '^' ||
		char == '&' ||
		char == '(' ||
		char == ')' ||
		char == '{' ||
		char == '}' ||
		char == '~'
}

func IsUnicodePathChar[T byte | rune](char T) bool {
	return IsUnicodeFilenameChar(char) || char == '/'
}
func IsBinChar[T byte | rune](char T) bool {
	return char == '0' || char == '1'
}
