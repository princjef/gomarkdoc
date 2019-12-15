package lang

import (
	"regexp"
	"strings"
)

// Doc provides access to the documentation comment contents for a package or
// symbol in a structured form.
type Doc struct {
	level  int
	blocks []*Block
}

var (
	multilineRegex      = regexp.MustCompile("\n(?:[\t\f ]*\n)+")
	headerRegex         = regexp.MustCompile(`^[A-Z][^!:;,{}\[\]<>.?]*\n(?:[\t\f ]*\n)`)
	spaceCodeBlockRegex = regexp.MustCompile(`^(?:(?:(?:(?:  ).*[^\s]+.*)|[\t\f ]*)\n)+`)
	tabCodeBlockRegex   = regexp.MustCompile(`^(?:(?:(?:\t.*[^\s]+.*)|[\t\f ]*)\n)+`)
	blankLineRegex      = regexp.MustCompile(`^[\t\f ]*\n`)
)

// NewDoc initializes a Doc struct from the provided raw documentation text and
// with headers rendered by default at the heading level provided. Documentation
// is separated into block level elements using the standard rules from golang's
// documentation conventions.
func NewDoc(text string, level int) *Doc {
	// Replace CRLF with LF
	rawText := []byte(normalizeDoc(text) + "\n")

	var blocks []*Block
	for len(rawText) > 0 {
		// Blank lines (ignore)
		if l, ok := parseBlankLine(rawText); ok {
			rawText = rawText[l:]
			continue
		}

		// Header
		if b, l, ok := parseHeaderBlock(rawText, level); ok {
			blocks = append(blocks, b)
			rawText = rawText[l:]
			continue
		}

		// Code block
		if b, l, ok := parseCodeBlock(rawText, level); ok {
			blocks = append(blocks, b)
			rawText = rawText[l:]
			continue
		}

		// Paragraph
		b, l := parseParagraph(rawText, level)
		blocks = append(blocks, b)
		rawText = rawText[l:]
	}

	return &Doc{level, blocks}
}

// Level provides the default level that headers within the documentation should
// be rendered
func (d *Doc) Level() int {
	return d.level
}

// Blocks holds the list of block elements that makes up the documentation
// contents.
func (d *Doc) Blocks() []*Block {
	return d.blocks
}

func parseBlankLine(text []byte) (length int, ok bool) {
	if l := blankLineRegex.Find(text); l != nil {
		// Ignore blank lines
		return len(l), true
	}

	return 0, false
}

func parseHeaderBlock(text []byte, level int) (block *Block, length int, ok bool) {
	if l := headerRegex.Find(text); l != nil {
		headerText := strings.TrimSpace(string(l))
		return NewBlock(HeaderBlock, headerText, level), len(l), true
	}

	return nil, 0, false
}

func parseCodeBlock(text []byte, level int) (block *Block, length int, ok bool) {
	l := spaceCodeBlockRegex.Find(text)
	var indent rune
	if l != nil {
		indent = ' '
	} else {
		l = tabCodeBlockRegex.Find(text)
		if l != nil {
			indent = '\t'
		} else {
			return nil, 0, false
		}
	}

	lines := strings.Split(string(l), "\n")

	minIndent := -1
	for _, line := range lines {
		for i, r := range line {
			if r != indent && (minIndent == -1 || i < minIndent) {
				minIndent = i
			}
		}
	}

	var trimmedBlock strings.Builder
	for i, line := range lines {
		if i > 0 {
			trimmedBlock.WriteRune('\n')
		}

		if len(strings.TrimSpace(line)) > 0 {
			trimmedBlock.WriteString(line[minIndent:])
		}
	}

	return NewBlock(CodeBlock, trimmedBlock.String(), level), len(l), true
}

func parseParagraph(text []byte, level int) (block *Block, length int) {
	if loc := multilineRegex.FindIndex(text); loc != nil {
		// Paragraph followed by something else
		paragraph := strings.TrimSpace(string(text[:loc[1]]))
		return NewBlock(ParagraphBlock, formatDocParagraph(paragraph), level), loc[1]
	}

	// Last paragraph
	paragraph := strings.TrimSpace(string(text))

	var mergedParagraph strings.Builder
	for i, line := range strings.Split(paragraph, "\n") {
		if i > 0 {
			mergedParagraph.WriteRune(' ')
		}

		mergedParagraph.WriteString(strings.TrimSpace(line))
	}

	return NewBlock(ParagraphBlock, mergedParagraph.String(), level), len(text)
}
