package repo_format

import (
	"fmt"
	"path/filepath"

	"github.com/princjef/gomarkdoc/lang"
)

type BitBucketRepoFormat struct{}

func (b BitBucketRepoFormat) CodeHref(loc lang.Location) (string, error) {
	var (
		relative string
		err      error
	)
	if filepath.IsAbs(loc.Filepath) {
		relative, err = filepath.Rel(loc.WorkDir, loc.Filepath)
		if err != nil {
			return "", err
		}
	} else {
		relative = loc.Filepath
	}

	full := filepath.Join(loc.Repo.PathFromRoot, relative)
	p, err := filepath.Rel(string(filepath.Separator), full)
	if err != nil {
		return "", err
	}

	var locStr string
	if loc.Start.Line == loc.End.Line {
		locStr = fmt.Sprintf("%d", loc.Start.Line)
	} else {
		locStr = fmt.Sprintf("%d-%d", loc.Start.Line, loc.End.Line)
	}

	return fmt.Sprintf(
		"%s/browse/%s?at=refs%%2Fheads%%2F%s#%s",
		loc.Repo.Remote,
		filepath.ToSlash(p),
		loc.Repo.DefaultBranch,
		locStr,
	), nil
}

func (b BitBucketRepoFormat) Supports(repoType lang.RepoType) bool {
	return repoType == lang.BitBucketRepoType
}
