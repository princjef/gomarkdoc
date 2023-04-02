package lang

// Doc provides access to the documentation comment contents for a package or
// symbol in a structured form.
type Doc struct {
	cfg    *Config
	blocks []*Block
}

// NewDoc initializes a Doc struct from the provided raw documentation text and
// with headers rendered by default at the heading level provided. Documentation
// is separated into block level elements using the standard rules from golang's
// documentation conventions.
func NewDoc(cfg *Config, text string) *Doc {
	// Replace CRLF with LF
	rawText := normalizeDoc(text)

	parsed := cfg.Pkg.Parser().Parse(rawText)

	blocks := ParseBlocks(cfg, parsed.Content, false)

	return &Doc{cfg, blocks}
}

// Level provides the default level that headers within the documentation should
// be rendered
func (d *Doc) Level() int {
	return d.cfg.Level
}

// Blocks holds the list of block elements that makes up the documentation
// contents.
func (d *Doc) Blocks() []*Block {
	return d.blocks
}
