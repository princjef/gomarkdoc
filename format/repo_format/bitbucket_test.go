package repo_format

import (
	"testing"

	"github.com/matryer/is"

	"github.com/princjef/gomarkdoc/lang"
)

func TestBitBucketRepoFormat_CodeHref(t *testing.T) {
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
					Type:          lang.BitBucketRepoType,
					Remote:        "https://some.bitbucket.repo.com/user/repo",
					DefaultBranch: "master",
					PathFromRoot:  "/path/to/sources",
				},
			},
			href: "https://some.bitbucket.repo.com/user/repo/browse/path/to/sources/file.go?at=refs%2Fheads%2Fmaster#15",
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
					Type:          lang.BitBucketRepoType,
					Remote:        "https://some.bitbucket.repo.com/user/repo",
					DefaultBranch: "main",
					PathFromRoot:  "/",
				},
			},
			href: "https://some.bitbucket.repo.com/user/repo/browse/file.go?at=refs%2Fheads%2Fmain#15-20",
			err:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			format := BitBucketRepoFormat{}
			href, err := format.CodeHref(test.location)
			is.NoErr(err)
			is.Equal(test.href, href)
		})
	}
}
