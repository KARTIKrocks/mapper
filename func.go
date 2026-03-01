package mapper

// Func is a type-safe mapping function with zero reflection overhead.
// It wraps a user-provided conversion function and provides convenience
// methods for single and slice mapping.
//
// Use Func when you want maximum performance and full compile-time safety.
type Func[S, D any] struct {
	fn func(S) D
}

// NewFunc creates a typed mapping function.
// It panics if fn is nil.
func NewFunc[S, D any](fn func(S) D) Func[S, D] {
	if fn == nil {
		panic("mapper: mapping function must not be nil")
	}
	return Func[S, D]{fn: fn}
}

// Map applies the mapping function to src and returns the result.
func (f Func[S, D]) Map(src S) D {
	return f.fn(src)
}

// MapSlice applies the mapping function to each element of src.
// Returns nil if src is nil.
func (f Func[S, D]) MapSlice(src []S) []D {
	if src == nil {
		return nil
	}
	result := make([]D, len(src))
	for i, s := range src {
		result[i] = f.fn(s)
	}
	return result
}
