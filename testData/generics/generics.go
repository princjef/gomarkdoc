package generics

// Generic is a generic struct.
type Generic[T any] struct {
	Field T
}

// NewGeneric produces a new [Generic] struct.
func NewGeneric[T any](param T) Generic[T] {
	return Generic[T]{Field: param}
}

// Method is a method of a generic type.
func (g Generic[T]) Method() {}

// Func is a generic function.
func Func[S int | float64](s S) S {
	return s
}
