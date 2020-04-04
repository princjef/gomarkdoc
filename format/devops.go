package format

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/princjef/gomarkdoc/lang"
)

// AzureDevOpsMarkdown provides a Format which is compatible with Azure
// DevOps's syntax and semantics. See the Azure DevOps documentation for more
// details about their markdown format:
// https://docs.microsoft.com/en-us/azure/devops/project/wiki/markdown-guidance?view=azure-devops
type AzureDevOpsMarkdown struct{}

// Bold converts the provided text to bold
func (f *AzureDevOpsMarkdown) Bold(text string) (string, error) {
	return bold(text), nil
}

// CodeBlock wraps the provided code as a code block and tags it with the
// provided language (or no language if the empty string is provided).
func (f *AzureDevOpsMarkdown) CodeBlock(language, code string) (string, error) {
	return gfmCodeBlock(language, code), nil
}

// Header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func (f *AzureDevOpsMarkdown) Header(level int, text string) (string, error) {
	return header(level, escape(text))
}

// RawHeader converts the provided text into a header of the provided level
// without escaping the header text. The level is expected to be at least 1.
func (f *AzureDevOpsMarkdown) RawHeader(level int, text string) (string, error) {
	return header(level, text)
}

var devOpsWhitespaceRegex = regexp.MustCompile(`\s`)

// LocalHref generates an href for navigating to a header with the given
// headerText located within the same document as the href itself. Link
// generation follows the guidelines here:
// https://docs.microsoft.com/en-us/azure/devops/project/wiki/markdown-guidance?view=azure-devops#anchor-links
func (f *AzureDevOpsMarkdown) LocalHref(headerText string) (string, error) {
	result := strings.ToLower(headerText)
	result = strings.TrimSpace(result)
	result = devOpsWhitespaceRegex.ReplaceAllString(result, "-")
	result = url.PathEscape(result)
	// We also have to escape the `:` character if present
	result = strings.ReplaceAll(result, ":", "%3A")

	return fmt.Sprintf("#%s", result), nil
}

// CodeHref generates an href to the provided code entry.
func (f *AzureDevOpsMarkdown) CodeHref(loc lang.Location) (string, error) {
	// If there's no repo, we can't compute an href
	if loc.Repo == nil {
		return "", nil
	}

	var (
		relative string
		err      error
	)
	if filepath.IsAbs(loc.Filepath) {
		relative, err = filepath.Rel(loc.WorkDir, loc.Filepath)
		if err != nil {
			return "", err
		}
	} else {
		relative = loc.Filepath
	}

	full := filepath.Join(loc.Repo.PathFromRoot, relative)
	p, err := filepath.Rel(string(filepath.Separator), full)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s?path=%s&version=GB%s&lineStyle=plain&line=%d&lineEnd=%d&lineStartColumn=%d&lineEndColumn=%d",
		loc.Repo.Remote,
		url.PathEscape(filepath.ToSlash(p)),
		loc.Repo.DefaultBranch,
		loc.Start.Line,
		loc.End.Line,
		loc.Start.Col,
		loc.End.Col,
	), nil
}

// Link generates a link with the given text and href values.
func (f *AzureDevOpsMarkdown) Link(text, href string) (string, error) {
	return link(text, href), nil
}

// ListEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func (f *AzureDevOpsMarkdown) ListEntry(depth int, text string) (string, error) {
	return listEntry(depth, text), nil
}

// Accordion generates a collapsible content. The accordion's visible title
// while collapsed is the provided title and the expanded content is the body.
func (f *AzureDevOpsMarkdown) Accordion(title, body string) (string, error) {
	return gfmAccordion(title, body), nil
}

// AccordionHeader generates the header visible when an accordion is collapsed.
//
// The AccordionHeader is expected to be used in conjunction with
// AccordionTerminator() when the demands of the body's rendering requires it to
// be generated independently. The result looks conceptually like the following:
//
//	accordion := format.AccordionHeader("Accordion Title") + "Accordion Body" + format.AccordionTerminator()
func (f *AzureDevOpsMarkdown) AccordionHeader(title string) (string, error) {
	return gfmAccordionHeader(title), nil
}

// AccordionTerminator generates the code necessary to terminate an accordion
// after the body. It is expected to be used in conjunction with
// AccordionHeader(). See AccordionHeader for a full description.
func (f *AzureDevOpsMarkdown) AccordionTerminator() (string, error) {
	return gfmAccordionTerminator(), nil
}

// Paragraph formats a paragraph with the provided text as the contents.
func (f *AzureDevOpsMarkdown) Paragraph(text string) (string, error) {
	return paragraph(text), nil
}

// Escape escapes special markdown characters from the provided text.
func (f *AzureDevOpsMarkdown) Escape(text string) string {
	return escape(text)
}
