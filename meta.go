// package meta is a extension for the goldmark(http://github.com/yuin/goldmark).
//
// This extension parses YAML metadata blocks and store metadata to a
// parser.Context.
package meta

import (
	"strings"

	"github.com/yuin/goldmark"
	ast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type metaParser struct {
}

var (
	contextKey        = parser.NewContextKey()
	defaultMetaParser = &metaParser{}
)

// NewParser returns a BlockParser that can parse YAML metadata blocks.
func NewParser() parser.BlockParser {
	return defaultMetaParser
}

func (b *metaParser) Trigger() []byte {
	return []byte("<!--")
}

func (b *metaParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	linenum, _ := reader.Position()
	if linenum != 0 {
		return nil, parser.NoChildren
	}
	line, _ := reader.PeekLine()
	if trim(string(line)) == "<!--" {
		return ast.NewTextBlock(), parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *metaParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if trim(string(line)) == "-->" {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *metaParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	data := make(map[string]string)
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		line := string(segment.Value(reader.Source()))
		if strings.ContainsAny(line, ":") {
			k, v := toEntry(line)
			data[k] = v
		}
	}

	pc.Set(contextKey, data)
	node.Parent().RemoveChild(node.Parent(), node)
}

func trim(str string) string {
	return strings.TrimSpace(str)
}

func toEntry(str string) (string, string) {
	kva := strings.Split(str, ":")
	return trim(kva[0]), trim(kva[1])
}

func (b *metaParser) CanInterruptParagraph() bool {
	return false
}

func (b *metaParser) CanAcceptIndentedLine() bool {
	return false
}

type meta struct {
}

// Meta is a extension for the goldmark.
var Meta = &meta{}

func (e *meta) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewParser(), 0),
		),
	)
}
