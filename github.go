package gomarkdoc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// GithubMarkdownFormatter provides a Formatter which is compatible with
// GitHub Flavored Markdown's syntax and semantics. See GitHub's documentation
// for more details about their markdown format:
// https://guides.github.com/features/mastering-markdown/
type GitHubMarkdownFormatter struct{}

// Bold converts the provided text to bold
func (f *GitHubMarkdownFormatter) Bold(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	return fmt.Sprintf("**%s**", escape(text)), nil
}

// CodeBlock wraps the provided code as a code block and tags it with the
// provided language (or no language if the empty string is provided).
func (f *GitHubMarkdownFormatter) CodeBlock(language, code string) (string, error) {
	return fmt.Sprintf("```%s\n%s\n```\n\n", language, strings.TrimSpace(code)), nil
}

// Header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func (f *GitHubMarkdownFormatter) Header(level int, text string) (string, error) {
	if level < 1 {
		return "", errors.New("formatter: header level cannot be less than 1")
	}

	switch level {
	case 1:
		return fmt.Sprintf("# %s\n\n", escape(text)), nil
	case 2:
		return fmt.Sprintf("## %s\n\n", escape(text)), nil
	case 3:
		return fmt.Sprintf("### %s\n\n", escape(text)), nil
	case 4:
		return fmt.Sprintf("#### %s\n\n", escape(text)), nil
	case 5:
		return fmt.Sprintf("##### %s\n\n", escape(text)), nil
	default:
		// Only go up to 6 levels. Anything higher is also level 6
		return fmt.Sprintf("###### %s\n\n", escape(text)), nil
	}
}

var (
	whitespaceRegex = regexp.MustCompile(`\s`)
	removeRegex     = regexp.MustCompile(`[^\pL-_\d]+`)
)

// LocalHref generates an href for navigating to a header with the given
// headerText located within the same document as the href itself.
func (f *GitHubMarkdownFormatter) LocalHref(headerText string) (string, error) {
	result := strings.ToLower(headerText)
	result = strings.TrimSpace(result)
	result = whitespaceRegex.ReplaceAllString(result, "-")
	result = removeRegex.ReplaceAllString(result, "")

	return fmt.Sprintf("#%s", result), nil
}

// Link generates a link with the given text and href values.
func (f *GitHubMarkdownFormatter) Link(text, href string) (string, error) {
	if text == "" {
		return "", nil
	}

	if href == "" {
		return text, nil
	}

	return fmt.Sprintf("[%s](<%s>)", text, href), nil
}

// ListEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func (f *GitHubMarkdownFormatter) ListEntry(depth int, text string) (string, error) {
	// TODO: this is a weird special case
	if text == "" {
		return "", nil
	}

	prefix := strings.Repeat("  ", depth)
	return fmt.Sprintf("%s- %s\n", prefix, text), nil
}

// Accordion generates a collapsible content. The accordion's visible title
// while collapsed is the provided title and the expanded content is the body.
func (f *GitHubMarkdownFormatter) Accordion(title, body string) (string, error) {
	return fmt.Sprintf("<details><summary>%s</summary>\n<p>\n\n%s</p>\n</details>\n\n", title, escape(body)), nil
}

// AccordionHeader generates the header visible when an accordion is collapsed.
//
// The AccordionHeader is expected to be used in conjunction with
// AccordionTerminator() when the demands of the body's rendering requires it to
// be generated independently. The result looks conceptually like the following:
//
//	accordion := formatter.AccordionHeader("Accordion Title") + "Accordion Body" + formatter.AccordionTerminator()
func (f *GitHubMarkdownFormatter) AccordionHeader(title string) string {
	return fmt.Sprintf("<details><summary>%s</summary>\n<p>\n\n", title)
}

// AccordionTerminator generates the code necessary to terminate an accordion
// after the body. It is expected to be used in conjunction with
// AccordionHeader(). See AccordionHeader for a full description.
func (f *GitHubMarkdownFormatter) AccordionTerminator() string {
	return "</p>\n</details>\n\n"
}

// Paragraph formats a paragraph with the provided text as the contents.
func (f *GitHubMarkdownFormatter) Paragraph(text string) (string, error) {
	return fmt.Sprintf("%s\n\n", escape(text)), nil
}

var specialCharacterRegex = regexp.MustCompile("([\\`*_{}\\[\\]()<>#+-.!])")

func escape(text string) string {
	return specialCharacterRegex.ReplaceAllString(text, "\\$1")
}
