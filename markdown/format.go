package markdown

import (
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

// Format formats Markdown.
func Format(src []byte, options ...renderer.Option) ([]byte, error) {
	ctx := parser.NewContext()
	parsed := myParser.Parse(text.NewReader(src), parser.WithContext(ctx))

	return Render(ctx, src, parsed, options...)
}
