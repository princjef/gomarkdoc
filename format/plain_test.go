package format_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/format"
	"github.com/princjef/gomarkdoc/lang"
)

func TestPlainMarkdown_Bold(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.Bold("sample text")
	is.NoErr(err)
	is.Equal(res, "**sample text**")
}

func TestPlainMarkdown_CodeBlock(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.CodeBlock("go", "Line 1\nLine 2")
	is.NoErr(err)
	is.Equal(res, "\tLine 1\n\tLine 2\n\n")
}

func TestPlainMarkdown_CodeBlock_noLanguage(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.CodeBlock("", "Line 1\nLine 2")
	is.NoErr(err)
	is.Equal(res, "\tLine 1\n\tLine 2\n\n")
}

func TestPlainMarkdown_Header(t *testing.T) {
	tests := []struct {
		text   string
		level  int
		result string
	}{
		{"header text", 1, "# header text\n\n"},
		{"level 2", 2, "## level 2\n\n"},
		{"level 3", 3, "### level 3\n\n"},
		{"level 4", 4, "#### level 4\n\n"},
		{"level 5", 5, "##### level 5\n\n"},
		{"level 6", 6, "###### level 6\n\n"},
		{"other level", 12, "###### other level\n\n"},
		{"with * escape", 2, "## with \\* escape\n\n"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s (level %d)", test.text, test.level), func(t *testing.T) {
			is := is.New(t)

			var f format.PlainMarkdown
			res, err := f.Header(test.level, test.text)
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestPlainMarkdown_Header_invalidLevel(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	_, err := f.Header(-1, "invalid")
	is.Equal(err.Error(), "format: header level cannot be less than 1")
}

func TestPlainMarkdown_RawHeader(t *testing.T) {
	tests := []struct {
		text   string
		level  int
		result string
	}{
		{"header text", 1, "# header text\n\n"},
		{"with * escape", 2, "## with * escape\n\n"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s (level %d)", test.text, test.level), func(t *testing.T) {
			is := is.New(t)

			var f format.PlainMarkdown
			res, err := f.RawHeader(test.level, test.text)
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestPlainMarkdown_LocalHref(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.LocalHref("Normal Header")
	is.NoErr(err)
	is.Equal(res, "")
}

func TestPlainMarkdown_CodeHref(t *testing.T) {
	is := is.New(t)

	wd, err := filepath.Abs(".")
	is.NoErr(err)
	locPath := filepath.Join(wd, "subdir", "file.go")

	var f format.PlainMarkdown
	res, err := f.CodeHref(lang.Location{
		Start:    lang.Position{Line: 12, Col: 1},
		End:      lang.Position{Line: 14, Col: 43},
		Filepath: locPath,
		WorkDir:  wd,
		Repo: &lang.Repo{
			Remote:        "https://dev.azure.com/org/project/_git/repo",
			DefaultBranch: "master",
			PathFromRoot:  "/",
		},
	})
	is.NoErr(err)
	is.Equal(res, "")
}

func TestPlainMarkdown_CodeHref_noRepo(t *testing.T) {
	is := is.New(t)

	wd, err := filepath.Abs(".")
	is.NoErr(err)
	locPath := filepath.Join(wd, "subdir", "file.go")

	var f format.PlainMarkdown
	res, err := f.CodeHref(lang.Location{
		Start:    lang.Position{Line: 12, Col: 1},
		End:      lang.Position{Line: 14, Col: 43},
		Filepath: locPath,
		WorkDir:  wd,
		Repo:     nil,
	})
	is.NoErr(err)
	is.Equal(res, "")
}

func TestPlainMarkdown_Link(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.Link("link text", "https://test.com/a/b/c")
	is.NoErr(err)
	is.Equal(res, "[link text](<https://test.com/a/b/c>)")
}

func TestPlainMarkdown_ListEntry(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.ListEntry(0, "list entry text")
	is.NoErr(err)
	is.Equal(res, "- list entry text\n")
}

func TestPlainMarkdown_ListEntry_nested(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.ListEntry(2, "nested text")
	is.NoErr(err)
	is.Equal(res, "    - nested text\n")
}

func TestPlainMarkdown_ListEntry_empty(t *testing.T) {
	is := is.New(t)

	var f format.PlainMarkdown
	res, err := f.ListEntry(0, "")
	is.NoErr(err)
	is.Equal(res, "")
}
