package main

import (
	"log"
	"net/http"

	"fmt"

	"gopkg.in/yaml.v2"
)

type t struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path, ok := pathsToUrls[request.URL.Path]
		if ok {
			http.Redirect(writer, request, path, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls := buildMap(yml)
	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(yml []byte) map[string]string {
	var t []t
	yaml.Unmarshal(yml, &t)
	pathsToUrls := make(map[string]string)
	for _, c := range t {
		pathsToUrls[c.Path] = c.Url
	}
	return pathsToUrls
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8081")
	log.Fatal(http.ListenAndServe(":8081", yamlHandler))
}
