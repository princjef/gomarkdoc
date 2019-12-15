package lang_test

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/lang"
)

func TestPackage_strings(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("strings")
	is.NoErr(err)

	pkg, err := lang.NewPackageFromBuild(buildPkg)
	is.NoErr(err)

	is.Equal(pkg.Level(), 1) // level should be root
	is.True(strings.HasSuffix(pkg.Dir(), filepath.FromSlash("/strings")))
	is.Equal(pkg.Dirname(), "strings")
	is.Equal(pkg.Name(), "strings")
	is.Equal(pkg.Import(), `import "strings"`)
	is.Equal(pkg.Summary(), "Package strings implements simple functions to manipulate UTF-8 encoded strings.")
	is.Equal(len(pkg.Consts()), 0)   // strings should have no constants
	is.Equal(len(pkg.Vars()), 0)     // strings should have no vars
	is.True(len(pkg.Funcs()) > 0)    // strings should have top-level functions
	is.True(len(pkg.Types()) > 0)    // strings should have top-level types
	is.Equal(len(pkg.Examples()), 0) // strings should have no top-level examples
}

func TestPackage_textScanner(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("text/scanner")
	is.NoErr(err)

	pkg, err := lang.NewPackageFromBuild(buildPkg)
	is.NoErr(err)

	is.Equal(pkg.Level(), 1) // level should be root
	is.True(strings.HasSuffix(pkg.Dir(), filepath.FromSlash("/text/scanner")))
	is.Equal(pkg.Dirname(), "scanner")
	is.Equal(pkg.Name(), "scanner")
	is.Equal(pkg.Import(), `import "text/scanner"`)
	is.Equal(pkg.Summary(), "Package scanner provides a scanner and tokenizer for UTF-8-encoded text.")
	is.True(len(pkg.Consts()) > 0)   // text/scanner should have constants
	is.Equal(len(pkg.Vars()), 0)     // text/scanner should have no vars
	is.True(len(pkg.Funcs()) > 0)    // text/scanner should have top-level functions
	is.True(len(pkg.Types()) > 0)    // text/scanner should have top-level types
	is.True(len(pkg.Examples()) > 0) // text/scanner should have top-level examples
}

func TestPackage_ioIoutil(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("io/ioutil")
	is.NoErr(err)

	pkg, err := lang.NewPackageFromBuild(buildPkg)
	is.NoErr(err)

	is.Equal(pkg.Level(), 1) // level should be root
	is.True(strings.HasSuffix(pkg.Dir(), filepath.FromSlash("/io/ioutil")))
	is.Equal(pkg.Dirname(), "ioutil")
	is.Equal(pkg.Name(), "ioutil")
	is.Equal(pkg.Import(), `import "io/ioutil"`)
	is.Equal(pkg.Summary(), "Package ioutil implements some I/O utility functions.")
	is.Equal(len(pkg.Consts()), 0)   // io/ioutil should have no constants
	is.True(len(pkg.Vars()) > 0)     // io/ioutil should have vars
	is.True(len(pkg.Funcs()) > 0)    // io/ioutil should have top-level functions
	is.Equal(len(pkg.Types()), 0)    // io/ioutil should have no top-level types
	is.Equal(len(pkg.Examples()), 0) // io/ioutil should have no top-level examples
}

func TestPackage_encoding(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("encoding")
	is.NoErr(err)

	pkg, err := lang.NewPackageFromBuild(buildPkg)
	is.NoErr(err)

	is.Equal(pkg.Level(), 1) // level should be root
	is.True(strings.HasSuffix(pkg.Dir(), filepath.FromSlash("/encoding")))
	is.Equal(pkg.Dirname(), "encoding")
	is.Equal(pkg.Name(), "encoding")
	is.Equal(pkg.Import(), `import "encoding"`)
	is.Equal(pkg.Summary(), "Package encoding defines interfaces shared by other packages that convert data to and from byte-level and textual representations.")
	is.Equal(len(pkg.Consts()), 0)   // encoding should have no constants
	is.Equal(len(pkg.Vars()), 0)     // encoding should have no vars
	is.Equal(len(pkg.Funcs()), 0)    // encoding should have no top-level functions
	is.True(len(pkg.Types()) > 0)    // encoding should have top-level types
	is.Equal(len(pkg.Examples()), 0) // encoding should have no top-level examples
}

func getBuildPackage(path string) (*build.Package, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return build.Import(path, wd, build.ImportComment)
}
