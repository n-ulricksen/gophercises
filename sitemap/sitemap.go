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

	// Verify starting URL is valid
	parsedInputUrl, err := url.Parse(urlInputString)
	if err != nil {
		log.Fatal("URL parse error:", err)
	}

	// Check each link to see if it is in the sitemap domain
	visitedLinks := map[string]bool{urlInputString: true}
	newLinks := []string{urlInputString}
	var sitemapLinks []string

	fmt.Println("INPUT URL:", parsedInputUrl)

	for len(newLinks) > 0 {
		link := newLinks[0]
		newLinks = newLinks[1:]

		// Parse the user specified URL, make HTTP request
		currentUrl, parseErr := url.Parse(link)
		httpBody, getErr := getHTTPBody(currentUrl.String())
		// Skip invalid links or pages requiring authentication
		if parseErr != nil || getErr != nil {
			continue
		}
		defer httpBody.Close()

		fmt.Println("-x- Parsing URL:", currentUrl.String())

		// Find all links on the page
		possibleLinks, err := linkparser.ParseHTMLLinks(httpBody)
		if err != nil {
			log.Fatal(err)
		}

		for _, l := range possibleLinks {
			u, err := url.Parse(l.Href)
			// Skip any invalid urls
			if err != nil {
				continue
			}

			// Check if link has been visited
			if visitedLinks[u.String()] {
				continue
			}
			visitedLinks[u.String()] = true

			// Do not include empty links or anchor tag links
			if u.String() == "" || isAnchorLink(u) {
				continue
			}

			var sitemapLink string
			if linkInDomain(u, parsedInputUrl) {
				sitemapLink = u.String()
			} else if isRelativeLink(u) {
				sitemapLink = relToAbsURL(u.String(), currentUrl.String())
			}

			if sitemapLink != "" {
				newLinks = append(newLinks, sitemapLink)
				sitemapLinks = append(sitemapLinks, sitemapLink)
				fmt.Println("\t-x- Adding Link:", sitemapLink)
			}
		}
	}

	// fmt.Println(sitemapLinks)

	// TODO:
	// (Optimize with goroutines)
	// Use linked list for bfs newLinks queue if necessary - "container/list"
	// Once no more new links can be found, build XML sitemap with found links
	// Output the built XML to stdout, or file specified by command-line flag
}

func trimURLPath(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	trimmed := fmt.Sprintf("%v://%v/", u.Scheme, u.Hostname())
	return trimmed, nil
}

func relToAbsURL(relUrl string, searchUrl string) string {
	var absUrl string

	baseUrl, err := trimURLPath(searchUrl)
	if err != nil {
		log.Print("ERROR:", err)
	}

	// Make sure searchUrl ends with '/'
	if searchUrl[len(searchUrl)-1] != byte('/') {
		searchUrl += "/"
	}

	// If the relative path begins with '/', it is specifies a path relative to
	// the baseUrl
	if relUrl[0] == byte('/') {
		relUrl = relUrl[1:]
		absUrl += baseUrl
	} else {
		absUrl += searchUrl
	}
	absUrl += relUrl

	return absUrl
}

func linkInDomain(u *url.URL, domainURL *url.URL) bool {
	if getDomainFromURL(u) == getDomainFromURL(domainURL) {
		return true
	}

	return false
}

func getHTTPBody(urlString string) (io.ReadCloser, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
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
	if l < 2 {
		return ""
	}

	return strings.Join(urlParts[l-2:], ".")
}

func printErrorMessage() {
	fmt.Println("Usage: ./sitemap <url>")
	fmt.Println()
	fmt.Println("This program builds a sitemap for the specified URL's domain")
	fmt.Println("using the standard sitemap protocol.")
}
