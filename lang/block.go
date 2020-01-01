package lang

type (
	// Block defines a single block element (e.g. paragraph, code block) in the
	// documentation for a symbol or package.
	Block struct {
		cfg  *Config
		kind BlockKind
		text string
	}

	// BlockKind identifies the type of block element represented by the
	// corresponding Block.
	BlockKind string
)

const (
	// ParagraphBlock defines a block that represents a paragraph of text.
	ParagraphBlock BlockKind = "paragraph"

	// CodeBlock defines a block that represents a section of code.
	CodeBlock BlockKind = "code"

	// HeaderBlock defines a block that represents a section header.
	HeaderBlock BlockKind = "header"
)

// NewBlock creates a new block element of the provided kind and with the given
// text contents.
func NewBlock(cfg *Config, kind BlockKind, text string) *Block {
	return &Block{cfg, kind, text}
}

// Level provides the default level that a block of kind HeaderBlock will render
// at in the output. The level is not used for other block types.
func (b *Block) Level() int {
	return b.cfg.Level
}

// Kind provides the kind of data that this block's text should be interpreted
// as.
func (b *Block) Kind() BlockKind {
	return b.kind
}

// Text provides the raw text of the block's contents. The text is pre-scrubbed
// and sanitized as determined by the block's Kind(), but it is not wrapped in
// any special constructs for rendering purposes (such as markdown code blocks).
func (b *Block) Text() string {
	return b.text
}
