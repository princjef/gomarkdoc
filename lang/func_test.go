package lang_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

func TestFunc_stringsCompare(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("strings")
	is.NoErr(err)

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	is.NoErr(err)

	var fn *lang.Func
	for _, f := range pkg.Funcs() {
		if f.Name() == "Compare" {
			fn = f
			break
		}
	}

	is.True(fn != nil) // didn't find the function we were looking for

	sig, err := fn.Signature()
	is.NoErr(err)

	is.Equal(fn.Name(), "Compare")
	is.Equal(fn.Level(), 2)
	is.Equal(fn.Title(), "func Compare")
	is.Equal(fn.Summary(), "Compare returns an integer comparing two strings lexicographically.")
	is.Equal(sig, "func Compare(a, b string) int")
	is.Equal(len(fn.Examples()), 1)
}

func TestFunc_textScannerInit(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("text/scanner")
	is.NoErr(err)

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	is.NoErr(err)

	var typ *lang.Type
	for _, t := range pkg.Types() {
		if t.Name() == "Scanner" {
			typ = t
			break
		}
	}
	is.True(typ != nil)

	var fn *lang.Func
	for _, f := range typ.Methods() {
		if f.Name() == "Init" {
			fn = f
			break
		}
	}

	is.True(fn != nil) // didn't find the function we were looking for

	sig, err := fn.Signature()
	is.NoErr(err)

	is.Equal(fn.Name(), "Init")
	is.Equal(fn.Level(), 3)
	is.Equal(fn.Title(), "func (*Scanner) Init")
	is.Equal(fn.Summary(), "Init initializes a Scanner with a new source and returns s.")
	is.Equal(sig, "func (s *Scanner) Init(src io.Reader) *Scanner")
	is.Equal(len(fn.Examples()), 0)
}

func TestFunc_ioIoutilTempFile(t *testing.T) {
	is := is.New(t)

	buildPkg, err := getBuildPackage("io/ioutil")
	is.NoErr(err)

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	is.NoErr(err)

	var fn *lang.Func
	for _, f := range pkg.Funcs() {
		if f.Name() == "TempFile" {
			fn = f
			break
		}
	}

	is.True(fn != nil) // didn't find the function we were looking for

	sig, err := fn.Signature()
	is.NoErr(err)

	is.Equal(fn.Name(), "TempFile")
	is.Equal(fn.Level(), 2)
	is.Equal(fn.Title(), "func TempFile")
	is.Equal(fn.Summary(), "TempFile creates a new temporary file in the directory dir, opens the file for reading and writing, and returns the resulting *os.File.")
	is.Equal(sig, "func TempFile(dir, pattern string) (f *os.File, err error)")
	is.Equal(len(fn.Examples()), 2)
}
