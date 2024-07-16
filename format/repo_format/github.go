package repo_format

import (
	"fmt"
	"path/filepath"

	"github.com/princjef/gomarkdoc/lang"
)

type GitHubRepoFormat struct{}

func (g GitHubRepoFormat) CodeHref(loc lang.Location) (string, error) {
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
		locStr = fmt.Sprintf("L%d", loc.Start.Line)
	} else {
		locStr = fmt.Sprintf("L%d-L%d", loc.Start.Line, loc.End.Line)
	}

	return fmt.Sprintf(
		"%s/blob/%s/%s#%s",
		loc.Repo.Remote,
		loc.Repo.DefaultBranch,
		filepath.ToSlash(p),
		locStr,
	), nil
}

func (g GitHubRepoFormat) Supports(repoType lang.RepoType) bool {
	return repoType == lang.GitHubRepoType
}
