package lang

import (
	"fmt"
	"go/doc"
	"go/printer"
	"strings"
)

// Example holds a single documentation example for a package or symbol.
type Example struct {
	cfg  *Config
	name string
	doc  *doc.Example
}

// NewExample creates a new example from the example function's name, its
// documentation example and the files holding code related to the example.
func NewExample(cfg *Config, name string, doc *doc.Example) *Example {
	return &Example{cfg, name, doc}
}

// Level provides the default level that headers for the example should be
// rendered.
func (ex *Example) Level() int {
	return ex.cfg.Level
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

// Location returns a representation of the node's location in a file within a
// repository.
func (ex *Example) Location() Location {
	return NewLocation(ex.cfg, ex.doc.Code)
}

// Summary provides the one-sentence summary of the example's documentation
// comment.
func (ex *Example) Summary() string {
	return extractSummary(ex.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// example.
func (ex *Example) Doc() *Doc {
	return NewDoc(ex.cfg.Inc(1), ex.doc.Doc)
}

// Code provides the raw text code representation of the example's contents.
func (ex *Example) Code() (string, error) {
	var codeNode interface{}
	if ex.doc.Play != nil {
		codeNode = ex.doc.Play
	} else {
		codeNode = &printer.CommentedNode{Node: ex.doc.Code, Comments: ex.doc.Comments}
	}

	var code strings.Builder
	p := &printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}
	err := p.Fprint(&code, ex.cfg.FileSet, codeNode)
	if err != nil {
		return "", err
	}

	str := code.String()

	// additional formatting if this is a function body
	if i := len(str); i >= 2 && str[0] == '{' && str[i-1] == '}' {
		// remove surrounding braces
		str = str[1 : i-1]
		// unindent
		str = strings.ReplaceAll(str, "\n\t", "\n")
	}

	return str, nil
}

// Output provides the code's example output.
func (ex *Example) Output() string {
	return ex.doc.Output
}

// HasOutput indicates whether the example contains any example output.
func (ex *Example) HasOutput() bool {
	return ex.doc.Output != "" || ex.doc.EmptyOutput
}
