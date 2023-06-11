package docs

// AnotherStruct has methods like [*AnotherStruct.GetField] and also has an
// initializer called [NewAnotherStruct].
type AnotherStruct struct {
	Field string
}

// NewAnotherStruct() makes [*AnotherStruct].
func NewAnotherStruct() *AnotherStruct {
	return &AnotherStruct{
		Field: "test",
	}
}

// GetField gets [*AnotherStruct.Field].
func (s *AnotherStruct) GetField() string {
	return s.Field
}
