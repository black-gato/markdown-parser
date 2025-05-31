package internal

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// TODO: Need to store each
// TODO:
// TODO:

func getNodeText(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	// Traverse all sibling text nodes directly within the list item
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok && textNode.Segment.Len() > 0 {
			buf.Write(textNode.Segment.Value(source))
		} else if child.Kind() == ast.KindParagraph || child.Kind() == ast.KindTextBlock {
			// Collect text from all Text nodes within Paragraph or TextBlock
			for textChild := child.FirstChild(); textChild != nil; textChild = textChild.NextSibling() {
				if textNode, ok := textChild.(*ast.Text); ok && textNode.Segment.Len() > 0 {
					buf.Write(textNode.Segment.Value(source))
				}
			}
		}
	}
	return strings.TrimSpace(buf.String())
}
func MarkdownTicks(level int) string {
	if level < 0 {
		return ""

	}
	return strings.Repeat("    ", level)

}
func PageLink(file string) string {
	header := fmt.Sprintf("### [%s](%s)\n", file, file)

	return header
}
func GatherContent(n ast.Node, level int, source []byte, tag, content string) string {

	if n.Kind() == ast.KindListItem {
		t := MarkdownTicks(level)
		content += fmt.Sprintf("%s- %s\n", t, getNodeText(n, source))

	}
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {

		if child.Kind() == ast.KindList {
			content = GatherContent(child, level+1, source, tag, content)

		} else {
			content = GatherContent(child, level, source, tag, content)

		}

	}
	return content
}

func SearchTag(content []byte, tag string) []ast.Node {
	var liNodes []ast.Node
	md := goldmark.New()
	reader := text.NewReader(content)
	doc := md.Parser().Parse(reader)
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() == ast.KindListItem && entering {
			if strings.Contains(getNodeText(n, content), tag) {
				liNodes = append(liNodes, n)
				return ast.WalkSkipChildren, nil
			}
		}
		return ast.WalkContinue, nil
	})

	return liNodes

}

func Parse(files []string, tags []string) (reference string, err error) {
	var a string
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return "", nil
		}

		reference += PageLink(f)
		for _, tag := range tags {
			nodes := SearchTag(content, tag)
			for _, n := range nodes {
				reference += GatherContent(n, 0, content, tag, a)

			}

		}

	}

	return reference, nil
}
