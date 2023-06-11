package lang

import (
	"fmt"
	"go/ast"
	"go/doc"
	"strings"
)

type (
	// Symbol provides identity information for a symbol in a package.
	Symbol struct {
		// Receiver holds the receiver for a method or field.
		Receiver string
		// Name holds the name of the symbol itself.
		Name string
		// Kind identifies the category of the symbol.
		Kind SymbolKind
		// Parent holds the linkable parent symbol which contains this one.
		Parent *Symbol
	}

	// SymbolKind identifies the type of symbol.
	SymbolKind int
)

// The list of valid symbol kinds.
const (
	TypeSymbolKind SymbolKind = iota + 1
	FuncSymbolKind
	ConstSymbolKind
	VarSymbolKind
	MethodSymbolKind
	FieldSymbolKind
)

// PackageSymbols gets the list of symbols for a doc package.
func PackageSymbols(pkg *doc.Package) map[string]Symbol {
	sym := make(map[string]Symbol)
	for _, c := range pkg.Consts {
		parent := Symbol{
			Name: c.Names[0],
			Kind: ConstSymbolKind,
		}

		for _, n := range c.Names {
			sym[symbolName("", n)] = Symbol{
				Name:   n,
				Kind:   ConstSymbolKind,
				Parent: &parent,
			}
		}
	}

	for _, v := range pkg.Vars {
		parent := Symbol{
			Name: v.Names[0],
			Kind: VarSymbolKind,
		}

		for _, n := range v.Names {
			sym[symbolName("", n)] = Symbol{
				Name:   n,
				Kind:   VarSymbolKind,
				Parent: &parent,
			}
		}
	}

	for _, v := range pkg.Funcs {
		sym[symbolName("", v.Name)] = Symbol{
			Name: v.Name,
			Kind: FuncSymbolKind,
		}
	}

	for _, t := range pkg.Types {
		typeSymbols(sym, t)
	}

	return sym
}

func typeSymbols(sym map[string]Symbol, t *doc.Type) {
	typeSym := Symbol{
		Name: t.Name,
		Kind: TypeSymbolKind,
	}

	sym[t.Name] = typeSym

	for _, f := range t.Methods {
		sym[symbolName(t.Name, f.Name)] = Symbol{
			Receiver: t.Name,
			Name:     f.Name,
			Kind:     MethodSymbolKind,
		}
	}

	for _, s := range t.Decl.Specs {
		typ, ok := s.(*ast.TypeSpec).Type.(*ast.StructType)
		if !ok {
			continue
		}

		for _, f := range typ.Fields.List {
			for _, n := range f.Names {
				sym[symbolName(t.Name, n.String())] = Symbol{
					Receiver: t.Name,
					Name:     n.String(),
					Kind:     FieldSymbolKind,
					Parent:   &typeSym,
				}
			}
		}
	}

	for _, f := range t.Funcs {
		sym[symbolName("", f.Name)] = Symbol{
			Name: f.Name,
			Kind: FuncSymbolKind,
		}
	}

	for _, c := range t.Consts {
		parent := Symbol{
			Name: c.Names[0],
			Kind: ConstSymbolKind,
		}

		for _, n := range c.Names {
			sym[symbolName("", n)] = Symbol{
				Name:   n,
				Kind:   ConstSymbolKind,
				Parent: &parent,
			}
		}
	}

	for _, v := range t.Vars {
		parent := Symbol{
			Name: v.Names[0],
			Kind: VarSymbolKind,
		}

		for _, n := range v.Names {
			sym[symbolName("", n)] = Symbol{
				Name:   n,
				Kind:   VarSymbolKind,
				Parent: &parent,
			}
		}
	}

}

// Anchor produces anchor text for the symbol.
func (s Symbol) Anchor() string {
	if s.Parent != nil {
		return s.Parent.Anchor()
	}

	switch s.Kind {
	case MethodSymbolKind, FieldSymbolKind:
		return fmt.Sprintf("%s.%s", strings.TrimLeft(s.Receiver, "*"), s.Name)
	default:
		return s.Name
	}
}

// symbolName returns the string representation of the symbol.
func symbolName(receiver string, name string) string {
	receiver = strings.TrimLeft(receiver, "*")
	name = strings.TrimLeft(name, "*")

	if receiver == "" {
		return name
	}

	return fmt.Sprintf("%s.%s", receiver, name)
}
