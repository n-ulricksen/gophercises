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

	for len(newLinks) > 0 {
		fmt.Println("QUEUE:", newLinks)
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
				sitemapLink = relToAbsURL(u.String(), urlInputString)
			}

			if sitemapLink != "" {
				newLinks = append(newLinks, sitemapLink)
				sitemapLinks = append(sitemapLinks, sitemapLink)
			}
		}

	}

	fmt.Println(sitemapLinks)

	// for site := range visitedLinks {
	// 	fmt.Println(site)
	// }

	// TODO:
	// (Optimize with goroutines)
	// Once no more new links can be found, build XML sitemap with found links
	// Output the built XML to stdout, or file specified by command-line flag
}

func relToAbsURL(relUrl string, baseUrl string) string {
	var absUrl string

	absUrl += baseUrl
	// Trim possible '/' at end of base URL
	if absUrl[len(absUrl)-1] == byte('/') {
		absUrl = absUrl[0 : len(absUrl)-1]
	}
	// Add '/' to beginning of path if necessary
	if relUrl[0] != byte('/') {
		absUrl += "/"
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
