package lang_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

func TestValue_Level(t *testing.T) {
	is := is.New(t)

	val, err := loadValue("../testData/lang/function", "Variable")
	is.NoErr(err)

	is.Equal(val.Level(), 2)
}

func TestValue_Summary(t *testing.T) {
	is := is.New(t)

	val, err := loadValue("../testData/lang/function", "Variable")
	is.NoErr(err)

	is.Equal(val.Summary(), "Variable is a package-level variable.")
}

func TestValue_Doc(t *testing.T) {
	is := is.New(t)

	val, err := loadValue("../testData/lang/function", "Variable")
	is.NoErr(err)

	doc := val.Doc()
	is.Equal(doc.Level(), 3)
	is.Equal(len(doc.Blocks()), 1)
	is.Equal(doc.Blocks()[0].Kind(), lang.ParagraphBlock)
	is.Equal(doc.Blocks()[0].Level(), 3)
	is.Equal(doc.Blocks()[0].Text(), "Variable is a package-level variable.")
}

func TestValue_Decl(t *testing.T) {
	is := is.New(t)

	val, err := loadValue("../testData/lang/function", "Variable")
	is.NoErr(err)

	decl, err := val.Decl()
	is.NoErr(err)

	is.Equal(decl, "var Variable = 5")
}

func TestValue_Location(t *testing.T) {
	is := is.New(t)

	val, err := loadValue("../testData/lang/function", "Variable")
	is.NoErr(err)

	loc := val.Location()
	is.Equal(loc.Start.Line, 4)
	is.Equal(loc.Start.Col, 1)
	is.Equal(loc.End.Line, 4)
	is.Equal(loc.End.Col, 17)
	is.True(strings.HasSuffix(loc.Filepath, "value.go"))
}

func loadValue(dir, name string) (*lang.Value, error) {
	buildPkg, err := getBuildPackage(dir)
	if err != nil {
		return nil, err
	}

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return nil, err
	}

	for _, v := range pkg.Vars() {
		d, err := v.Decl()
		if err == nil && strings.Contains(d, name) {
			return v, nil
		}
	}

	for _, v := range pkg.Consts() {
		d, err := v.Decl()
		if err == nil && strings.Contains(d, name) {
			return v, nil
		}
	}

	return nil, errors.New("value not found")
}
