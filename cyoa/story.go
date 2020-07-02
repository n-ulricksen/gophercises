package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	// Template must be acceptable at runtime
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var defaultHandlerTemplate = `
<!DOCTYPE HTML>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title></title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
    {{end}}
    <ul>
		{{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
		{{end}}
    </ul>
</body>
</html>
`

func NewHandler(story Story) http.Handler {
	return handler{story}
}

type handler struct {
	story Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		// Start story from beginning
		path = "/intro"
	}
	// trim preceding slash ('/')
	path = path[1:]

	if chapter, ok := h.story[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func JsonStory(jsonFile io.Reader) (Story, error) {
	jsonDecoder := json.NewDecoder(jsonFile)
	var story Story
	if err := jsonDecoder.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []struct {
		Text    string `json:"text"`
		Chapter string `json:"arc"`
	} `json:"options"`
}

func main() {
	fmt.Println("vim-go")

	// TODO:
	// Parse JSON
	// Build http.Handler
	// Parse paths
	// Style HTML
	// Custom Templates
	// Functional Options
	// Custom Paths
}
