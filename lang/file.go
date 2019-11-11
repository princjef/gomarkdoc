package lang

// File holds information for rendering a single file that contains one or more
// packages.
type File struct {
	Header   string
	Footer   string
	Packages []*Package
}

// NewFile creates a new instance of File with the provided information.
func NewFile(header, footer string, packages []*Package) *File {
	return &File{
		Header:   header,
		Footer:   footer,
		Packages: packages,
	}
}
