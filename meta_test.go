package meta

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

func TestMeta(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Meta,
		),
	)
	source := `<!--
Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags: markdown, goldmark
-->

# Hello goldmark-meta
`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := context.Get(contextKey).(map[string]string)
	if metaData["Title"] != "goldmark-meta" {
		t.Error("Title not found in meta data or is not a string")
	}
	if metaData["Summary"] != "Add YAML metadata to the document" {
		t.Error("Summary not found in meta data or is not a string")
	}
	if metaData["Tags"] != "markdown, goldmark" {
		t.Error("Tags not found in meta data or is not a string")
	}

	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("should render '<h1>Hello goldmark-meta</h1>', but '%s'", buf.String())
	}

}
