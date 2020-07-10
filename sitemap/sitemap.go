package main

import (
	"fmt"
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

	// Parse the user specified URL, get the sitemap domain
	parsedUrl, err := url.Parse(urlInputString)
	if err != nil {
		log.Fatal("URL parse error:", err)
	}
	sitemapDomain := getDomainFromURL(parsedUrl)
	fmt.Println("SITEMAP DOMAIN:", sitemapDomain)

	// Make HTTP request, get the message body
	resp, err := http.Get(urlInputString)
	if err != nil {
		log.Fatal("URL parse error:", err)
	}
	respReader := resp.Body
	defer respReader.Close()

	// Find all links on the page
	links, err := linkparser.ParseHTMLLinks(respReader)
	if err != nil {
		log.Fatal(err)
	}

	// Check each link to see if it is in the sitemap domain
	var sitemapLinks []string
	for _, link := range links {
		u, err := url.Parse(link.Href)
		if err != nil {
			log.Fatal(err)
		}

		if isRelativeLink(u) || getDomainFromURL(u) == sitemapDomain {
			sitemapLinks = append(sitemapLinks, link.Href)
			continue
		}
	}
	for _, l := range sitemapLinks {
		fmt.Println(l)
	}

	// TODO:
	// Repeatedly parse all new sitemap links for new links
	// Once no more new links can be found, build XML sitemap with found links
	// Output the built XML to stdout, or file specified by command-line flag
}

// Determine if the link is a relative link to a site on the same domain
func isRelativeLink(u *url.URL) bool {
	return !u.IsAbs() && u.Hostname() == ""
}

// Get the domain name from a url.URL
func getDomainFromURL(u *url.URL) string {
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
