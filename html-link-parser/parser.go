package linkparser

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseHTMLLinks(r io.Reader) []Link {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
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

	return links
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
	var textElements []string

	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			// Check the length of the trimmed string to see if the string
			// contains any meaningful characters.
			trimmedText := strings.TrimSpace(child.Data)
			if child.Type == html.TextNode && len(trimmedText) > 0 {
				// This is the link's text!
				textElements = append(textElements, trimmedText)
			}
			dfs(child)
		}
	}
	dfs(n)

	return strings.Join(textElements, " ")
}
