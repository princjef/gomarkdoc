package format

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

var specialCharacterRegex = regexp.MustCompile("([\\\\`*_{}\\[\\]()<>#+-.!])")

func escape(text string) string {
	return specialCharacterRegex.ReplaceAllString(text, "\\$1")
}
