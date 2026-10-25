package atype

import "regexp"

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isHexByte(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func IsUUID(s string) bool {
	// 32: hex without dashes. 36: canonical 8-4-4-4-12 with version/variant nibbles.
	switch len(s) {
	case 32:
		for i := 0; i < 32; i++ {
			if !isHexByte(s[i]) {
				return false
			}
		}
		v := s[12]
		if v < '1' || v > '5' {
			return false
		}
		switch s[16] {
		case '8', '9', 'a', 'b', 'A', 'B':
			return true
		default:
			return false
		}
	case 36:
		if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			return false
		}
		for _, i := range [...]int{0, 1, 2, 3, 4, 5, 6, 7, 9, 10, 11, 12, 14, 15, 16, 17, 19, 20, 21, 22, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35} {
			if !isHexByte(s[i]) {
				return false
			}
		}
		v := s[14]
		if v < '1' || v > '5' {
			return false
		}
		switch s[19] {
		case '8', '9', 'a', 'b', 'A', 'B':
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func IsNumberString(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsNumber(char) {
			return false
		}
	}
	return true
}

func IsLowers(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsLower(char) {
			return false
		}
	}
	return true
}

func IsUppers(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsUpper(char) {
			return false
		}
	}
	return true
}

func IsAlphabetical(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsAlphabeticalChar(char) {
			return false
		}
	}
	return true
}

func IsAlphaDigit(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsAlphaDigitChar(char) {
			return false
		}
	}
	return true
}

func IsWord(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsWordChar(char) {
			return false
		}
	}
	return true
}

func IsFilename(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsFilenameChar(char) {
			return false
		}
	}
	return true
}
func IsUnicodeFilename(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsUnicodeFilenameChar(char) {
			return false
		}
	}
	return true
}

func IsPath(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsPathChar(char) {
			return false
		}
	}
	return true
}
func IsUnicodePath(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsUnicodePathChar(char) {
			return false
		}
	}
	return true
}

// IsEmail checks if a string is a valid email address.
// It uses a regular expression that matches most common email formats,
// but does not cover 100% of all possible RFC-compliant email addresses.
// For stricter validation, consider using a dedicated validation library.
func IsEmail(s string) bool {
	if s == "" {
		return false
	}
	return emailPattern.MatchString(s)
}

func IsBin(s string) bool {
	if s == "" {
		return false
	}
	for _, char := range s {
		if !IsBinChar(char) {
			return false
		}
	}
	return true
}
