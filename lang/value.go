package lang

import (
	"go/doc"
	"go/token"
)

// Value holds documentation for a var or const declaration within a package.
type Value struct {
	level int
	doc   *doc.Value
	fs    *token.FileSet
}

// NewValue creates a new Value from the raw const or var documentation and the
// token.FileSet of files for the containing package.
func NewValue(doc *doc.Value, fs *token.FileSet, level int) *Value {
	return &Value{level, doc, fs}
}

// Level provides the default level that headers for the value should be
// rendered.
func (v *Value) Level() int {
	return v.level
}

// Summary provides the one-sentence summary of the value's documentation
// comment.
func (v *Value) Summary() string {
	return extractSummary(v.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// example.
func (v *Value) Doc() *Doc {
	return NewDoc(v.doc.Doc, v.level+1)
}

// Decl provides the raw text representation of the code for declaring the const
// or var.
func (v *Value) Decl() (string, error) {
	return printNode(v.doc.Decl, v.fs)
}
