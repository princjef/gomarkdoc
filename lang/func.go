package lang

import (
	"fmt"
	"go/doc"
	"strings"
)

// Func holds documentation information for a single func declaration within a
// package or type.
type Func struct {
	cfg      *Config
	doc      *doc.Func
	examples []*doc.Example
}

// NewFunc creates a new Func from the corresponding documentation construct
// from the standard library, the related token.FileSet for the package and
// the list of examples for the package.
func NewFunc(cfg *Config, doc *doc.Func, examples []*doc.Example) *Func {
	return &Func{cfg, doc, examples}
}

// Level provides the default level at which headers for the func should be
// rendered in the final documentation.
func (fn *Func) Level() int {
	return fn.cfg.Level
}

// Name provides the name of the function.
func (fn *Func) Name() string {
	return fn.doc.Name
}

// Title provides the formatted name of the func. It is primarily designed for
// generating headers.
func (fn *Func) Title() string {
	if fn.doc.Recv != "" {
		return fmt.Sprintf("func (%s) %s", fn.doc.Recv, fn.doc.Name)
	}

	return fmt.Sprintf("func %s", fn.doc.Name)
}

// Receiver provides the type of the receiver for the function, or empty string
// if there is no receiver type.
func (fn *Func) Receiver() string {
	return fn.doc.Recv
}

// Location returns a representation of the node's location in a file within a
// repository.
func (fn *Func) Location() Location {
	return NewLocation(fn.cfg, fn.doc.Decl)
}

// Summary provides the one-sentence summary of the function's documentation
// comment
func (fn *Func) Summary() string {
	return extractSummary(fn.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// function.
func (fn *Func) Doc() *Doc {
	return NewDoc(fn.cfg.Inc(1), fn.doc.Doc)
}

// Signature provides the raw text representation of the code for the
// function's signature.
func (fn *Func) Signature() (string, error) {
	return printNode(fn.doc.Decl, fn.cfg.FileSet)
}

// Examples provides the list of examples from the list given on initialization
// that pertain to the function.
func (fn *Func) Examples() (examples []*Example) {
	var fullName string
	if fn.doc.Recv != "" {
		fullName = fmt.Sprintf("%s_%s", fn.rawRecv(), fn.doc.Name)
	} else {
		fullName = fn.doc.Name
	}
	underscorePrefix := fmt.Sprintf("%s_", fullName)

	for _, example := range fn.examples {
		var name string
		switch {
		case example.Name == fullName:
			name = ""
		case strings.HasPrefix(example.Name, underscorePrefix):
			name = underscorePrefix[len(underscorePrefix):]
		default:
			// TODO: better filtering
			continue
		}

		examples = append(examples, NewExample(fn.cfg.Inc(1), name, example))
	}

	return
}

func (fn *Func) rawRecv() string {
	if strings.HasPrefix(fn.doc.Recv, "*") {
		return fn.doc.Recv[1:]
	}

	return fn.doc.Recv
}
