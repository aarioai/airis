package atype

import "regexp"

func IsUUID(s string) bool {
	if len(s) != 32 {
		return false
	}
	pattern := `^[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?[1-5][0-9a-fA-F]{3}-?[89abAB][0-9a-fA-F]{3}-?[0-9a-fA-F]{12}$`
	return regexp.MustCompile(pattern).MatchString(s)
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
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(s)
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
