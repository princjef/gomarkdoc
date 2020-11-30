package format

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"
	"mvdan.cc/xurls/v2"
)

// bold converts the provided text to bold
func bold(text string) string {
	if text == "" {
		return ""
	}

	return fmt.Sprintf("**%s**", escape(text))
}

// codeBlock wraps the provided code as a code block. Language syntax
// highlighting is not supported.
func codeBlock(code string) string {
	var builder strings.Builder

	lines := strings.Split(code, "\n")
	for _, line := range lines {
		builder.WriteRune('\t')
		builder.WriteString(line)
	}

	return builder.String()
}

// gfmCodeBlock wraps the provided code as a code block and tags it with the
// provided language (or no language if the empty string is provided), using
// the triple backtick format from GitHub Flavored Markdown.
func gfmCodeBlock(language, code string) string {
	return fmt.Sprintf("```%s\n%s\n```\n\n", language, strings.TrimSpace(code))
}

// header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func header(level int, text string) (string, error) {
	if level < 1 {
		return "", errors.New("format: header level cannot be less than 1")
	}

	switch level {
	case 1:
		return fmt.Sprintf("# %s\n\n", text), nil
	case 2:
		return fmt.Sprintf("## %s\n\n", text), nil
	case 3:
		return fmt.Sprintf("### %s\n\n", text), nil
	case 4:
		return fmt.Sprintf("#### %s\n\n", text), nil
	case 5:
		return fmt.Sprintf("##### %s\n\n", text), nil
	default:
		// Only go up to 6 levels. Anything higher is also level 6
		return fmt.Sprintf("###### %s\n\n", text), nil
	}
}

// link generates a link with the given text and href values.
func link(text, href string) string {
	if text == "" {
		return ""
	}

	if href == "" {
		return text
	}

	return fmt.Sprintf("[%s](<%s>)", text, href)
}

// listEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func listEntry(depth int, text string) string {
	// TODO: this is a weird special case
	if text == "" {
		return ""
	}

	prefix := strings.Repeat("  ", depth)
	return fmt.Sprintf("%s- %s\n", prefix, text)
}

// gfmAccordion generates a collapsible content. The accordion's visible title
// while collapsed is the provided title and the expanded content is the body.
func gfmAccordion(title, body string) string {
	return fmt.Sprintf("<details><summary>%s</summary>\n<p>\n\n%s</p>\n</details>\n\n", title, escape(body))
}

// gfmAccordionHeader generates the header visible when an accordion is
// collapsed.
//
// The gfmAccordionHeader is expected to be used in conjunction with
// gfmAccordionTerminator() when the demands of the body's rendering requires
// it to be generated independently. The result looks conceptually like the
// following:
//
//	accordion := gfmAccordionHeader("Accordion Title") + "Accordion Body" + gfmAccordionTerminator()
func gfmAccordionHeader(title string) string {
	return fmt.Sprintf("<details><summary>%s</summary>\n<p>\n\n", title)
}

// gfmAccordionTerminator generates the code necessary to terminate an
// accordion after the body. It is expected to be used in conjunction with
// gfmAccordionHeader(). See gfmAccordionHeader for a full description.
func gfmAccordionTerminator() string {
	return "</p>\n</details>\n\n"
}

// paragraph formats a paragraph with the provided text as the contents.
func paragraph(text string) string {
	return fmt.Sprintf("%s\n\n", escape(text))
}

var (
	specialCharacterRegex = regexp.MustCompile("([\\\\`*_{}\\[\\]()<>#+-.!~])")
	urlRegex              = xurls.Strict() // Require a scheme in URLs
)

func escape(text string) string {
	b := []byte(text)

	var (
		cursor  int
		builder strings.Builder
	)

	for _, urlLoc := range urlRegex.FindAllIndex(b, -1) {
		// Walk through each found URL, escaping the text before the URL and
		// leaving the text in the URL unchanged.
		if urlLoc[0] > cursor {
			// Escape the previous section if its length is nonzero
			builder.Write(escapeRaw(b[cursor:urlLoc[0]]))
		}

		// Add the unescaped URL to the end of it
		builder.Write(b[urlLoc[0]:urlLoc[1]])

		// Move the cursor forward for the next iteration
		cursor = urlLoc[1]
	}

	// Escape the end of the string after the last URL if there's anything left
	if len(b) > cursor {
		builder.Write(escapeRaw(b[cursor:]))
	}

	return builder.String()
}

func escapeRaw(segment []byte) []byte {
	return specialCharacterRegex.ReplaceAll(segment, []byte("\\$1"))
}

// plainText converts a markdown string to the plain text that appears in the
// rendered output.
func plainText(text string) string {
	md := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	node := md.Parse([]byte(text))

	var builder strings.Builder
	plainTextInner(node, &builder)

	return builder.String()
}

func plainTextInner(node *blackfriday.Node, builder *strings.Builder) {
	// Only text nodes produce output
	if node.Type == blackfriday.Text {
		builder.Write(node.Literal)
	}

	// Run the children first
	if node.FirstChild != nil {
		plainTextInner(node.FirstChild, builder)
	}

	// Then run any other siblings
	if node.Next != nil {
		// Add extra space if necessary between nodes
		if node.Type == blackfriday.Paragraph ||
			node.Type == blackfriday.CodeBlock ||
			node.Type == blackfriday.Heading {
			builder.WriteRune(' ')
		}

		plainTextInner(node.Next, builder)
	}
}
