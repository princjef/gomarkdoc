package lang

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type (
	// Config defines contextual information used to resolve documentation for
	// a construct.
	Config struct {
		FileSet *token.FileSet
		Level   int
		Repo    *Repo
		PkgDir  string
	}

	// Repo represents information about a repository relevant to documentation
	// generation.
	Repo struct {
		Remote        string
		DefaultBranch string
		RootDir       string
	}

	// Location holds information for identifying a position within a file and
	// repository, if present.
	Location struct {
		Start    Position
		End      Position
		Filepath string
		Repo     *Repo
	}

	// Position represents a line and column number within a file.
	Position struct {
		Line int
		Col  int
	}
)

// NewConfig generates a Config for the provided package directory. It will
// resolve the filepath and attempt to determine the repository containing the
// directory. If no repository is found, the Repo field will be set to nil. An
// error is returned if the provided directory is invalid.
func NewConfig(pkgDir string) (*Config, error) {
	dir, err := filepath.Abs(pkgDir)
	if err != nil {
		return nil, err
	}

	repo, err := getRepoForDir(dir)
	if err != nil {
		return &Config{
			FileSet: token.NewFileSet(),
			Level:   1,
			PkgDir:  dir,
		}, nil
	}

	return &Config{
		FileSet: token.NewFileSet(),
		Level:   1,
		Repo:    repo,
		PkgDir:  dir,
	}, nil
}

// Inc copies the Config and increments the level by the provided step.
func (c *Config) Inc(step int) *Config {
	return &Config{
		FileSet: c.FileSet,
		Level:   c.Level + step,
		Repo:    c.Repo,
	}
}

func getRepoForDir(dir string) (*Repo, error) {
	var ri *Repo

	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}

	for _, r := range remotes {
		// TODO: configurable remote name?
		if r.Config().Name != "origin" {
			continue
		}

		if repo, ok := processRemote(r); ok {
			ri = repo
			break
		}
	}

	// If there's no "origin", just use the first remote
	if ri == nil {
		if len(remotes) == 0 {
			return nil, errors.New("no remotes found for repository")
		}

		repo, ok := processRemote(remotes[0])
		if !ok {
			return nil, errors.New("no remotes found for repository")
		}

		ri = repo
	}

	t, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	ri.RootDir = t.Filesystem.Root()

	return ri, nil
}

func processRemote(remote *git.Remote) (*Repo, bool) {
	repo := &Repo{}

	c := remote.Config()

	// TODO: configurable remote name?
	if c.Name != "origin" || len(c.URLs) == 0 {
		return nil, false
	}

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, false
	}

	for _, ref := range refs {
		if ref.Name() == plumbing.HEAD && strings.HasPrefix(string(ref.Target()), "refs/heads/") {
			repo.DefaultBranch = strings.TrimPrefix(string(ref.Target()), "refs/heads/")
			break
		}
	}

	if repo.DefaultBranch == "" {
		return nil, false
	}

	normalized, ok := normalizeRemote(c.URLs[0])
	if !ok {
		return nil, false
	}

	repo.Remote = normalized
	return repo, true
}

var (
	sshRemoteRegex       = regexp.MustCompile(`^[\w-]+@([^:]+):(.+?)(?:\.git)?$`)
	httpsRemoteRegex     = regexp.MustCompile(`^(https?://)(?:[^@/]+@)?([\w-.]+)(/.+?)?(?:\.git)?$`)
	devOpsSSHV3PathRegex = regexp.MustCompile(`^v3/([^/]+)/([^/]+)/([^/]+)$`)
	devOpsHTTPSPathRegex = regexp.MustCompile(`^/([^/]+)/([^/]+)/_git/([^/]+)$`)
)

func normalizeRemote(remote string) (string, bool) {
	if match := sshRemoteRegex.FindStringSubmatch(remote); match != nil {
		switch match[1] {
		case "ssh.dev.azure.com", "vs-ssh.visualstudio.com":
			if pathMatch := devOpsSSHV3PathRegex.FindStringSubmatch(match[2]); pathMatch != nil {
				// DevOps v3
				return fmt.Sprintf(
					"https://dev.azure.com/%s/%s/_git/%s",
					pathMatch[1],
					pathMatch[2],
					pathMatch[3],
				), true
			}

			return "", false
		default:
			// GitHub and friends
			return fmt.Sprintf("https://%s/%s", match[1], match[2]), true
		}
	}

	if match := httpsRemoteRegex.FindStringSubmatch(remote); match != nil {
		switch {
		case match[2] == "dev.azure.com":
			if pathMatch := devOpsHTTPSPathRegex.FindStringSubmatch(match[3]); pathMatch != nil {
				// DevOps
				return fmt.Sprintf(
					"https://dev.azure.com/%s/%s/_git/%s",
					pathMatch[1],
					pathMatch[2],
					pathMatch[3],
				), true
			}

			return "", false
		case strings.HasSuffix(match[2], ".visualstudio.com"):
			if pathMatch := devOpsHTTPSPathRegex.FindStringSubmatch(match[3]); pathMatch != nil {
				// DevOps (old domain)

				// Pull off the beginning of the domain
				org := strings.SplitN(match[2], ".", 2)[0]
				return fmt.Sprintf(
					"https://dev.azure.com/%s/%s/_git/%s",
					org,
					pathMatch[2],
					pathMatch[3],
				), true
			}

			return "", false
		default:
			// GitHub and friends
			return fmt.Sprintf("%s%s%s", match[1], match[2], match[3]), true
		}
	}

	// TODO: error instead?
	return "", false
}

// NewLocation returns a location for the provided Config and ast.Node
// combination. This is typically not called directly, but is made available via
// the Location() methods of various lang constructs.
func NewLocation(cfg *Config, node ast.Node) Location {
	start := cfg.FileSet.Position(node.Pos())
	end := cfg.FileSet.Position(node.End())

	return Location{
		Start:    Position{start.Line, start.Column},
		End:      Position{end.Line, end.Column},
		Filepath: start.Filename,
		Repo:     cfg.Repo,
	}
}
