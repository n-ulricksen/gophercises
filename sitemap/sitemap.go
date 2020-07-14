package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ulricksennick/gophercises/linkparser"
)

func main() {
	// File to store the generated sitemap upon completion
	var xmlFilename string = "sitemap.xml"

	// Flags
	filename := flag.String("file", xmlFilename,
		"used to specify file in to store XML sitemap to")
	searchDepth := flag.Int("depth", 0, "how many layers deep to search for links")
	helpFlag := flag.Bool("help", false, "display the help message")
	flag.Parse()

	// Get sitemap URL from command line arguments
	args := flag.Args()
	if len(args) != 1 {
		printHelpMessage()
		os.Exit(1)
	}
	if *helpFlag {
		printHelpMessage()
		os.Exit(0)
	}
	urlInputString := args[0]

	// Verify starting URL is valid
	parsedInputUrl, err := url.Parse(urlInputString)
	if err != nil {
		log.Fatal("URL parse error:", err)
	}
	if parsedInputUrl.Scheme == "" || parsedInputUrl.Hostname() == "" {
		fmt.Println("Invalid URL...")
		os.Exit(1)
	}

	// Get sitemap links from input URL by traversing links using BFS, stopping
	// when no new links can be found
	sitemapLinks := getSitemapLinks(parsedInputUrl, *searchDepth)

	// Once no more new links can be found, build XML sitemap with found links
	xmlBytes := generateSitemapXML(sitemapLinks)

	err = ioutil.WriteFile(*filename, xmlBytes, 0664)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sitemap creation complete... file saved at %v.\n", *filename)

	// TODO:
	// Use linked list for bfs newLinks queue if necessary - "container/list"
	// (Optimize with goroutines)
}

type sitemapXML struct {
	XMLName   xml.Name          `xml:"urlset"`
	Locations []sitemapLocation `xml:"url"`
}

type sitemapLocation struct {
	Location string `xml:"loc"`
}

func generateSitemapXML(links []string) []byte {
	xmlBytes := []byte(xml.Header)

	sitemap := &sitemapXML{}
	for _, link := range links {
		sitemap.Locations = append(sitemap.Locations, sitemapLocation{link})
	}

	output, err := xml.MarshalIndent(sitemap, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	xmlBytes = append(xmlBytes, output...)

	return xmlBytes
}

// TODO: Break apart this monster function
func getSitemapLinks(parsedInputUrl *url.URL, depth int) []string {
	// Check each link to see if it is in the sitemap domain
	visitedLinks := map[string]bool{parsedInputUrl.String(): true}
	queue := []string{}
	nextQueue := []string{parsedInputUrl.String()}
	var sitemapLinks []string

	fmt.Println("Creating sitemap from URL:", parsedInputUrl)

	for i := 0; i <= depth; i++ {
		// fmt.Println("-x- DEPTH:", i)
		queue = nextQueue
		for len(queue) > 0 {
			currentLink := queue[0]
			queue = queue[1:]

			// Parse the user specified URL, make HTTP request
			currentUrl, parseErr := url.Parse(currentLink)
			httpBody, httpErr := getHTTPBody(currentUrl.String())
			// Skip invalid links or pages requiring authentication
			if parseErr != nil || httpErr != nil {
				continue
			}
			defer httpBody.Close()

			// fmt.Println("-x- Parsing URL:", currentUrl.String())

			// Find all links on the page
			linksOnPage, err := linkparser.ParseHTMLLinks(httpBody)
			if err != nil {
				log.Fatal(err)
			}

			for _, l := range linksOnPage {
				u, err := url.Parse(l.Href)
				// Skip any invalid urls
				if err != nil {
					continue
				}

				currentPage := u.String()
				// Check if link has been visited
				if visitedLinks[currentPage] {
					continue
				}
				visitedLinks[currentPage] = true

				// Do not include empty links or anchor tag links
				if currentPage == "" || isAnchorLink(u) {
					continue
				}

				var sitemapLink string
				if linkInDomain(u, parsedInputUrl) {
					sitemapLink = currentPage
				} else if isRelativeLink(u) {
					sitemapLink = relToAbsURL(currentPage, currentUrl.String())
				}

				if sitemapLink != "" {
					nextQueue = append(nextQueue, sitemapLink)
					sitemapLinks = append(sitemapLinks, sitemapLink)
					// fmt.Println("\t-x- Adding Link:", sitemapLink)
				}
			}
		}
	}

	return sitemapLinks
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

func printHelpMessage() {
	fmt.Println("This program builds a sitemap for the specified URL's domain")
	fmt.Println("using the standard sitemap protocol.")
	fmt.Println()
	fmt.Println("Usage: ./sitemap (-file=<filename>) <url>")
}
