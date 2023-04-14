package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Story map[string]Chapter

type Chapter struct {
	Title     string   `json:"title"`
	Paragraph []string `json:"story"`
	Options   []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

var temp = ` <!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
	
      <h1>{{.Title}}</h1>
	  <h4>{{.Paragraph}}</h4>
	  <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
		  {{end}}
		  </ul>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
  </body>
</html>`

func readChapter(title string) Chapter {
	var cap Chapter
	r, err := os.ReadFile("./gopher.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var s Story
	json.Unmarshal(r, &s)
	if err != nil {
		fmt.Println("error:", err)
	}
	_, ok := s[title]
	if ok {
		cap = s[title]
	}
	return cap
}

func GetHendle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.New("template").Parse(temp))
		path := request.URL.Path
		pathNew := strings.Replace(path, "/", "", -1)
		chapter := readChapter(pathNew)
		tmpl.Execute(writer, chapter)
	}
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", GetHendle()))
	fmt.Println("Starting the server on :8080")
}
