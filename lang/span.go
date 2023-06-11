package lang

import (
	"fmt"
	"go/doc/comment"
	"regexp"
	"strings"
)

type (
	// Span defines a single text span in a block for documentation of a symbol
	// or package.
	Span struct {
		cfg  *Config
		kind SpanKind
		text string
		url  string
	}

	// SpanKind identifies the type of span element represented by the
	// corresponding Span.
	SpanKind string
)

const (
	// TextSpan defines a span that represents plain text.
	TextSpan SpanKind = "text"

	// RawTextSpan defines a span that represents plain text that should be
	// displayed as-is.
	RawTextSpan SpanKind = "rawText"

	// LinkSpan defines a span that represents a link.
	LinkSpan SpanKind = "link"

	// AutolinkSpan defines a span that represents text which is itself a link.
	AutolinkSpan SpanKind = "autolink"
)

// NewSpan creates a new span.
func NewSpan(cfg *Config, kind SpanKind, text string, url string) *Span {
	return &Span{cfg, kind, text, url}
}

// Kind provides the kind of data that this span represents.
func (s *Span) Kind() SpanKind {
	return s.kind
}

// Text provides the raw text for the span.
func (s *Span) Text() string {
	return s.text
}

// URL provides the url associated with the span, if any.
func (s *Span) URL() string {
	return s.url
}

// ParseSpans turns a set of *comment.Text entries into a slice of spans.
func ParseSpans(cfg *Config, texts []comment.Text) []*Span {
	var s []*Span
	for _, t := range texts {
		switch v := t.(type) {
		case comment.Plain:
			s = append(s, NewSpan(cfg.Inc(0), TextSpan, collapseWhitespace(string(v)), ""))
		case comment.Italic:
			s = append(s, NewSpan(cfg.Inc(0), TextSpan, collapseWhitespace(string(v)), ""))
		case *comment.DocLink:
			var b strings.Builder
			printText(&b, v.Text...)
			str := collapseWhitespace(b.String())

			// Replace local links as needed
			if v.ImportPath == "" {
				name := symbolName(v.Recv, v.Name)
				if sym, ok := cfg.Symbols[name]; ok {
					s = append(s, NewSpan(cfg.Inc(0), LinkSpan, str, fmt.Sprintf("#%s", sym.Anchor())))
				} else {
					cfg.Log.Warnf("Unable to find symbol %s", name)
					s = append(s, NewSpan(cfg.Inc(0), TextSpan, collapseWhitespace(str), ""))
				}
				break
			}

			s = append(s, NewSpan(cfg.Inc(0), LinkSpan, str, v.DefaultURL("https://pkg.go.dev/")))
		case *comment.Link:
			var b strings.Builder
			printText(&b, v.Text...)
			str := collapseWhitespace(b.String())

			if v.Auto {
				s = append(s, NewSpan(cfg.Inc(0), AutolinkSpan, str, str))
			}

			s = append(s, NewSpan(cfg.Inc(0), LinkSpan, str, v.URL))
		}
	}

	return s
}

func printText(b *strings.Builder, text ...comment.Text) {
	for i, t := range text {
		if i > 0 {
			b.WriteRune(' ')
		}

		switch v := t.(type) {
		case comment.Plain:
			b.WriteString(string(v))
		case comment.Italic:
			b.WriteString(string(v))
		case *comment.DocLink:
			printText(b, v.Text...)
		case *comment.Link:
			// No need to linkify implicit links
			if v.Auto {
				printText(b, v.Text...)
				continue
			}

			b.WriteRune('[')
			printText(b, v.Text...)
			b.WriteRune(']')
			b.WriteRune('(')
			b.WriteString(v.URL)
			b.WriteRune(')')
		}
	}
}

var whitespaceRegex = regexp.MustCompile(`\s+`)

func collapseWhitespace(s string) string {
	return string(whitespaceRegex.ReplaceAll([]byte(s), []byte(" ")))
}
