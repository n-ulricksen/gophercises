package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ulricksennick/gophercises/linkparser"
)

func main() {
	html := `
		<a href="/dog">
  			<span>Something in a span</span>
  			Text not in a span
  			<b>Bold text!</b>
		</a>
	`
	r := strings.NewReader(html)

	htmlURL := "https://ulricksen.codes/"
	// 	htmlURL := "https://en.wikipedia.org/wiki/Democratic_Republic_of_the_Congo"
	response, err := http.Get(htmlURL)
	if err != nil {
		log.Fatal(err)
	}
	htmlReader := response.Body
	defer htmlReader.Close()

	htmlLinks := linkparser.ParseHTMLLinks(htmlReader)
	for _, link := range htmlLinks {
		fmt.Println(link)
	}

	htmlLinks2 := linkparser.ParseHTMLLinks(r)
	for _, link := range htmlLinks2 {
		fmt.Println(link)
	}
}
