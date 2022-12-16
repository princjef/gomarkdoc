package lang

import (
	"go/doc/comment"
	"strings"
)

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
	d := cfg.docParser.Parse(text)
	prt := cfg.docPrinter

	var blocks []*Block

loop:
	for i := 0; i < len(d.Content); {
		switch d.Content[i].(type) {
		case *comment.Heading:
			markdown := string(prt.Markdown(&comment.Doc{
				Links:   d.Links,
				Content: []comment.Block{d.Content[i]},
			}))
			hdr := strings.TrimSpace(strings.TrimLeft(markdown, "#")) // Remove leading '#' as they are added later by gomarkdoc
			blocks = append(blocks, NewBlock(cfg.Inc(0), HeaderBlock, hdr))
		case *comment.Code:
			markdown := string(prt.Markdown(&comment.Doc{
				Links:   d.Links,
				Content: []comment.Block{d.Content[i]},
			}))
			blocks = append(blocks, NewBlock(cfg.Inc(0), CodeBlock, strings.Trim(markdown, "\n")))
		default:
			paragraph := &comment.Doc{
				Links:   d.Links,
				Content: []comment.Block{},
			}
			for i < len(d.Content) {
				switch d.Content[i].(type) {
				case *comment.Heading, *comment.Code:
					blocks = append(blocks, NewBlock(cfg.Inc(0), ParagraphBlock, strings.Trim(string(prt.Markdown(paragraph)), "\n")))
					continue loop
				}
				paragraph.Content = append(paragraph.Content, d.Content[i])
				i++
			}
			blocks = append(blocks, NewBlock(cfg.Inc(0), ParagraphBlock, strings.Trim(string(prt.Markdown(paragraph)), "\n")))
		}
		i++
	}
	return &Doc{
		cfg:    cfg,
		blocks: blocks,
	}
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
