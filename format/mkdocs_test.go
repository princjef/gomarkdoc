package format_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	"github.com/princjef/gomarkdoc/format"
	"github.com/princjef/gomarkdoc/format/repo_format"
	"github.com/princjef/gomarkdoc/lang"
)

func TestMkDocsMarkdown_Bold(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.Bold("sample text")
	is.NoErr(err)
	is.Equal(res, "**sample text**")
}

func TestMkDocsMarkdown_CodeBlock(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.CodeBlock("go", "Line 1\nLine 2")
	is.NoErr(err)
	is.Equal(res, "```go\nLine 1\nLine 2\n```")
}

func TestMkDocsMarkdown_CodeBlock_noLanguage(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.CodeBlock("", "Line 1\nLine 2")
	is.NoErr(err)
	is.Equal(res, "```\nLine 1\nLine 2\n```")
}

func TestMkDocsMarkdown_Header(t *testing.T) {
	tests := []struct {
		text   string
		level  int
		result string
	}{
		{"header text", 1, "# header text"},
		{"level 2", 2, "## level 2"},
		{"level 3", 3, "### level 3"},
		{"level 4", 4, "#### level 4"},
		{"level 5", 5, "##### level 5"},
		{"level 6", 6, "###### level 6"},
		{"other level", 12, "###### other level"},
		{"with * escape", 2, "## with \\* escape"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s (level %d)", test.text, test.level), func(t *testing.T) {
			is := is.New(t)

			var f format.MkDocsMarkdown
			res, err := f.Header(test.level, test.text)
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestMkDocsMarkdown_Header_invalidLevel(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	_, err := f.Header(-1, "invalid")
	is.Equal(err.Error(), "format: header level cannot be less than 1")
}

func TestMkDocsMarkdown_RawHeader(t *testing.T) {
	tests := []struct {
		text   string
		level  int
		result string
	}{
		{"header text", 1, "# header text"},
		{"with * escape", 2, "## with * escape"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s (level %d)", test.text, test.level), func(t *testing.T) {
			is := is.New(t)

			var f format.MkDocsMarkdown
			res, err := f.RawHeader(test.level, test.text)
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestMkDocsMarkdown_LocalHref(t *testing.T) {
	tests := map[string]string{
		"Normal Header":          "#normal-header",
		" Leading whitespace":    "#leading-whitespace",
		"Multiple	 whitespace":   "#multiple--whitespace",
		"Special(#)%^Characters": "#specialcharacters",
		"With:colon":             "#withcolon",
	}

	for input, output := range tests {
		t.Run(input, func(t *testing.T) {
			is := is.New(t)

			var f format.MkDocsMarkdown
			res, err := f.LocalHref(input)
			is.NoErr(err)
			is.Equal(res, output)
		})
	}
}

type mockRepoFormat struct {
	t lang.RepoType
	r string
}

func (m *mockRepoFormat) CodeHref(_ lang.Location) (string, error) {
	return m.r, nil
}

func (m *mockRepoFormat) Supports(repoType lang.RepoType) bool {
	return m.t == repoType
}

func TestMkDocsMarkdown_CodeHref(t *testing.T) {
	tests := []struct {
		name     string
		repoType lang.RepoType
		result   string
	}{
		{
			name:     "github",
			repoType: lang.GitHubRepoType,
			result:   "from github",
		},
		{
			name:     "bitbucket",
			repoType: lang.BitBucketRepoType,
			result:   "from bitbucket",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)

			wd, err := filepath.Abs(".")
			is.NoErr(err)
			locPath := filepath.Join(wd, "subdir", "file.go")

			f := &format.MkDocsMarkdown{
				RepoFormats: []repo_format.RepoFormat{
					&mockRepoFormat{
						t: lang.GitHubRepoType,
						r: "from github",
					},
					&mockRepoFormat{
						t: lang.BitBucketRepoType,
						r: "from bitbucket",
					},
				},
			}
			res, err := f.CodeHref(lang.Location{
				Start:    lang.Position{Line: 12, Col: 1},
				End:      lang.Position{Line: 14, Col: 43},
				Filepath: locPath,
				WorkDir:  wd,
				Repo: &lang.Repo{
					Type:          test.repoType,
					Remote:        "https://dev.azure.com/org/project/_git/repo",
					DefaultBranch: "master",
					PathFromRoot:  "/",
				},
			})
			is.NoErr(err)
			is.Equal(res, test.result)
		})
	}
}

func TestMkDocsMarkdown_CodeHref_noRepo(t *testing.T) {
	is := is.New(t)

	wd, err := filepath.Abs(".")
	is.NoErr(err)
	locPath := filepath.Join(wd, "subdir", "file.go")

	var f format.MkDocsMarkdown
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

func TestMkDocsMarkdown_Link(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.Link("link text", "https://test.com/a/b/c")
	is.NoErr(err)
	is.Equal(res, "[link text](<https://test.com/a/b/c>)")
}

func TestMkDocsMarkdown_ListEntry(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.ListEntry(0, "list entry text")
	is.NoErr(err)
	is.Equal(res, "- list entry text")
}

func TestMkDocsMarkdown_ListEntry_nested(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.ListEntry(2, "nested text")
	is.NoErr(err)
	is.Equal(res, "    - nested text")
}

func TestMkDocsMarkdown_ListEntry_empty(t *testing.T) {
	is := is.New(t)

	var f format.MkDocsMarkdown
	res, err := f.ListEntry(0, "")
	is.NoErr(err)
	is.Equal(res, "")
}
