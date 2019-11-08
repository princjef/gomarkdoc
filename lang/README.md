# lang

```go
import "github.com/princjef/gomarkdoc/lang"
```

Package lang provides constructs for defining golang language constructs and extracting information from them for documentation purposes\.

## Index

- [type Block](<#type-block>)
  - [func NewBlock(kind BlockKind, text string, level int) *Block](<#func-newblock>)
  - [func (b *Block) Kind() BlockKind](<#func-block-kind>)
  - [func (b *Block) Level() int](<#func-block-level>)
  - [func (b *Block) Text() string](<#func-block-text>)
- [type BlockKind](<#type-blockkind>)
- [type Doc](<#type-doc>)
  - [func NewDoc(text string, level int) *Doc](<#func-newdoc>)
  - [func (d *Doc) Blocks() []*Block](<#func-doc-blocks>)
  - [func (d *Doc) Level() int](<#func-doc-level>)
- [type Example](<#type-example>)
  - [func NewExample(name string, doc *doc.Example, fs *token.FileSet, level int) *Example](<#func-newexample>)
  - [func (ex *Example) Code() (string, error)](<#func-example-code>)
  - [func (ex *Example) Doc() *Doc](<#func-example-doc>)
  - [func (ex *Example) Level() int](<#func-example-level>)
  - [func (ex *Example) Name() string](<#func-example-name>)
  - [func (ex *Example) Summary() string](<#func-example-summary>)
  - [func (ex *Example) Title() string](<#func-example-title>)
- [type Func](<#type-func>)
  - [func NewFunc(doc *doc.Func, fs *token.FileSet, examples []*doc.Example, level int) *Func](<#func-newfunc>)
  - [func (fn *Func) Doc() *Doc](<#func-func-doc>)
  - [func (fn *Func) Examples() (examples []*Example)](<#func-func-examples>)
  - [func (fn *Func) Level() int](<#func-func-level>)
  - [func (fn *Func) Name() string](<#func-func-name>)
  - [func (fn *Func) Signature() (string, error)](<#func-func-signature>)
  - [func (fn *Func) Summary() string](<#func-func-summary>)
  - [func (fn *Func) Title() string](<#func-func-title>)
- [type Package](<#type-package>)
  - [func NewPackage(doc *doc.Package, fs *token.FileSet, examples []*doc.Example, dir string) (*Package, error)](<#func-newpackage>)
  - [func NewPackageFromBuild(pkg *build.Package, opts ...PackageOption) (*Package, error)](<#func-newpackagefrombuild>)
  - [func (pkg *Package) Consts() (consts []*Value)](<#func-package-consts>)
  - [func (pkg *Package) Dir() string](<#func-package-dir>)
  - [func (pkg *Package) Dirname() string](<#func-package-dirname>)
  - [func (pkg *Package) Doc() *Doc](<#func-package-doc>)
  - [func (pkg *Package) Examples() (examples []*Example)](<#func-package-examples>)
  - [func (pkg *Package) Funcs() (funcs []*Func)](<#func-package-funcs>)
  - [func (pkg *Package) Import() string](<#func-package-import>)
  - [func (pkg *Package) Level() int](<#func-package-level>)
  - [func (pkg *Package) Name() string](<#func-package-name>)
  - [func (pkg *Package) Summary() string](<#func-package-summary>)
  - [func (pkg *Package) Types() (types []*Type)](<#func-package-types>)
  - [func (pkg *Package) Vars() (vars []*Value)](<#func-package-vars>)
- [type PackageOption](<#type-packageoption>)
  - [func PackageWithUnexportedIncluded() PackageOption](<#func-packagewithunexportedincluded>)
- [type PackageOptions](<#type-packageoptions>)
- [type Type](<#type-type>)
  - [func NewType(doc *doc.Type, fs *token.FileSet, examples []*doc.Example, level int) *Type](<#func-newtype>)
  - [func (typ *Type) Decl() (string, error)](<#func-type-decl>)
  - [func (typ *Type) Doc() *Doc](<#func-type-doc>)
  - [func (typ *Type) Examples() (examples []*Example)](<#func-type-examples>)
  - [func (typ *Type) Funcs() (funcs []*Func)](<#func-type-funcs>)
  - [func (typ *Type) Level() int](<#func-type-level>)
  - [func (typ *Type) Methods() (methods []*Func)](<#func-type-methods>)
  - [func (typ *Type) Name() string](<#func-type-name>)
  - [func (typ *Type) Summary() string](<#func-type-summary>)
  - [func (typ *Type) Title() string](<#func-type-title>)
- [type Value](<#type-value>)
  - [func NewValue(doc *doc.Value, fs *token.FileSet, level int) *Value](<#func-newvalue>)
  - [func (v *Value) Decl() (string, error)](<#func-value-decl>)
  - [func (v *Value) Doc() *Doc](<#func-value-doc>)
  - [func (v *Value) Level() int](<#func-value-level>)
  - [func (v *Value) Summary() string](<#func-value-summary>)


## type Block

Block defines a single block element \(e\.g\. paragraph\, code block\) in the documentation for a symbol or package\.

```go
type Block struct {
    // contains filtered or unexported fields
}
```

### func NewBlock

```go
func NewBlock(kind BlockKind, text string, level int) *Block
```

NewBlock creates a new block element of the provided kind and with the given text contents\.

### func \(\*Block\) Kind

```go
func (b *Block) Kind() BlockKind
```

Kind provides the kind of data that this block's text should be interpreted as\.

### func \(\*Block\) Level

```go
func (b *Block) Level() int
```

Level provides the default level that a block of kind HeaderBlock will render at in the output\. The level is not used for other block types\.

### func \(\*Block\) Text

```go
func (b *Block) Text() string
```

Text provides the raw text of the block's contents\. The text is pre\-scrubbed and sanitized as determined by the block's Kind\(\)\, but it is not wrapped in any special constructs for rendering purposes \(such as markdown code blocks\)\.

## type BlockKind

BlockKind identifies the type of block element represented by the corresponding Block\.

```go
type BlockKind string
```

## type Doc

Doc provides access to the documentation comment contents for a package or symbol in a structured form\.

```go
type Doc struct {
    // contains filtered or unexported fields
}
```

### func NewDoc

```go
func NewDoc(text string, level int) *Doc
```

NewDoc initializes a Doc struct from the provided raw documentation text and with headers rendered by default at the heading level provided\. Documentation is separated into block level elements using the standard rules from golang's documentation conventions\.

### func \(\*Doc\) Blocks

```go
func (d *Doc) Blocks() []*Block
```

Blocks holds the list of block elements that makes up the documentation contents\.

### func \(\*Doc\) Level

```go
func (d *Doc) Level() int
```

Level provides the default level that headers within the documentation should be rendered

## type Example

Example holds a single documentation example for a package or symbol\.

```go
type Example struct {
    // contains filtered or unexported fields
}
```

### func NewExample

```go
func NewExample(name string, doc *doc.Example, fs *token.FileSet, level int) *Example
```

NewExample creates a new example from the example function's name\, its documentation example and the files holding code related to the example\.

### func \(\*Example\) Code

```go
func (ex *Example) Code() (string, error)
```

Code provides the raw text code representation of the example's contents\.

### func \(\*Example\) Doc

```go
func (ex *Example) Doc() *Doc
```

Doc provides the structured contents of the documentation comment for the example\.

### func \(\*Example\) Level

```go
func (ex *Example) Level() int
```

Level provides the default level that headers for the example should be rendered\.

### func \(\*Example\) Name

```go
func (ex *Example) Name() string
```

Name provides a pretty\-printed name for the specific example\, if one was provided\.

### func \(\*Example\) Summary

```go
func (ex *Example) Summary() string
```

Summary provides the one\-sentence summary of the example's documentation comment\.

### func \(\*Example\) Title

```go
func (ex *Example) Title() string
```

Title provides a formatted string to print as the title of the example\. It incorporates the example's name\, if present\.

## type Func

Func holds documentation information for a single func declaration within a package or type\.

```go
type Func struct {
    // contains filtered or unexported fields
}
```

### func NewFunc

```go
func NewFunc(doc *doc.Func, fs *token.FileSet, examples []*doc.Example, level int) *Func
```

NewFunc creates a new Func from the corresponding documentation construct from the standard library\, the related token\.FileSet for the package and the list of examples for the package\.

### func \(\*Func\) Doc

```go
func (fn *Func) Doc() *Doc
```

Doc provides the structured contents of the documentation comment for the function\.

### func \(\*Func\) Examples

```go
func (fn *Func) Examples() (examples []*Example)
```

Examples provides the list of examples from the list given on initialization that pertain to the function\.

### func \(\*Func\) Level

```go
func (fn *Func) Level() int
```

Level provides the default level at which headers for the func should be rendered in the final documentation\.

### func \(\*Func\) Name

```go
func (fn *Func) Name() string
```

Name provides the name of the function\.

### func \(\*Func\) Signature

```go
func (fn *Func) Signature() (string, error)
```

Signature provides the raw text representation of the code for the function's signature\.

### func \(\*Func\) Summary

```go
func (fn *Func) Summary() string
```

Summary provides the one\-sentence summary of the function's documentation comment

### func \(\*Func\) Title

```go
func (fn *Func) Title() string
```

Title provides the formatted name of the func\. It is primarily designed for generating headers\.

## type Package

Package holds documentation information for a package and all of the symbols contained within it\.

```go
type Package struct {
    // contains filtered or unexported fields
}
```

### func NewPackage

```go
func NewPackage(doc *doc.Package, fs *token.FileSet, examples []*doc.Example, dir string) (*Package, error)
```

NewPackage creates a representation of a package's documentation from the raw documentation constructs provided by the standard library\. This is only recommended for advanced scenarios\. Most consumers will find it easier to use NewPackageFromBuild instead\.

### func NewPackageFromBuild

```go
func NewPackageFromBuild(pkg *build.Package, opts ...PackageOption) (*Package, error)
```

NewPackageFromBuild creates a representation of a package's documentation from the build metadata for that package\. It can be configured using the provided options\.

### func \(\*Package\) Consts

```go
func (pkg *Package) Consts() (consts []*Value)
```

Consts lists the top\-level constants provided by the package\.

### func \(\*Package\) Dir

```go
func (pkg *Package) Dir() string
```

Dir provides the name of the full directory in which the package is located\.

### func \(\*Package\) Dirname

```go
func (pkg *Package) Dirname() string
```

Dirname provides the name of the leaf directory in which the package is located\.

### func \(\*Package\) Doc

```go
func (pkg *Package) Doc() *Doc
```

Doc provides the structured contents of the documentation comment for the package\.

### func \(\*Package\) Examples

```go
func (pkg *Package) Examples() (examples []*Example)
```

Examples provides the package\-level examples that have been defined\. This does not include examples that are associated with symbols contained within the package\.

### func \(\*Package\) Funcs

```go
func (pkg *Package) Funcs() (funcs []*Func)
```

Funcs lists the top\-level functions provided by the package\.

### func \(\*Package\) Import

```go
func (pkg *Package) Import() string
```

Import provides the raw text for the import declaration that is used to import code from the package\. If your package's documentation is generated from a local path and does not use Go Modules\, this will typically print \`import "\."\`\.

### func \(\*Package\) Level

```go
func (pkg *Package) Level() int
```

Level provides the default level that headers for the package's root documentation should be rendered\.

### func \(\*Package\) Name

```go
func (pkg *Package) Name() string
```

Name provides the name of the package as it would be seen from another package importing it\.

### func \(\*Package\) Summary

```go
func (pkg *Package) Summary() string
```

Summary provides the one\-sentence summary of the package's documentation comment\.

### func \(\*Package\) Types

```go
func (pkg *Package) Types() (types []*Type)
```

Types lists the top\-level types provided by the package\.

### func \(\*Package\) Vars

```go
func (pkg *Package) Vars() (vars []*Value)
```

Vars lists the top\-level variables provided by the package\.

## type PackageOption

PackageOption configures one or more options for the package\.

```go
type PackageOption func(opts *PackageOptions) error
```

### func PackageWithUnexportedIncluded

```go
func PackageWithUnexportedIncluded() PackageOption
```

PackageWithUnexportedIncluded can be used along with the NewPackageFromBuild function to specify that all symbols\, including unexported ones\, should be included in the documentation for the package\.

## type PackageOptions

PackageOptions holds options related to the configuration of the package and its documentation on creation\.

```go
type PackageOptions struct {
    // contains filtered or unexported fields
}
```

## type Type

Type holds documentation information for a type declaration\.

```go
type Type struct {
    // contains filtered or unexported fields
}
```

### func NewType

```go
func NewType(doc *doc.Type, fs *token.FileSet, examples []*doc.Example, level int) *Type
```

NewType creates a Type from the raw documentation representation of the type\, the token\.FileSet for the package's files and the full list of examples from the containing package\.

### func \(\*Type\) Decl

```go
func (typ *Type) Decl() (string, error)
```

Decl provides the raw text representation of the code for the type's declaration\.

### func \(\*Type\) Doc

```go
func (typ *Type) Doc() *Doc
```

Doc provides the structured contents of the documentation comment for the type\.

### func \(\*Type\) Examples

```go
func (typ *Type) Examples() (examples []*Example)
```

Examples lists the examples pertaining to the type from the set provided on initialization\.

### func \(\*Type\) Funcs

```go
func (typ *Type) Funcs() (funcs []*Func)
```

Funcs lists the funcs related to the type\. This only includes functions which return an instance of the type or its pointer\.

### func \(\*Type\) Level

```go
func (typ *Type) Level() int
```

Level provides the default level that headers for the type should be rendered\.

### func \(\*Type\) Methods

```go
func (typ *Type) Methods() (methods []*Func)
```

Methods lists the funcs that use the type as a value or pointer receiver\.

### func \(\*Type\) Name

```go
func (typ *Type) Name() string
```

Name provides the name of the type

### func \(\*Type\) Summary

```go
func (typ *Type) Summary() string
```

Summary provides the one\-sentence summary of the type's documentation comment\.

### func \(\*Type\) Title

```go
func (typ *Type) Title() string
```

Title provides a formatted name suitable for use in a header identifying the type\.

## type Value

Value holds documentation for a var or const declaration within a package\.

```go
type Value struct {
    // contains filtered or unexported fields
}
```

### func NewValue

```go
func NewValue(doc *doc.Value, fs *token.FileSet, level int) *Value
```

NewValue creates a new Value from the raw const or var documentation and the token\.FileSet of files for the containing package\.

### func \(\*Value\) Decl

```go
func (v *Value) Decl() (string, error)
```

Decl provides the raw text representation of the code for declaring the const or var\.

### func \(\*Value\) Doc

```go
func (v *Value) Doc() *Doc
```

Doc provides the structured contents of the documentation comment for the example\.

### func \(\*Value\) Level

```go
func (v *Value) Level() int
```

Level provides the default level that headers for the value should be rendered\.

### func \(\*Value\) Summary

```go
func (v *Value) Summary() string
```

Summary provides the one\-sentence summary of the value's documentation comment\.

