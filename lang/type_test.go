package lang_test

import (
	"errors"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
)

func TestType_Examples(t *testing.T) {
	is := is.New(t)

	typ, err := loadType("../testData/lang/function", "Receiver")
	is.NoErr(err)

	ex := typ.Examples()
	is.Equal(len(ex), 2)

	is.Equal(ex[0].Name(), "")
	is.Equal(ex[1].Name(), "Sub Test")
}

func loadType(dir, name string) (*lang.Type, error) {
	buildPkg, err := getBuildPackage(dir)
	if err != nil {
		return nil, err
	}

	log := logger.New(logger.ErrorLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	if err != nil {
		return nil, err
	}

	for _, t := range pkg.Types() {
		if t.Name() == name {
			return t, nil
		}
	}

	return nil, errors.New("type not found")
}
