package format

import (
	"fmt"

	"github.com/princjef/gomarkdoc/format/formatcore"
	"github.com/princjef/gomarkdoc/lang"
)

// PlainMarkdown provides a Format which is compatible with the base Markdown
// format specification.
type PlainMarkdown struct{}

// Bold converts the provided text to bold
func (f *PlainMarkdown) Bold(text string) (string, error) {
	return formatcore.Bold(text), nil
}

// CodeBlock wraps the provided code as a code block. The provided language is
// ignored as it is not supported in plain markdown.
func (f *PlainMarkdown) CodeBlock(language, code string) (string, error) {
	return formatcore.CodeBlock(code), nil
}

// Header converts the provided text into a header of the provided level. The
// level is expected to be at least 1.
func (f *PlainMarkdown) Header(level int, text string) (string, error) {
	return formatcore.Header(level, formatcore.Escape(text))
}

// RawHeader converts the provided text into a header of the provided level
// without escaping the header text. The level is expected to be at least 1.
func (f *PlainMarkdown) RawHeader(level int, text string) (string, error) {
	return formatcore.Header(level, text)
}

// LocalHref always returns the empty string, as header links are not supported
// in plain markdown.
func (f *PlainMarkdown) LocalHref(headerText string) (string, error) {
	return "", nil
}

// CodeHref always returns the empty string, as there is no defined file linking
// format in standard markdown.
func (f *PlainMarkdown) CodeHref(loc lang.Location) (string, error) {
	return "", nil
}

// Link generates a link with the given text and href values.
func (f *PlainMarkdown) Link(text, href string) (string, error) {
	return formatcore.Link(text, href), nil
}

// ListEntry generates an unordered list entry with the provided text at the
// provided zero-indexed depth. A depth of 0 is considered the topmost level of
// list.
func (f *PlainMarkdown) ListEntry(depth int, text string) (string, error) {
	return formatcore.ListEntry(depth, text), nil
}

// Accordion generates a collapsible content. Since accordions are not supported
// by plain markdown, this generates a level 6 header followed by a paragraph.
func (f *PlainMarkdown) Accordion(title, body string) (string, error) {
	h, err := formatcore.Header(6, title)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", h, formatcore.Paragraph(body)), nil
}

// AccordionHeader generates the header visible when an accordion is collapsed.
// Since accordions are not supported in plain markdown, this generates a level
// 6 header.
//
// The AccordionHeader is expected to be used in conjunction with
// AccordionTerminator() when the demands of the body's rendering requires it to
// be generated independently. The result looks conceptually like the following:
//
//	accordion := format.AccordionHeader("Accordion Title") + "Accordion Body" + format.AccordionTerminator()
func (f *PlainMarkdown) AccordionHeader(title string) (string, error) {
	return formatcore.Header(6, title)
}

// AccordionTerminator generates the code necessary to terminate an accordion
// after the body. Since accordions are not supported in plain markdown, this
// completes a paragraph section. It is expected to be used in conjunction with
// AccordionHeader(). See AccordionHeader for a full description.
func (f *PlainMarkdown) AccordionTerminator() (string, error) {
	return "\n\n", nil
}

// Paragraph formats a paragraph with the provided text as the contents.
func (f *PlainMarkdown) Paragraph(text string) (string, error) {
	return formatcore.Paragraph(text), nil
}

// Escape escapes special markdown characters from the provided text.
func (f *PlainMarkdown) Escape(text string) string {
	return formatcore.Escape(text)
}
