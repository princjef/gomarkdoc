package format

import "github.com/princjef/gomarkdoc/lang"

// Format is a generic interface for formatting documentation contents in a
// particular way.
type Format interface {
	// Bold converts the provided text to bold
	Bold(text string) (string, error)

	// CodeBlock wraps the provided code as a code block and tags it with the
	// provided language (or no language if the empty string is provided).
	CodeBlock(language, code string) (string, error)

	// Header converts the provided text into a header of the provided level.
	// The level is expected to be at least 1.
	Header(level int, text string) (string, error)

	// RawHeader converts the provided text into a header of the provided level
	// without escaping the header text. The level is expected to be at least 1.
	RawHeader(level int, text string) (string, error)

	// LocalHref generates an href for navigating to a header with the given
	// headerText located within the same document as the href itself.
	LocalHref(headerText string) (string, error)

	// Link generates a link with the given text and href values.
	Link(text, href string) (string, error)

	// CodeHref generates an href to the provided code entry.
	CodeHref(loc lang.Location) (string, error)

	// ListEntry generates an unordered list entry with the provided text at the
	// provided zero-indexed depth. A depth of 0 is considered the topmost level
	// of list.
	ListEntry(depth int, text string) (string, error)

	// Accordion generates a collapsible content. The accordion's visible title
	// while collapsed is the provided title and the expanded content is the
	// body.
	Accordion(title, body string) (string, error)

	// AccordionHeader generates the header visible when an accordion is
	// collapsed.
	//
	// The AccordionHeader is expected to be used in conjunction with
	// AccordionTerminator() when the demands of the body's rendering requires
	// it to be generated independently. The result looks conceptually like the
	// following:
	//
	//	accordion := formatter.AccordionHeader("Accordion Title") + "Accordion Body" + formatter.AccordionTerminator()
	AccordionHeader(title string) (string, error)

	// AccordionTerminator generates the code necessary to terminate an
	// accordion after the body. It is expected to be used in conjunction with
	// AccordionHeader(). See AccordionHeader for a full description.
	AccordionTerminator() (string, error)

	// Paragraph formats a paragraph with the provided text as the contents.
	Paragraph(text string) (string, error)

	// Escape escapes special markdown characters from the provided text.
	Escape(text string) string
}
