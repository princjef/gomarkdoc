package format

import (
	"testing"

	"github.com/matryer/is"
)

func TestPlainText(t *testing.T) {
	tests := map[string]string{
		"plain text":                    "plain text",
		"[linked](https://foo.bar)":     "linked",
		"[linked 2](<https://foo.bar>)": "linked 2",
		"type [foo](<https://foo.bar>)": "type foo",
		"**bold** text":                 "bold text",
		"*italicized* text":             "italicized text",
		"~~strikethrough~~ text":        "strikethrough text",
		"paragraph 1\n\nparagraph 2":    "paragraph 1 paragraph 2",
		"# header\n\nparagraph":         "header paragraph",
	}

	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			is := is.New(t)
			is.Equal(plainText(in), out) // Wrong output for plainText()
		})
	}
}
