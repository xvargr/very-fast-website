package vdoc

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/xvargr/very-fast-website/internal/logger"
	"golang.org/x/net/html"
)

type VirtualDocument struct {
	doc *goquery.Document
}

type VirtualExtractor struct {
	Meta         Meta
	HeadNodes    []*html.Node
	ContentNodes []*html.Node
}

type Meta struct {
	Title string
}

func NewVirtualDocument() *VirtualDocument {
	root := &html.Node{
		Type: html.DocumentNode,
	}

	root.AppendChild(&html.Node{
		Type: html.DoctypeNode,
		Data: "html",
	})

	htmlNode := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		Attr: []html.Attribute{
			{Key: "lang", Val: "en"},
		},
	}

	root.AppendChild(htmlNode)

	headNode := &html.Node{
		Type: html.ElementNode,
		Data: "head",
	}

	headNode.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "title",
		FirstChild: &html.Node{
			Type: html.TextNode,
			Data: "Untitled",
		},
	})

	htmlNode.AppendChild(headNode)

	htmlNode.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "body",
	})

	htmlNode.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "outlet",
	})

	return &VirtualDocument{
		doc: goquery.NewDocumentFromNode(root),
	}
}

func (doc *VirtualDocument) RenderHtml() string {
	html, err := doc.doc.Html()
	if err != nil {
		logger.Console(logger.SeverityError, "failed to generate html "+err.Error())
	}

	return html
}

func Extract(raw []byte) *VirtualExtractor {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(raw))
	if err != nil {
		logger.Console(logger.SeverityError, "failed to parse html "+err.Error())
	}

	head := doc.Find("head").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Nodes[0].Data != "title"
	})
	body := doc.Find("body").Contents()

	// fmt.Println("head", len(head.Nodes))
	// fmt.Println("body", len(body.Nodes))

	return &VirtualExtractor{
		Meta: Meta{
			Title: doc.Find("title").Text(),
		},
		HeadNodes:    head.Nodes,
		ContentNodes: body.Nodes,
	}
}

func cloneNode(n *html.Node) *html.Node {
	newNode := &html.Node{
		Type:     n.Type,
		Data:     n.Data,
		DataAtom: n.DataAtom,
		Attr:     append([]html.Attribute(nil), n.Attr...), // Deep copy attributes
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		newNode.AppendChild(cloneNode(child)) // Recursively clone children
	}
	return newNode
}

func (d *VirtualDocument) Merge(extr *VirtualExtractor) {
	// fmt.Println("merge b4", d.RenderHtml())
	d.doc.Find("outlet").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithNodes(extr.ContentNodes...)
	})
	d.doc.Find("head").Each(func(i int, s *goquery.Selection) {
		s.AppendNodes(extr.HeadNodes...)
	})
	if extr.Meta.Title != "" {
		d.doc.Find("title").Get(0).FirstChild.Data = extr.Meta.Title
	}
	// fmt.Println("merge after", d.RenderHtml())
}
