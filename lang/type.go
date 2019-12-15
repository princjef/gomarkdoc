package lang

import (
	"fmt"
	"go/doc"
	"go/token"
	"strings"
)

// Type holds documentation information for a type declaration.
type Type struct {
	level    int
	doc      *doc.Type
	fs       *token.FileSet
	examples []*doc.Example
}

// NewType creates a Type from the raw documentation representation of the type,
// the token.FileSet for the package's files and the full list of examples from
// the containing package.
func NewType(doc *doc.Type, fs *token.FileSet, examples []*doc.Example, level int) *Type {
	return &Type{level, doc, fs, examples}
}

// Level provides the default level that headers for the type should be
// rendered.
func (typ *Type) Level() int {
	return typ.level
}

// Name provides the name of the type
func (typ *Type) Name() string {
	return typ.doc.Name
}

// Title provides a formatted name suitable for use in a header identifying the
// type.
func (typ *Type) Title() string {
	return fmt.Sprintf("type %s", typ.doc.Name)
}

// Summary provides the one-sentence summary of the type's documentation
// comment.
func (typ *Type) Summary() string {
	return extractSummary(typ.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// type.
func (typ *Type) Doc() *Doc {
	return NewDoc(typ.doc.Doc, typ.level+1)
}

// Decl provides the raw text representation of the code for the type's
// declaration.
func (typ *Type) Decl() (string, error) {
	return printNode(typ.doc.Decl, typ.fs)
}

// Examples lists the examples pertaining to the type from the set provided on
// initialization.
func (typ *Type) Examples() (examples []*Example) {
	underscorePrefix := fmt.Sprintf("%s_", typ.doc.Name)
	for _, example := range typ.examples {
		var name string
		switch {
		case example.Name == typ.doc.Name:
			name = ""
		case strings.HasPrefix(example.Name, underscorePrefix) && !typ.isSubexample(example.Name):
			name = example.Name[len(underscorePrefix):]
		default:
			// TODO: better filtering
			continue
		}

		examples = append(examples, NewExample(name, example, typ.fs, typ.level+1))
	}

	return
}

func (typ *Type) isSubexample(exampleName string) bool {
	for _, m := range typ.doc.Methods {
		fullName := fmt.Sprintf("%s_%s", typ.doc.Name, m.Name)
		underscorePrefix := fmt.Sprintf("%s_", fullName)
		if exampleName == fullName || strings.HasPrefix(exampleName, underscorePrefix) {
			return true
		}
	}

	return false
}

// Funcs lists the funcs related to the type. This only includes functions which
// return an instance of the type or its pointer.
func (typ *Type) Funcs() (funcs []*Func) {
	for _, fn := range typ.doc.Funcs {
		funcs = append(funcs, NewFunc(fn, typ.fs, typ.examples, typ.level+1))
	}

	return
}

// Methods lists the funcs that use the type as a value or pointer receiver.
func (typ *Type) Methods() (methods []*Func) {
	for _, fn := range typ.doc.Methods {
		methods = append(methods, NewFunc(fn, typ.fs, typ.examples, typ.level+1))
	}

	return
}
