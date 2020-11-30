package format

import (
	"testing"

	"github.com/matryer/is"
)

func TestPlainText(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{
			in:  "plain text",
			out: "plain text",
		},
		{
			in:  "[linked](https://foo.bar)",
			out: "linked",
		},
		{
			in:  "[linked 2](<https://foo.bar>)",
			out: "linked 2",
		},
		{
			in:  "type [foo](<https://foo.bar>)",
			out: "type foo",
		},
		{
			in:  "**bold** text",
			out: "bold text",
		},
		{
			in:  "*italicized* text",
			out: "italicized text",
		},
		{
			in:  "~~strikethrough~~ text",
			out: "strikethrough text",
		},
		{
			in:  "paragraph 1\n\nparagraph 2",
			out: "paragraph 1 paragraph 2",
		},
		{
			in:  "# header\n\nparagraph",
			out: "header paragraph",
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			is := is.New(t)
			is.Equal(plainText(test.in), test.out) // Wrong output for plainText()
		})
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{
			in:  "plain text",
			out: `plain text`,
		},
		{
			in:  "**bold** text",
			out: `\*\*bold\*\* text`,
		},
		{
			in:  "*italicized* text",
			out: `\*italicized\* text`,
		},
		{
			in:  "~~strikethrough~~ text",
			out: `\~\~strikethrough\~\~ text`,
		},
		{
			in:  "# header",
			out: `\# header`,
		},
		{
			in:  "markdown [link](https://foo.bar)",
			out: `markdown \[link\]\(https://foo.bar\)`,
		},
		{
			in: "# header then complex URL: http://abc.def/sdfklj/sdf?key=value&special=%323%20sd " +
				"with http://simple.url and **bold** after",
			out: `\# header then complex URL: http://abc.def/sdfklj/sdf?key=value&special=%323%20sd ` +
				`with http://simple.url and \*\*bold\*\* after`,
		},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			is := is.New(t)
			is.Equal(escape(test.in), test.out) // Wrong output for escape()
		})
	}
}
