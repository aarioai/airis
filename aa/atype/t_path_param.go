package atype

func NewUUID(s string) (UUID, bool) {
	return UUID(s), IsUUID(s)
}

func (s UUID) Valid() bool {
	return IsUUID(string(s))
}

func NewNumberString(s string) (NumberString, bool) {
	return NumberString(s), IsNumberString(s)
}

func (s NumberString) Valid() bool {
	return IsNumberString(string(s))
}

func NewLowers(s string) (Lowers, bool) {
	return Lowers(s), IsLowers(s)
}
func (s Lowers) Valid() bool {
	return IsLowers(string(s))
}

func NewUppers(s string) (Uppers, bool) {
	return Uppers(s), IsUppers(s)
}

func (s Uppers) Valid() bool {
	return IsUppers(string(s))
}

func NewAlphabetical(s string) (Alphabetical, bool) {
	return Alphabetical(s), IsAlphabetical(s)
}

func (s Alphabetical) Valid() bool {
	return IsAlphabetical(string(s))
}

func NewAlphaDigit(s string) (AlphaDigit, bool) {
	return AlphaDigit(s), IsAlphaDigit(s)
}

func (s AlphaDigit) Valid() bool {
	return IsAlphaDigit(string(s))
}

func NewWord(s string) (Word, bool) {
	return Word(s), IsWord(s)
}

func (s Word) Valid() bool {
	return IsWord(string(s))
}

func NewFilename(s string) (Filename, bool) {
	return Filename(s), IsFilename(s)
}

func (s Filename) Valid() bool {
	return IsFilename(string(s))
}

func NewUnicodeFilename(s string) (UnicodeFilename, bool) {
	return UnicodeFilename(s), IsUnicodeFilename(s)
}

func (s UnicodeFilename) Valid() bool {
	return IsUnicodeFilename(string(s))
}

func NewPath(s string) (Path, bool) {
	return Path(s), IsPath(s)
}

func (s Path) Valid() bool {
	return IsPath(string(s))
}

func NewUnicodePath(s string) (UnicodePath, bool) {
	return UnicodePath(s), IsUnicodePath(s)
}

func (s UnicodePath) Valid() bool {
	return IsUnicodePath(string(s))
}

func NewEmail(s string) (Email, bool) {
	return Email(s), IsEmail(s)
}

func (s Email) Valid() bool {
	return IsEmail(string(s))
}
