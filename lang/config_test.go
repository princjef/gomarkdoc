package lang

import (
	"testing"

	"github.com/matryer/is"
)

func TestNormalizeRemote(t *testing.T) {
	tests := map[string]struct {
		raw        string
		normalized string
	}{
		"GitHub https": {
			raw:        "https://github.com/org/repo.git",
			normalized: "https://github.com/org/repo",
		},
		"GitHub ssh": {
			raw:        "git@github.com:org/repo.git",
			normalized: "https://github.com/org/repo",
		},
		"Azure DevOps https": {
			raw:        "https://org@dev.azure.com/org/project/_git/repo",
			normalized: "https://dev.azure.com/org/project/_git/repo",
		},
		"Azure DevOps ssh": {
			raw:        "git@ssh.dev.azure.com:v3/org/project/repo",
			normalized: "https://dev.azure.com/org/project/_git/repo",
		},
		"Azure DevOps https (visualstudio.com)": {
			raw:        "https://org.visualstudio.com/DefaultCollection/project/_git/repo",
			normalized: "https://dev.azure.com/org/project/_git/repo",
		},
		"Azure DevOps ssh (visualstudio.com)": {
			raw:        "org@vs-ssh.visualstudio.com:v3/org/project/repo",
			normalized: "https://dev.azure.com/org/project/_git/repo",
		},
		"GitLab https": {
			raw:        "https://gitlab.com/org/repo.git",
			normalized: "https://gitlab.com/org/repo",
		},
		"GitLab ssh": {
			raw:        "git@gitlab.com:org/repo.git",
			normalized: "https://gitlab.com/org/repo",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			normalized, ok := normalizeRemote(test.raw)
			is.True(ok)                           // Didn't produce a normlized value
			is.Equal(normalized, test.normalized) // Wrong normalized value
		})
	}
}
