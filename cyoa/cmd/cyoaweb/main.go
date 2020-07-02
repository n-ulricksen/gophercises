package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ulricksennick/Gophercises/cyoa"
)

func main() {
	port := flag.Int("port", 8000, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("using the story in %s\n", *filename)

	jsonFile, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	story, err := cyoa.JsonStory(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	h := cyoa.NewHandler(story, nil)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
