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
	tests := []string{
		"./simple",
		"./lang/function",
		"./docs",
		"./untagged",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			is := is.New(t)

			err := os.Chdir(filepath.Join(wd, "../../testData"))
			is.NoErr(err)

			harness(t, test, []string{
				"gomarkdoc", test,
				"--repository.url", "https://github.com/princjef/gomarkdoc",
				"--repository.default-branch", "master",
				"--repository.path", "/testData/",
			})
		})
	}
}

func TestCommand_check(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./simple",
		"-c",
		"-o", "{{.Dir}}/README-github.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup(t, "simple")

	main()
}

func TestCommand_nested(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./nested/...",
		"-o", "{{.Dir}}/README-github-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup(t, "nested")
	cleanup(t, "nested/inner")

	main()

	verify(t, "nested", "github")
	verify(t, "nested/inner", "github")
}

func TestCommand_unexported(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	harness(t, "unexported", []string{
		"gomarkdoc", "./unexported",
		"-u",
		"-o", "{{.Dir}}/README-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	})
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
	cleanup(t, "simple")

	cmd := buildCommand()
	err = cmd.Execute()
	t.Log(err.Error())

	is.Equal(err.Error(), "gomarkdoc: check mode cannot be run without an output set")
}

func TestCommand_defaultDirectory(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData/simple"))
	is.NoErr(err)

	harness(t, ".", []string{
		"gomarkdoc",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/simple/",
	})
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

	harness(t, "tags", []string{
		"gomarkdoc", "./tags",
		"--tags", "tagged",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	})
}

func TestCommand_tagsWithGOFLAGS(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Setenv("GOFLAGS", "-tags=tagged")
	os.Args = []string{
		"gomarkdoc", "./tags",
		"--config", "../.gomarkdoc-empty.yml",
		"-o", "{{.Dir}}/README-github-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup(t, "tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verify(t, "./tags", "github")
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
		"-o", "{{.Dir}}/README-github-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup(t, "tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verifyNotEqual(t, "./tags", "github")
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
		"-o", "{{.Dir}}/README-github-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup(t, "tags")

	cmd := buildCommand()
	err = cmd.Execute()
	is.NoErr(err)

	verifyNotEqual(t, "./tags", "github")
}

func TestCommand_embed(t *testing.T) {
	is := is.New(t)

	err := os.Chdir(filepath.Join(wd, "../../testData"))
	is.NoErr(err)

	os.Args = []string{
		"gomarkdoc", "./embed",
		"--embed",
		"-o", "{{.Dir}}/README-github-test.md",
		"--repository.url", "https://github.com/princjef/gomarkdoc",
		"--repository.default-branch", "master",
		"--repository.path", "/testData/",
	}
	cleanup(t, "embed")

	data, err := os.ReadFile("./embed/README-template.md")
	is.NoErr(err)

	err = os.WriteFile("./embed/README-github-test.md", data, 0664)
	is.NoErr(err)

	main()

	verify(t, "./embed", "github")
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

func verify(t *testing.T, dir, format string) {
	is := is.New(t)

	data, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("README-%s.md", format)))
	is.NoErr(err)

	data2, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("README-%s-test.md", format)))
	is.NoErr(err)

	is.Equal(string(data), string(data2))
}

func verifyNotEqual(t *testing.T, dir, format string) {
	is := is.New(t)

	data, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("README-%s.md", format)))
	is.NoErr(err)

	data2, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf("README-%s-test.md", format)))
	is.NoErr(err)

	is.True(string(data) != string(data2))
}

func cleanup(t *testing.T, dir string) {
	f, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, n := range f {
		if n.IsDir() {
			continue
		}

		if !strings.HasPrefix(n.Name(), "README") || !strings.HasSuffix(n.Name(), "-test.md") {
			continue
		}

		os.Remove(filepath.Join(dir, n.Name()))

	}
}

// harness runs the test for all formats. Omit the --output and --format args to
// the command when running this as it will fill them in for you
func harness(t *testing.T, dir string, args []string) {
	for _, format := range []string{"plain", "github", "azure-devops"} {
		os.Args = args
		os.Args = append(os.Args, "-o", fmt.Sprintf("{{.Dir}}/README-%s-test.md", format))
		os.Args = append(os.Args, "--format", format)

		cleanup(t, dir)

		main()

		verify(t, dir, format)
	}
}
