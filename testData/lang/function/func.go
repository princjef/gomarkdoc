package function

// Standalone provides a function that is not part of a type.
//
// Additional description can be provided in subsequent paragraphs, including
// code blocks and headers
//
// Header A
//
// This section contains a code block.
//
// 	Code Block
// 	More of Code Block
func Standalone(p1 int, p2 string) (int, error) {
	return p1, nil
}

// Receiver is a type used to demonstrate functions with receivers.
type Receiver struct{}

// New is an initializer for Receiver.
func New() Receiver {
	return Receiver{}
}

// WithReceiver has a receiver.
func (r Receiver) WithReceiver() {}

// WithPtrReceiver has a pointer receiver.
func (r *Receiver) WithPtrReceiver() {}

// Generic is a struct with a generic type.
type Generic[T any] struct{}

// WithGenericReceiver has a receiver with a generic type.
func (r Generic[T]) WithGenericReceiver() {}
