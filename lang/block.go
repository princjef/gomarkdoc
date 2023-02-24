package lang

import (
	"go/doc/comment"
	"regexp"
	"strings"
)

type (
	// Block defines a single block element (e.g. paragraph, code block) in the
	// documentation for a symbol or package.
	Block struct {
		cfg    *Config
		kind   BlockKind
		text   string
		list   *List
		inline bool
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

	// ListBlock defines a block that represents an ordered or unordered list.
	ListBlock BlockKind = "list"
)

// NewBlock creates a new block element of the provided kind and with the given
// text contents and a flag indicating whether this block is part of an inline
// element.
func NewBlock(cfg *Config, kind BlockKind, text string, inline bool) *Block {
	return &Block{cfg, kind, text, nil, inline}
}

// NewListBlock creates a new list block element and with the given list
// definition and a flag indicating whether this block is part of an inline
// element.
func NewListBlock(cfg *Config, list *List, inline bool) *Block {
	return &Block{cfg, ListBlock, "", list, inline}
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

// List provides the list contents for a list block. Only relevant for blocks of
// type ListBlock.
func (b *Block) List() *List {
	return b.list
}

// Inline indicates whether the block is part of an inline element, such as a
// list item.
func (b *Block) Inline() bool {
	return b.inline
}

// ParseBlocks produces a set of blocks from the corresponding comment blocks.
// It also takes a flag indicating whether the blocks are part of an inline
// element such as a list item.
func ParseBlocks(cfg *Config, blocks []comment.Block, inline bool) []*Block {
	res := make([]*Block, len(blocks))
	for i, b := range blocks {
		switch v := b.(type) {
		case *comment.Code:
			res[i] = NewBlock(cfg.Inc(0), CodeBlock, v.Text, inline)
		case *comment.Heading:
			var b strings.Builder
			printText(&b, v.Text...)
			res[i] = NewBlock(cfg.Inc(0), HeaderBlock, b.String(), inline)
		case *comment.List:
			list := NewList(cfg.Inc(0), v)
			res[i] = NewListBlock(cfg.Inc(0), list, inline)
		case *comment.Paragraph:
			var b strings.Builder
			printText(&b, v.Text...)
			text := collapseWhitespace(b.String())
			res[i] = NewBlock(cfg.Inc(0), ParagraphBlock, text, inline)
		}
	}

	return res
}

func printText(b *strings.Builder, text ...comment.Text) {
	for i, t := range text {
		if i > 0 {
			b.WriteRune(' ')
		}

		switch v := t.(type) {
		case comment.Plain:
			b.WriteString(string(v))
		case comment.Italic:
			b.WriteString(string(v))
		case *comment.DocLink:
			printText(b, v.Text...)
		case *comment.Link:
			printText(b, v.Text...)
		}
	}
}

var whitespaceRegex = regexp.MustCompile(`\s+`)

func collapseWhitespace(s string) string {
	return string(whitespaceRegex.ReplaceAll([]byte(s), []byte(" ")))
}
