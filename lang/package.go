package lang

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/princjef/gomarkdoc/logger"
)

type (
	// Package holds documentation information for a package and all of the
	// symbols contained within it.
	Package struct {
		cfg      *Config
		doc      *doc.Package
		examples []*doc.Example
	}

	// PackageOptions holds options related to the configuration of the package
	// and its documentation on creation.
	PackageOptions struct {
		includeUnexported   bool
		repositoryOverrides *Repo
	}

	// PackageOption configures one or more options for the package.
	PackageOption func(opts *PackageOptions) error
)

// NewPackage creates a representation of a package's documentation from the
// raw documentation constructs provided by the standard library. This is only
// recommended for advanced scenarios. Most consumers will find it easier to use
// NewPackageFromBuild instead.
func NewPackage(cfg *Config, doc *doc.Package, examples []*doc.Example) *Package {
	return &Package{cfg, doc, examples}
}

// NewPackageFromBuild creates a representation of a package's documentation
// from the build metadata for that package. It can be configured using the
// provided options.
func NewPackageFromBuild(log logger.Logger, pkg *build.Package, opts ...PackageOption) (*Package, error) {
	var options PackageOptions
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfig(log, wd, pkg.Dir, ConfigWithRepoOverrides(options.repositoryOverrides))
	if err != nil {
		return nil, err
	}

	docPkg, err := getDocPkg(pkg, cfg.FileSet, options.includeUnexported)
	if err != nil {
		return nil, err
	}

	files, err := parsePkgFiles(pkg, cfg.FileSet)
	if err != nil {
		return nil, err
	}

	examples := doc.Examples(files...)

	return NewPackage(cfg, docPkg, examples), nil
}

// PackageWithUnexportedIncluded can be used along with the NewPackageFromBuild
// function to specify that all symbols, including unexported ones, should be
// included in the documentation for the package.
func PackageWithUnexportedIncluded() PackageOption {
	return func(opts *PackageOptions) error {
		opts.includeUnexported = true
		return nil
	}
}

// PackageWithRepositoryOverrides can be used along with the NewPackageFromBuild
// function to define manual overrides to the automatic repository detection
// logic.
func PackageWithRepositoryOverrides(repo *Repo) PackageOption {
	return func(opts *PackageOptions) error {
		opts.repositoryOverrides = repo
		return nil
	}
}

// Level provides the default level that headers for the package's root
// documentation should be rendered.
func (pkg *Package) Level() int {
	return pkg.cfg.Level
}

// Dir provides the name of the full directory in which the package is located.
func (pkg *Package) Dir() string {
	return pkg.cfg.PkgDir
}

// Dirname provides the name of the leaf directory in which the package is
// located.
func (pkg *Package) Dirname() string {
	return filepath.Base(pkg.cfg.PkgDir)
}

// Name provides the name of the package as it would be seen from another
// package importing it.
func (pkg *Package) Name() string {
	return pkg.doc.Name
}

// Import provides the raw text for the import declaration that is used to
// import code from the package. If your package's documentation is generated
// from a local path and does not use Go Modules, this will typically print
// `import "."`.
func (pkg *Package) Import() string {
	return fmt.Sprintf(`import "%s"`, pkg.doc.ImportPath)
}

// Summary provides the one-sentence summary of the package's documentation
// comment.
func (pkg *Package) Summary() string {
	return extractSummary(pkg.doc.Doc)
}

// Doc provides the structured contents of the documentation comment for the
// package.
func (pkg *Package) Doc() *Doc {
	// TODO: level should only be + 1, but we have special knowledge for rendering
	return NewDoc(pkg.cfg.Inc(2), pkg.doc.Doc)
}

// Consts lists the top-level constants provided by the package.
func (pkg *Package) Consts() (consts []*Value) {
	for _, c := range pkg.doc.Consts {
		consts = append(consts, NewValue(pkg.cfg.Inc(1), c))
	}

	return
}

// Vars lists the top-level variables provided by the package.
func (pkg *Package) Vars() (vars []*Value) {
	for _, v := range pkg.doc.Vars {
		vars = append(vars, NewValue(pkg.cfg.Inc(1), v))
	}

	return
}

// Funcs lists the top-level functions provided by the package.
func (pkg *Package) Funcs() (funcs []*Func) {
	for _, fn := range pkg.doc.Funcs {
		funcs = append(funcs, NewFunc(pkg.cfg.Inc(1), fn, pkg.examples))
	}

	return
}

// Types lists the top-level types provided by the package.
func (pkg *Package) Types() (types []*Type) {
	for _, typ := range pkg.doc.Types {
		types = append(types, NewType(pkg.cfg.Inc(1), typ, pkg.examples))
	}

	return
}

// Examples provides the package-level examples that have been defined. This
// does not include examples that are associated with symbols contained within
// the package.
func (pkg *Package) Examples() (examples []*Example) {
	for _, example := range pkg.examples {
		var name string
		switch {
		case example.Name == "":
			name = ""
		case strings.HasPrefix(example.Name, "_"):
			name = example.Name[1:]
		default:
			// TODO: better filtering
			continue
		}

		examples = append(examples, NewExample(pkg.cfg.Inc(1), name, example))
	}

	return
}

var goModRegex = regexp.MustCompile(`^\s*module ([^\s]+)`)

// findImportPath attempts to find an import path for the contents of the
// provided dir by walking up to the nearest go.mod file and constructing an
// import path from it. If the directory is not in a Go Module, the second
// return value will be false.
func findImportPath(dir string) (string, bool) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}

	f, ok := findFileInParent(absDir, "go.mod", false)
	if !ok {
		return "", false
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", false
	}

	m := goModRegex.FindSubmatch(b)
	if m == nil {
		return "", false
	}

	relative, err := filepath.Rel(filepath.Dir(f.Name()), absDir)
	if err != nil {
		return "", false
	}

	relative = filepath.ToSlash(relative)

	return path.Join(string(m[1]), relative), true
}

// findFileInParent looks for a file or directory of the given name within the
// provided dir. The returned os.File is opened and must be closed by the
// caller to avoid a memory leak.
func findFileInParent(dir, filename string, fileIsDir bool) (*os.File, bool) {
	initial := dir
	current := initial

	for {
		p := filepath.Join(current, filename)
		if f, err := os.Open(p); err == nil {
			if s, err := f.Stat(); err == nil && (fileIsDir && s.Mode().IsDir() || !fileIsDir && s.Mode().IsRegular()) {
				return f, true
			}
		}

		// Walk up a dir
		next := filepath.Join(current, "..")

		// If we didn't change dirs, there's no more to search
		if current == next {
			break
		}

		current = next
	}

	return nil, false
}

func getDocPkg(pkg *build.Package, fs *token.FileSet, includeUnexported bool) (*doc.Package, error) {
	pkgs, err := parser.ParseDir(
		fs,
		pkg.Dir,
		func(info os.FileInfo) bool {
			for _, name := range pkg.GoFiles {
				if name == info.Name() {
					return true
				}
			}

			for _, name := range pkg.CgoFiles {
				if name == info.Name() {
					return true
				}
			}

			return false
		},
		parser.ParseComments,
	)

	if err != nil {
		return nil, fmt.Errorf("gomarkdoc: failed to parse package: %w", err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("gomarkdoc: no source-code package in directory %s", pkg.Dir)
	}

	if len(pkgs) > 1 {
		return nil, fmt.Errorf("gomarkdoc: multiple packages in directory %s", pkg.Dir)
	}

	astPkg := pkgs[pkg.Name]

	if !includeUnexported {
		ast.PackageExports(astPkg)
	}

	importPath := pkg.ImportPath
	if pkg.ImportComment != "" {
		importPath = pkg.ImportComment
	}

	if importPath == "." {
		if modPath, ok := findImportPath(pkg.Dir); ok {
			importPath = modPath
		}
	}

	return doc.New(astPkg, importPath, doc.AllDecls), nil
}

func parsePkgFiles(pkg *build.Package, fs *token.FileSet) ([]*ast.File, error) {
	rawFiles, err := ioutil.ReadDir(pkg.Dir)
	if err != nil {
		return nil, fmt.Errorf("gomarkdoc: error reading package dir: %w", err)
	}

	var files []*ast.File
	for _, f := range rawFiles {
		if !strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), ".cgo") {
			continue
		}

		p := path.Join(pkg.Dir, f.Name())

		fi, err := os.Stat(p)
		if err != nil || !fi.Mode().IsRegular() {
			continue
		}

		parsed, err := parser.ParseFile(fs, p, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("gomarkdoc: failed to parse package file %s", f.Name())
		}

		files = append(files, parsed)
	}

	return files, nil
}
