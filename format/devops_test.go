package format_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/princjef/gomarkdoc/format"
	"github.com/princjef/gomarkdoc/lang"
)

func TestBold(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.Bold("sample text")
	is.NoErr(err)
	is.Equal(res, "**sample text**")
}

func TestCodeBlock(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.CodeBlock("go", "Line 1\nLine 2")
	is.NoErr(err)
	is.Equal(res, "```go\nLine 1\nLine 2\n```\n\n")
}

func TestCodeBlock_noLanguage(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.CodeBlock("", "Line 1\nLine 2")
	is.NoErr(err)
	is.Equal(res, "```\nLine 1\nLine 2\n```\n\n")
}

func TestHeader(t *testing.T) {
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

			var f format.AzureDevOpsMarkdown
			res, err := f.Header(test.level, test.text)
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestHeader_invalidLevel(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	_, err := f.Header(-1, "invalid")
	is.Equal(err.Error(), "format: header level cannot be less than 1")
}

func TestRawHeader(t *testing.T) {
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

			var f format.AzureDevOpsMarkdown
			res, err := f.RawHeader(test.level, test.text)
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestLocalHref(t *testing.T) {
	tests := map[string]string{
		"Normal Header":       "#normal-header",
		" Leading whitespace": "#leading-whitespace",
		"Multiple	 whitespace": "#multiple--whitespace",
		"Special(#)%^Characters": "#special%28%23%29%25%5Echaracters",
		"With:colon":             "#with%3Acolon",
	}

	for input, output := range tests {
		t.Run(input, func(t *testing.T) {
			is := is.New(t)

			var f format.AzureDevOpsMarkdown
			res, err := f.LocalHref(input)
			is.NoErr(err)
			is.Equal(res, output)
		})
	}
}

func TestCodeHref(t *testing.T) {
	is := is.New(t)

	repoDir, err := filepath.Abs(".")
	is.NoErr(err)

	locPath, err := filepath.Abs(".")
	is.NoErr(err)
	locPath = filepath.Join(locPath, "subdir", "file.go")

	var f format.AzureDevOpsMarkdown
	res, err := f.CodeHref(lang.Location{
		Start:    lang.Position{Line: 12, Col: 1},
		End:      lang.Position{Line: 14, Col: 43},
		Filepath: locPath,
		Repo: &lang.Repo{
			Remote:        "https://dev.azure.com/org/project/_git/repo",
			DefaultBranch: "master",
			RootDir:       repoDir,
		},
	})
	is.NoErr(err)
	is.Equal(res, "https://dev.azure.com/org/project/_git/repo?path=subdir%2Ffile.go&version=GBmaster&lineStyle=plain&line=12&lineEnd=14&lineStartColumn=1&lineEndColumn=43")
}

func TestCodeHref_noRepo(t *testing.T) {
	is := is.New(t)

	locPath, err := filepath.Abs(".")
	is.NoErr(err)
	locPath = filepath.Join(locPath, "subdir", "file.go")

	var f format.AzureDevOpsMarkdown
	res, err := f.CodeHref(lang.Location{
		Start:    lang.Position{Line: 12, Col: 1},
		End:      lang.Position{Line: 14, Col: 43},
		Filepath: locPath,
		Repo:     nil,
	})
	is.NoErr(err)
	is.Equal(res, "")
}

func TestLink(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.Link("link text", "https://test.com/a/b/c")
	is.NoErr(err)
	is.Equal(res, "[link text](<https://test.com/a/b/c>)")
}

func TestListEntry(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.ListEntry(0, "list entry text")
	is.NoErr(err)
	is.Equal(res, "- list entry text\n")
}

func TestListEntry_nested(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.ListEntry(2, "nested text")
	is.NoErr(err)
	is.Equal(res, "    - nested text\n")
}

func TestListEntry_empty(t *testing.T) {
	is := is.New(t)

	var f format.AzureDevOpsMarkdown
	res, err := f.ListEntry(0, "")
	is.NoErr(err)
	is.Equal(res, "")
}
