package urlshort

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if redirectUrl, ok := pathsToUrls[requestPath]; ok {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		}

		fallback.ServeHTTP(w, r)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	yamlMap := buildMapFromYaml(parsedYAML)

	return MapHandler(yamlMap, fallback), nil
}

type yamlUrlPath struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYAML(yml []byte) ([]yamlUrlPath, error) {
	urlPaths := []yamlUrlPath{}

	err := yaml.Unmarshal(yml, &urlPaths)
	if err != nil {
		return nil, err
	}

	return urlPaths, nil
}

func buildMapFromYaml(ymlPaths []yamlUrlPath) map[string]string {
	yamlMap := map[string]string{}
	for _, urlPath := range ymlPaths {
		path := urlPath.Path
		url := urlPath.URL
		yamlMap[path] = url
	}

	return yamlMap
}
