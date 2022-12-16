package lang_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

func TestFunc_Level_standalone(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "Standalone")
	is.NoErr(err)

	is.Equal(fn.Level(), 2)
}

func TestFunc_Level_receiver(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "WithReceiver")
	is.NoErr(err)

	is.Equal(fn.Level(), 3)
}

func TestFunc_Level_initializer(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "New")
	is.NoErr(err)

	is.Equal(fn.Level(), 3)
}

func TestFunc_Name_standalone(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "Standalone")
	is.NoErr(err)

	is.Equal(fn.Name(), "Standalone")
}

func TestFunc_Name_receiver(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "WithReceiver")
	is.NoErr(err)

	is.Equal(fn.Name(), "WithReceiver")
}

func TestFunc_Receiver_standalone(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "Standalone")
	is.NoErr(err)

	is.Equal(fn.Receiver(), "")
}

func TestFunc_Receiver_receiver(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "WithReceiver")
	is.NoErr(err)

	is.Equal(fn.Receiver(), "Receiver")
}

func TestFunc_Doc(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "Standalone")
	is.NoErr(err)

	doc := fn.Doc()
	blocks := doc.Blocks()
	is.Equal(len(blocks), 4)

	is.Equal(blocks[0].Kind(), lang.ParagraphBlock)
	is.Equal(blocks[0].Level(), 3)
	is.Equal(blocks[0].Text(), "Standalone provides a function that is not part of a type.\n\nAdditional description can be provided in subsequent paragraphs, including code blocks and headers")

	is.Equal(blocks[1].Kind(), lang.HeaderBlock)
	is.Equal(blocks[1].Level(), 3)
	is.Equal(blocks[1].Text(), "Header A")

	is.Equal(blocks[2].Kind(), lang.ParagraphBlock)
	is.Equal(blocks[2].Level(), 3)
	fmt.Printf("%q\n", blocks[2].Text())
	is.Equal(blocks[2].Text(), "This section contains a code block.")

	is.Equal(blocks[3].Kind(), lang.CodeBlock)
	is.Equal(blocks[3].Level(), 3)
	is.Equal(blocks[3].Text(), "\tCode Block\n\tMore of Code Block")
}

func TestFunc_Location(t *testing.T) {
	is := is.New(t)

	fn, err := loadFunc("../testData/lang/function", "Standalone")
	is.NoErr(err)

	loc := fn.Location()
	is.Equal(loc.Start.Line, 14)
	is.Equal(loc.Start.Col, 1)
	is.Equal(loc.End.Line, 14)
	is.Equal(loc.End.Col, 48)
	is.True(strings.HasSuffix(loc.Filepath, "func.go"))
}

func TestFunc_Examples_generic(t *testing.T) {
	is := is.New(t)
	fn, err := loadFunc("../testData/lang/function", "WithGenericReceiver")
	is.NoErr(err)

	examples := fn.Examples()
	is.Equal(len(examples), 1)

	ex := examples[0]
	is.Equal(ex.Name(), "")
}

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

func loadFunc(dir, name string) (*lang.Func, error) {
	buildPkg, err := getBuildPackage(dir)
	if err != nil {
		return nil, err
	}

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return nil, err
	}

	for _, f := range pkg.Funcs() {
		if f.Name() == name {
			return f, nil
		}
	}

	for _, t := range pkg.Types() {
		for _, f := range t.Funcs() {
			if f.Name() == name {
				return f, nil
			}
		}

		for _, f := range t.Methods() {
			if f.Name() == name {
				return f, nil
			}
		}
	}

	return nil, errors.New("func not found")
}
