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
		if l := blankLineRegex.Find(rawText); l != nil {
			// Ignore blank lines
			rawText = rawText[len(l):]
		} else if l := headerRegex.Find(rawText); l != nil {
			headerText := strings.SplitN(string(l), "\n", 1)[0]
			blocks = append(blocks, NewBlock(HeaderBlock, headerText, level))
			rawText = rawText[len(l):]
		} else if l := spaceCodeBlockRegex.Find(rawText); l != nil {
			lines := strings.Split(string(l), "\n")

			minIndent := -1
			for _, line := range lines {
				for i, r := range line {
					if r != ' ' && (minIndent == -1 || i < minIndent) {
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

			blocks = append(blocks, NewBlock(CodeBlock, trimmedBlock.String(), level))
			rawText = rawText[len(l):]
		} else if l := tabCodeBlockRegex.Find(rawText); l != nil {
			lines := strings.Split(string(l), "\n")

			minIndent := -1
			for _, line := range lines {
				for i, r := range line {
					if r != '\t' && (minIndent == -1 || i < minIndent) {
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

			blocks = append(blocks, NewBlock(CodeBlock, trimmedBlock.String(), level))
			rawText = rawText[len(l):]
		} else if loc := multilineRegex.FindIndex(rawText); loc != nil {
			// Paragraph followed by something else
			paragraph := strings.TrimSpace(string(rawText[:loc[1]]))
			blocks = append(blocks, NewBlock(ParagraphBlock, formatDocParagraph(paragraph), level))
			rawText = rawText[loc[1]:]
		} else {
			// Last paragraph
			paragraph := strings.TrimSpace(string(rawText))

			var mergedParagraph strings.Builder
			for i, line := range strings.Split(paragraph, "\n") {
				if i > 0 {
					mergedParagraph.WriteRune(' ')
				}

				mergedParagraph.WriteString(strings.TrimSpace(line))
			}

			blocks = append(blocks, NewBlock(ParagraphBlock, mergedParagraph.String(), level))
			rawText = []byte{}
		}
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
