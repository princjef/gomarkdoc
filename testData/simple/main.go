// Package simple contains, some simple code to exercise basic scenarios
// for documentation purposes.
package simple

// Num is a number.
//
// It is just a test type so that we can make sure this works.
type Num int

// Add adds the other num to this one.
func (n Num) Add(num Num) Num {
	return addInternal(n, num)
}

// AddNums adds two Nums together.
func AddNums(num1, num2 Num) Num {
	return addInternal(num1, num2)
}

// addInternal is a private version of AddNums.
func addInternal(num1, num2 Num) Num {
	return num1 + num2
}
