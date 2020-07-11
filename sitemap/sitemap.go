package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ulricksennick/gophercises/linkparser"
)

func main() {
	// Get sitemap URL from command line arguments
	args := os.Args[1:]
	if len(args) != 1 {
		printErrorMessage()
		os.Exit(1)
	}
	urlInputString := args[0]

	// Parse the user specified URL, make HTTP request
	parsedInputUrl, err := url.Parse(urlInputString)
	if err != nil {
		log.Fatal("URL parse error:", err)
	}
	httpBody := getHTTPBody(urlInputString)
	defer httpBody.Close()

	// Find all links on the page
	links, err := linkparser.ParseHTMLLinks(httpBody)
	if err != nil {
		log.Fatal(err)
	}

	// Check each link to see if it is in the sitemap domain
	sitemapLinks := filterLinksByDomain(links, parsedInputUrl)
	for _, l := range sitemapLinks {
		fmt.Println(l)
	}

	// TODO:
	// Repeatedly parse all new sitemap links for new links
	// Once no more new links can be found, build XML sitemap with found links
	// Output the built XML to stdout, or file specified by command-line flag
}

func filterLinksByDomain(links []linkparser.Link, domainURL *url.URL) []string {
	var filteredLinks []string
	for _, link := range links {
		u, err := url.Parse(link.Href)
		if err != nil {
			log.Fatal(err)
		}

		// Skip anchor-tag links
		if isAnchorLink(u) {
			continue
		}

		if isRelativeLink(u) {
			var absLink string

			absLink += domainURL.String()
			// Trim possible '/' at end of base URL
			if absLink[len(absLink)-1] == byte('/') {
				absLink = absLink[0 : len(absLink)-1]
			}
			// Add '/' to beginning of path if necessary
			if u.String()[0] != byte('/') {
				absLink += "/"
			}
			absLink += u.String()

			filteredLinks = append(filteredLinks, absLink)
			continue
		}

		if getDomainFromURL(u) == getDomainFromURL(domainURL) {
			filteredLinks = append(filteredLinks, u.String())
		}
	}

	return filteredLinks
}

func getHTTPBody(urlString string) io.ReadCloser {
	resp, err := http.Get(urlString)
	if err != nil {
		log.Fatal("URL parse error:", err)
	}
	return resp.Body
}

func isAnchorLink(u *url.URL) bool {
	return u.Fragment != ""
}

// Determine if the link is a relative link to a site on the same domain
func isRelativeLink(u *url.URL) bool {
	return !u.IsAbs() && u.Hostname() == ""
}

// Get the domain name from a url.URL
func getDomainFromURL(u *url.URL) string {
	if u.Hostname() == "" {
		return ""
	}

	urlParts := strings.Split(u.Hostname(), ".")
	l := len(urlParts)

	return strings.Join(urlParts[l-2:], ".")
}

func printErrorMessage() {
	fmt.Println("Usage: ./sitemap <url>")
	fmt.Println()
	fmt.Println("This program builds a sitemap for the specified URL's domain")
	fmt.Println("using the standard sitemap protocol.")
}
