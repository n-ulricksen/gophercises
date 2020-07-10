package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ulricksennick/gophercises/linkparser"
)

func main() {
	htmlURL := "https://ulricksen.codes/"
	// htmlURL := "https://en.wikipedia.org/wiki/Democratic_Republic_of_the_Congo"
	response, err := http.Get(htmlURL)
	if err != nil {
		log.Fatal(err)
	}
	htmlReader := response.Body
	defer htmlReader.Close()

	htmlLinks, err := linkparser.ParseHTMLLinks(htmlReader)
	if err != nil {
		log.Fatal(err)
	}
	for _, link := range htmlLinks {
		fmt.Println(link)
	}

	testHtml := `
		<a href="/dog">
  			<span>Something in a span</span>
  			Text not in a span
  			<b>Bold text!</b>
		</a>
	`
	r := strings.NewReader(testHtml)

	htmlLinks2, err := linkparser.ParseHTMLLinks(r)
	if err != nil {
		log.Fatal(err)
	}
	for _, link := range htmlLinks2 {
		fmt.Println(link)
	}
}
