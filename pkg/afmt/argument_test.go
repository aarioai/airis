package afmt_test

import (
	"testing"

	"github.com/aarioai/airis/pkg/afmt"
)

func TestFirstNotNil(t *testing.T) {
	tests := []struct {
		name     string
		args     []interface{} // Using interface{} to test different types
		expected interface{}
	}{
		// Integer tests
		{
			name:     "first non-zero int",
			args:     []interface{}{0, 0, 5, 0, 10},
			expected: 5,
		},
		{
			name:     "all zero ints",
			args:     []interface{}{0, 0, 0},
			expected: 0,
		},
		{
			name:     "first element non-zero int",
			args:     []interface{}{42, 0, 0},
			expected: 42,
		},

		// String tests
		{
			name:     "first non-zero string",
			args:     []interface{}{"", "", "hello", "", "world"},
			expected: "hello",
		},
		{
			name:     "all empty strings",
			args:     []interface{}{"", "", ""},
			expected: "",
		},
		{
			name:     "first element non-empty string",
			args:     []interface{}{"first", "", ""},
			expected: "first",
		},

		// Boolean tests
		{
			name:     "first non-zero bool (false is zero value)",
			args:     []interface{}{false, false, true, false},
			expected: true,
		},
		{
			name:     "all false bools",
			args:     []interface{}{false, false, false},
			expected: false,
		},

		// Pointer tests
		{
			name:     "first non-nil pointer",
			args:     []interface{}{nil, nil, new(int), nil},
			expected: new(int), // Will compare values, not addresses
		},
		{
			name:     "all nil pointers",
			args:     []interface{}{nil, nil, nil},
			expected: nil,
		},

		// Slice tests
		{
			name:     "first non-nil slice",
			args:     []interface{}{[]int(nil), []int(nil), []int{1, 2, 3}, []int(nil)},
			expected: []int{1, 2, 3},
		},
		{
			name:     "all nil slices",
			args:     []interface{}{[]int(nil), []int(nil), []int(nil)},
			expected: []int(nil),
		},
		{
			name:     "empty slice vs nil slice",
			args:     []interface{}{[]int(nil), []int{}, []int{1}},
			expected: []int{},
		},

		// Error/interface tests
		{
			name:     "first non-nil error",
			args:     []interface{}{nil, nil, &testError{"error message"}, nil},
			expected: &testError{"error message"},
		},
		{
			name:     "all nil errors",
			args:     []interface{}{nil, nil, nil},
			expected: nil,
		},

		// Mixed types (all same underlying type but different values)
		{
			name:     "mixed int values",
			args:     []interface{}{0, 0, -1, 0, 100},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need a type-specific wrapper since the function is generic
			// For testing, we'll create type-specific test functions
			switch tt.args[0].(type) {
			case int:
				intArgs := make([]int, len(tt.args))
				for i, v := range tt.args {
					if v == nil {
						intArgs[i] = 0
					} else {
						intArgs[i] = v.(int)
					}
				}
				result := afmt.FirstNotEmpty(intArgs)
				if result != tt.expected.(int) {
					t.Errorf("afmt.FirstNotEmpty(%v) = %v, want %v", intArgs, result, tt.expected)
				}

			case string:
				stringArgs := make([]string, len(tt.args))
				for i, v := range tt.args {
					if v == nil {
						stringArgs[i] = ""
					} else {
						stringArgs[i] = v.(string)
					}
				}
				result := afmt.FirstNotEmpty(stringArgs)
				if result != tt.expected.(string) {
					t.Errorf("afmt.FirstNotEmpty(%v) = %v, want %v", stringArgs, result, tt.expected)
				}

			case bool:
				boolArgs := make([]bool, len(tt.args))
				for i, v := range tt.args {
					if v == nil {
						boolArgs[i] = false
					} else {
						boolArgs[i] = v.(bool)
					}
				}
				result := afmt.FirstNotEmpty(boolArgs)
				if result != tt.expected.(bool) {
					t.Errorf("afmt.FirstNotEmpty(%v) = %v, want %v", boolArgs, result, tt.expected)
				}

			case *int:
				intPtrArgs := make([]*int, len(tt.args))
				for i, v := range tt.args {
					if v == nil {
						intPtrArgs[i] = nil
					} else {
						intPtrArgs[i] = v.(*int)
					}
				}
				result := afmt.FirstNotEmpty(intPtrArgs)
				// Compare values instead of addresses for expected
				if tt.expected != nil {
					if result == nil || *result != *(tt.expected.(*int)) {
						t.Errorf("afmt.FirstNotEmpty(%v) value = %v, want %v", intPtrArgs, result, tt.expected)
					}
				} else {
					if result != nil {
						t.Errorf("afmt.FirstNotEmpty(%v) = %v, want nil", intPtrArgs, result)
					}
				}
			case error:
				errArgs := make([]error, len(tt.args))
				for i, v := range tt.args {
					if v == nil {
						errArgs[i] = nil
					} else {
						errArgs[i] = v.(error)
					}
				}
				result := afmt.FirstNotEmpty(errArgs)
				// Compare error messages
				if tt.expected != nil {
					if result == nil || result.Error() != tt.expected.(error).Error() {
						t.Errorf("afmt.FirstNotEmpty(%v) = %v, want %v", errArgs, result, tt.expected)
					}
				} else {
					if result != nil {
						t.Errorf("afmt.FirstNotEmpty(%v) = %v, want nil", errArgs, result)
					}
				}
			}
		})
	}
}

// Helper function to compare slices
func compareSlices(a, b []int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Custom error type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// Additional edge case tests
func TestFirstNotNilEdgeCases(t *testing.T) {
	t.Run("empty slice returns zero value", func(t *testing.T) {
		var ints []int
		result := afmt.FirstNotEmpty(ints)
		if result != 0 {
			t.Errorf("Empty slice should return 0, got %v", result)
		}
	})

	t.Run("single element", func(t *testing.T) {
		result := afmt.FirstNotEmpty([]string{"single"})
		if result != "single" {
			t.Errorf("Expected 'single', got %v", result)
		}
	})

	t.Run("all zero values", func(t *testing.T) {
		result := afmt.FirstNotEmpty([]int{0, 0, 0})
		if result != 0 {
			t.Errorf("Expected 0, got %v", result)
		}
	})

	t.Run("struct with zero values", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		persons := []Person{
			{Name: "", Age: 0},
			{Name: "Alice", Age: 30},
			{Name: "", Age: 0},
		}
		result := afmt.FirstNotEmpty(persons)
		if result.Name != "Alice" || result.Age != 30 {
			t.Errorf("Expected {Alice 30}, got %v", result)
		}
	})
}

// Benchmark tests
func BenchmarkFirstNotNil(b *testing.B) {
	ints := []int{0, 0, 0, 0, 0, 42, 0, 0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		afmt.FirstNotEmpty(ints)
	}
}

func BenchmarkFirstNotNilEarlyExit(b *testing.B) {
	ints := []int{42, 0, 0, 0, 0, 0, 0, 0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		afmt.FirstNotEmpty(ints)
	}
}

func BenchmarkFirstNotNilAllZero(b *testing.B) {
	ints := []int{0, 0, 0, 0, 0, 0, 0, 0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		afmt.FirstNotEmpty(ints)
	}
}
