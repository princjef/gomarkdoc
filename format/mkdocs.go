package format

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/princjef/gomarkdoc/format/formatcore"
	"github.com/princjef/gomarkdoc/format/repo_format"
	"github.com/princjef/gomarkdoc/lang"
)

// MkDocsMarkdown provides a Format which is compatible with BitBucket
// Markdown's syntax and semantics
type MkDocsMarkdown struct {
	RepoFormats []repo_format.RepoFormat
}

// Bold converts the provided text to bold
func (f *MkDocsMarkdown) Bold(text string) (string, error) {
	return formatcore.Bold(text), nil
}

// CodeBlock wraps the provided code as a code block and tags it with the
// provided language (or no language if the empty string is provided).
func (f *MkDocsMarkdown) CodeBlock(language, code string) (string, error) {
	return formatcore.GFMCodeBlock(language, code), nil
}

// Anchor produces an anchor for the provided link.
func (f *MkDocsMarkdown) Anchor(anchor string) string {
	return formatcore.Anchor(anchor)
}

// AnchorHeader converts the provided text and custom anchor link into a header
// of the provided level. The level is expected to be at least 1.
func (f *MkDocsMarkdown) AnchorHeader(level int, text, anchor string) (string, error) {
	return formatcore.AnchorHeader(level, f.Escape(text), anchor)
}

// Header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func (f *MkDocsMarkdown) Header(level int, text string) (string, error) {
	return formatcore.Header(level, f.Escape(text))
}

// RawAnchorHeader converts the provided text and custom anchor link into a
// header of the provided level without escaping the header text. The level is
// expected to be at least 1.
func (f *MkDocsMarkdown) RawAnchorHeader(level int, text, anchor string) (string, error) {
	return formatcore.AnchorHeader(level, text, anchor)
}

// RawHeader converts the provided text into a header of the provided level
// without escaping the header text. The level is expected to be at least 1.
func (f *MkDocsMarkdown) RawHeader(level int, text string) (string, error) {
	return formatcore.Header(level, text)
}

var (
	mkDocsWhitespaceRegex = regexp.MustCompile(`\s`)
	mkDocsRemoveRegex     = regexp.MustCompile(`[^\pL-_\d]+`)
)

// LocalHref generates an href for navigating to a header with the given
// headerText located within the same document as the href itself.
func (f *MkDocsMarkdown) LocalHref(headerText string) (string, error) {
	result := formatcore.PlainText(headerText)
	result = strings.ToLower(result)
	result = strings.TrimSpace(result)
	result = mkDocsWhitespaceRegex.ReplaceAllString(result, "-")
	result = mkDocsRemoveRegex.ReplaceAllString(result, "")

	return fmt.Sprintf("#%s", result), nil
}

// RawLocalHref generates an href within the same document but with a direct
// link provided instead of text to slugify.
func (f *MkDocsMarkdown) RawLocalHref(anchor string) string {
	return fmt.Sprintf("#%s", anchor)
}

// Link generates a link with the given text and href values.
func (f *MkDocsMarkdown) Link(text, href string) (string, error) {
	return formatcore.Link(text, href), nil
}

// CodeHref generates an href to the provided code entry.
func (f *MkDocsMarkdown) CodeHref(loc lang.Location) (string, error) {
	// If there's no repo, we can't compute an href
	if loc.Repo == nil {
		return "", nil
	}

	for _, repoFormat := range f.RepoFormats {
		if repoFormat.Supports(loc.Repo.Type) {
			return repoFormat.CodeHref(loc)
		}
	}

	// if there is no any supported repo format
	return "", nil
}

// ListEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func (f *MkDocsMarkdown) ListEntry(depth int, text string) (string, error) {
	return formatcore.ListEntry(depth, text), nil
}

// Accordion generates a collapsible content. The accordion's visible title
// while collapsed is the provided title and the expanded content is the body.
func (f *MkDocsMarkdown) Accordion(title, body string) (string, error) {
	return formatcore.GFMAccordion(title, body), nil
}

// AccordionHeader generates the header visible when an accordion is collapsed.
//
// The AccordionHeader is expected to be used in conjunction with
// AccordionTerminator() when the demands of the body's rendering requires it to
// be generated independently. The result looks conceptually like the following:
//
//	accordion := format.AccordionHeader("Accordion Title") + "Accordion Body" + format.AccordionTerminator()
func (f *MkDocsMarkdown) AccordionHeader(title string) (string, error) {
	return formatcore.GFMAccordionHeader(title), nil
}

// AccordionTerminator generates the code necessary to terminate an accordion
// after the body. It is expected to be used in conjunction with
// AccordionHeader(). See AccordionHeader for a full description.
func (f *MkDocsMarkdown) AccordionTerminator() (string, error) {
	return formatcore.GFMAccordionTerminator(), nil
}

// Escape escapes special markdown characters from the provided text.
func (f *MkDocsMarkdown) Escape(text string) string {
	return formatcore.Escape(text)
}
