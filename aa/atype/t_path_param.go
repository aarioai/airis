package atype

import "slices"

// Iris path parameter types
// See https://iris-go.gitbook.io/iris/contents/routing/routing-path-parameter-types

var ParamTypes = []string{
	":string", ":uuid",
	":int8", ":int16", ":int32", ":int64", // no :int24 and :uint24
	":uint8", ":uint16", ":uint32", ":uint64",
	":bool", ":alphabetical",
	":email", ":mail", // mail is same to email, but mail without server domain validation
	":weekday",
}

func NewPathParamType(t string) (PathParamType, bool) {
	return PathParamType(t), slices.Contains(ParamTypes, t)
}

func (t PathParamType) Valid() bool {
	return slices.Contains(ParamTypes, string(t))
}

func NewUUID(s string) (UUID, bool) {
	return UUID(s), IsUUID(s)
}

func (s UUID) Valid() bool {
	return IsUUID(string(s))
}

func NewNumberString(s string) (Numerical, bool) {
	return Numerical(s), IsNumberString(s)
}

func (s Numerical) Valid() bool {
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

func (s Email) Valid() bool {
	return IsEmail(string(s))
}
