package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
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
    <h1>{{.Title}}</h1>_
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
	fmt.Println("handling...")
	err := tpl.Execute(w, h.story["intro"])
	if err != nil {
		log.Fatal(err)
	}
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
