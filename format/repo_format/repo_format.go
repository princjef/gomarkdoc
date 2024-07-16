package repo_format

import "github.com/princjef/gomarkdoc/lang"

// RepoFormat represent the format of concrete repository management system (RMS for short).
// Because different RMSs follow different rules (e.g., the URL pattern), we need to differentiate them
// from the Markdown rendering software to generate correct RMS related stuff
type RepoFormat interface {
	// CodeHref returns a formatted URL to the specified [lang.Location]
	//
	// It's guaranteed that [RepoFormat.Supports](loc.Repo.Type) is true
	CodeHref(loc lang.Location) (string, error)

	// Supports determines whether this [RepoFormat] supports specified [lang.RepoType]
	Supports(repoType lang.RepoType) bool
}
