package repo_format

import (
	"testing"

	"github.com/matryer/is"

	"github.com/princjef/gomarkdoc/lang"
)

func TestGitHubRepoFormat_CodeHref(t *testing.T) {
	tests := []struct {
		name     string
		location lang.Location
		href     string
		err      error
	}{
		{
			name: "concrete_location",
			location: lang.Location{
				Start: lang.Position{
					Col:  8,
					Line: 15,
				},
				End: lang.Position{
					Col:  60,
					Line: 15,
				},
				Filepath: "/some/path/to/workdir/file.go",
				WorkDir:  "/some/path/to/workdir",
				Repo: &lang.Repo{
					Type:          lang.GitHubRepoType,
					Remote:        "https://some.git.repo.com/user/repo",
					DefaultBranch: "master",
					PathFromRoot:  "/path/to/sources",
				},
			},
			href: "https://some.git.repo.com/user/repo/blob/master/path/to/sources/file.go#L15",
			err:  nil,
		},
		{
			name: "range",
			location: lang.Location{
				Start: lang.Position{
					Col:  8,
					Line: 15,
				},
				End: lang.Position{
					Col:  60,
					Line: 20,
				},
				Filepath: "/some/path/to/workdir/file.go",
				WorkDir:  "/some/path/to/workdir",
				Repo: &lang.Repo{
					Type:          lang.GitHubRepoType,
					Remote:        "https://some.git.repo.com/user/repo",
					DefaultBranch: "main",
					PathFromRoot:  "/",
				},
			},
			href: "https://some.git.repo.com/user/repo/blob/main/file.go#L15-L20",
			err:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			format := GitHubRepoFormat{}
			href, err := format.CodeHref(test.location)
			is.NoErr(err)
			is.Equal(test.href, href)
		})
	}
}
