package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var wd, _ = os.Getwd()

func TestCommand(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./simple",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("simple")

	main()

	verify(t, "simple")
}

func TestCommand_check(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./simple",
		"-c",
		"-o", "{{.Dir}}/README.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("simple")

	main()
}

func TestCommand_nested(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./nested/...",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("nested")
	cleanup("nested/inner")

	main()

	verify(t, "nested")
	verify(t, "nested/inner")
}

func TestCommand_unexported(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./unexported",
		"-u",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("unexported")

	main()

	verify(t, "unexported")
}

func TestCommand_version(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{"gomarkdoc", "--version"}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	main()
	w.Close()

	data, err := io.ReadAll(r)
	is.NoErr(err)

	is.Equal(strings.TrimSpace(string(data)), "(devel)")
}

func TestCommand_invalidCheck(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./simple",
		"-c",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("simple")

	cmd := buildCommand()
	err = cmd.Execute()
	t.Log(err.Error())

	is.Equal(err.Error(), "gomarkdoc: check mode cannot be run without an output set")
}

func TestCommand_defaultDirectory(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData/simple"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/simple/",
	}
	cleanup(".")

	main()

	verify(t, ".")
}

func TestCommand_nonexistant(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./nonexistant",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}

	cmd := buildCommand()
	err = cmd.Execute()
	t.Log(err.Error())
	is.Equal(err.Error(), fmt.Sprintf("gomarkdoc: invalid package in directory: .%snonexistant", string(filepath.Separator)))
}

func TestCommand_tags(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./tags",
		"--tags", "tagged",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verify(t, "./tags")
}

func TestCommand_tagsWithGOFLAGS(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Setenv("GOFLAGS", "-tags=tagged")
	os.Args = []string{
		"gomarkdoc", "./tags",
		"--config", "../.gomarkdoc-empty.yml",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verify(t, "./tags")
}

func TestCommand_tagsWithGOFLAGSNoTags(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	err = os.Setenv("GOFLAGS", "-other=foo")
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./tags",
		"--config", "../.gomarkdoc-empty.yml",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verifyNotEqual(t, "./tags")
}

func TestCommand_tagsWithGOFLAGSNoParse(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	err = os.Setenv("GOFLAGS", "invalid")
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./tags",
		"--config", "../.gomarkdoc-empty.yml",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verifyNotEqual(t, "./tags")
}

func TestCommand_embed(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./embed",
		"--embed",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("embed")

	data, err := os.ReadFile("./embed/README-template.md")
	is.NoErr(err)

	err = os.WriteFile("./embed/README-test.md", data, 0664)
	is.NoErr(err)

	main()

	verify(t, "./embed")
}

func TestCommand_untagged(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./untagged",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup("untagged")

	main()

	verify(t, "./untagged")
}

func TestCompare(t *testing.T) {
	tests := []struct {
		b1, b2 []byte
		equal  bool
	}{
		{[]byte("abc"), []byte("abc"), true},
		{[]byte("abc"), []byte("def"), false},
		{[]byte{}, []byte{}, true},
		{[]byte("abc"), []byte{}, false},
		{[]byte{}, []byte("abc"), false},
	}

	for _, test := range tests {
		name := fmt.Sprintf(`"%s" == "%s"`, string(test.b1), string(test.b2))
		if !test.equal {
			name = fmt.Sprintf(`"%s" != "%s"`, string(test.b1), string(test.b2))
		}

		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			eq, err := compare(bytes.NewBuffer(test.b1), bytes.NewBuffer(test.b2))
			is.NoErr(err)

			is.Equal(eq, test.equal)
		})
	}
}

func verify(t *testing.T, dir string) {
	is := is.New(t)

	data, err := os.ReadFile(filepath.Join(dir, "README.md"))
	is.NoErr(err)

	data2, err := os.ReadFile(filepath.Join(dir, "README-test.md"))
	is.NoErr(err)

	is.Equal(string(data), string(data2))
}

func verifyNotEqual(t *testing.T, dir string) {
	is := is.New(t)

	data, err := os.ReadFile(filepath.Join(dir, "README.md"))
	is.NoErr(err)

	data2, err := os.ReadFile(filepath.Join(dir, "README-test.md"))
	is.NoErr(err)

	is.True(string(data) != string(data2))
}

func cleanup(dir string) {
	os.Remove(filepath.Join(dir, "README-test.md"))
}
