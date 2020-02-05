package lang

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/princjef/gomarkdoc/logger"
	"gopkg.in/src-d/go-git.v4"
)

type (
	// Config defines contextual information used to resolve documentation for
	// a construct.
	Config struct {
		FileSet *token.FileSet
		Level   int
		Repo    *Repo
		PkgDir  string
		Log     logger.Logger
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
func NewConfig(log logger.Logger, pkgDir string) (*Config, error) {
	dir, err := filepath.Abs(pkgDir)
	if err != nil {
		return nil, err
	}

	repo, err := getRepoForDir(log, dir)
	if err != nil {
		log.Infof("unable to resolve repository due to error: %s", err)
		return &Config{
			FileSet: token.NewFileSet(),
			Level:   1,
			PkgDir:  dir,
			Log:     log,
		}, nil
	}

	log.Debugf(
		"resolved repository with remote %s, default branch %s, root directory %s",
		repo.Remote,
		repo.DefaultBranch,
		repo.RootDir,
	)
	return &Config{
		FileSet: token.NewFileSet(),
		Level:   1,
		Repo:    repo,
		PkgDir:  dir,
		Log:     log,
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

func getRepoForDir(log logger.Logger, dir string) (*Repo, error) {
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
		if repo, ok := processRemote(log, repo, r); ok {
			ri = repo
			break
		}
	}

	// If there's no "origin", just use the first remote
	if ri == nil {
		if len(remotes) == 0 {
			return nil, errors.New("no remotes found for repository")
		}

		repo, ok := processRemote(log, repo, remotes[0])
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

func processRemote(log logger.Logger, repository *git.Repository, remote *git.Remote) (*Repo, bool) {
	repo := &Repo{}

	c := remote.Config()

	// TODO: configurable remote name?
	if c.Name != "origin" || len(c.URLs) == 0 {
		log.Debugf("skipping remote because it is not the origin or it has no URLs")
		return nil, false
	}

	refs, err := repository.References()
	if err != nil {
		log.Debugf("skipping remote %s because listing its refs failed: %s", c.URLs[0], err)
		return nil, false
	}

	prefix := fmt.Sprintf("refs/remotes/%s/", c.Name)
	headRef := fmt.Sprintf("refs/remotes/%s/HEAD", c.Name)

	for {
		ref, err := refs.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Debugf("skipping remote %s because listing its refs failed: %s", c.URLs[0], err)
			return nil, false
		}
		defer refs.Close()

		if ref == nil {
			break
		}

		if string(ref.Name()) == headRef && strings.HasPrefix(string(ref.Target()), prefix) {
			repo.DefaultBranch = strings.TrimPrefix(string(ref.Target()), prefix)
			log.Debugf("found default branch %s for remote %s", repo.DefaultBranch, c.URLs[0])
			break
		}
	}

	if repo.DefaultBranch == "" {
		log.Debugf("skipping remote %s because no default branch was found", c.URLs[0])
		return nil, false
	}

	normalized, ok := normalizeRemote(c.URLs[0])
	if !ok {
		log.Debugf("skipping remote %s because its remote URL could not be normalized", c.URLs[0])
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
