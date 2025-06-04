package arrmap

import "fmt"

// IsBase checks if the digits slice has the specified length and contains no duplicate values
// Example:
// IsBase(36, []byte("0123456789abcdefghijklmnopqrstuvwxyz"))   // check is base36 base bytes
func IsBase(base int, digits []byte) bool {
	if len(digits) != base {
		return false
	}
	return !HasDuplicates(digits)
}

func PanicIfNotBase(name string, base int, digits []byte) {
	if !IsBase(base, digits) {
		panic(fmt.Sprintf("%s: %s is not the base bytes of base%d", name, string(digits), base))
	}
}
