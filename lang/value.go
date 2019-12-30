package lang

import (
	"go/doc"
)

// Value holds documentation for a var or const declaration within a package.
type Value struct {
	cfg *Config
	doc *doc.Value
}

// NewValue creates a new Value from the raw const or var documentation and the
// token.FileSet of files for the containing package.
func NewValue(cfg *Config, doc *doc.Value) *Value {
	return &Value{cfg, doc}
}

// Level provides the default level that headers for the value should be
// rendered.
func (v *Value) Level() int {
	return v.cfg.Level
}

// Location returns a representation of the node's location in a file within a
// repository.
func (v *Value) Location() Location {
	return NewLocation(v.cfg, v.doc.Decl)
}

// Summary provides the one-sentence summary of the value's documentation
// comment.
func (v *Value) Summary() string {
	return extractSummary(v.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// example.
func (v *Value) Doc() *Doc {
	return NewDoc(v.cfg.Inc(1), v.doc.Doc)
}

// Decl provides the raw text representation of the code for declaring the const
// or var.
func (v *Value) Decl() (string, error) {
	return printNode(v.doc.Decl, v.cfg.FileSet)
}
