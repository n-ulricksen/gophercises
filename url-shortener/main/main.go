package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ulricksennick/Gophercises/url-shortener"
)

func main() {
	mux := defaultMux()

	// Parse flags
	defaultYamlFile := "paths.yaml"
	defaultJsonFile := "paths.json"

	yamlFlag := flag.String("yaml", defaultYamlFile, "yaml path-url file location")
	jsonFlag := flag.String("json", defaultJsonFile, "json path-url file location")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the MapHandler as the fallback
	yamlBytes, err := ioutil.ReadFile(*yamlFlag)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yamlBytes, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the YAMLHandler as the fallback
	jsonBytes, err := ioutil.ReadFile(*jsonFlag)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonBytes, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
