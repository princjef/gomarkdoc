package lang

import (
	"fmt"
	"go/doc"
	"go/format"
	"go/token"
	"strings"
)

// Example holds a single documentation example for a package or symbol.
type Example struct {
	name  string
	doc   *doc.Example
	fs    *token.FileSet
	level int
}

// NewExample creates a new example from the example function's name, its
// documentation example and the files holding code related to the example.
func NewExample(name string, doc *doc.Example, fs *token.FileSet, level int) *Example {
	return &Example{name, doc, fs, level}
}

// Level provides the default level that headers for the example should be
// rendered.
func (ex *Example) Level() int {
	return ex.level
}

// Name provides a pretty-printed name for the specific example, if one was
// provided.
func (ex *Example) Name() string {
	return splitCamel(ex.name)
}

// Title provides a formatted string to print as the title of the example. It
// incorporates the example's name, if present.
func (ex *Example) Title() string {
	name := ex.Name()
	if name == "" {
		return "Example"
	}

	return fmt.Sprintf("Example (%s)", name)
}

// Summary provides the one-sentence summary of the example's documentation
// comment.
func (ex *Example) Summary() string {
	return extractSummary(ex.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// example.
func (ex *Example) Doc() *Doc {
	return NewDoc(ex.doc.Doc, ex.level+1)
}

// Code provides the raw text code representation of the example's contents.
func (ex *Example) Code() (string, error) {
	var codeNode interface{}
	if ex.doc.Play != nil {
		codeNode = ex.doc.Play
	} else {
		codeNode = ex.doc.Code
	}

	var code strings.Builder
	if err := format.Node(&code, ex.fs, codeNode); err != nil {
		return "", err
	}

	return code.String(), nil
}
