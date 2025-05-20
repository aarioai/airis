package basic

// Ter a ternary
func Ter[T any](cond bool, a T, b T) T {
	if cond {
		return a
	}
	return b
}
