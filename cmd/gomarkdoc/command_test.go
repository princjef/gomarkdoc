package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		b1, b2 []byte
		equal  bool
	}{
		{[]byte("abc"), []byte("abc"), true},
		{[]byte("abc"), []byte("def"), false},
		{[]byte{}, []byte{}, true},
		{[]byte("abc"), []byte{}, false},
		{[]byte{}, []byte("abc"), false},
	}

	for _, test := range tests {
		name := fmt.Sprintf(`"%s" == "%s"`, string(test.b1), string(test.b2))
		if !test.equal {
			name = fmt.Sprintf(`"%s" != "%s"`, string(test.b1), string(test.b2))
		}

		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			eq, err := compare(bytes.NewBuffer(test.b1), bytes.NewBuffer(test.b2))
			is.NoErr(err)

			is.Equal(eq, test.equal)
		})
	}
}
