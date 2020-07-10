package linkparser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseHTMLLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var links []Link

	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// This is a link!
			newLink := createLink(n)
			links = append(links, newLink)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(doc)

	return links, nil
}

func createLink(n *html.Node) Link {
	var link Link

	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			link.Text = getLinkText(n)
		}
	}

	return link
}

func getLinkText(n *html.Node) string {
	var text string

	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				// This is the link's text!
				text += child.Data
			}
			dfs(child)
		}
	}
	dfs(n)

	return strings.Join(strings.Fields(text), " ")
}
