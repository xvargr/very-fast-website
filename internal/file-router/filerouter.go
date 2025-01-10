package filerouter

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

type Document struct {
	Path      string
	Extension string
}

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format(time.DateTime), r.Method, html.EscapeString(r.RequestURI))

		doc := resolveDocument(r.RequestURI)

		// AddTypeHeader(w, doc)
		http.ServeFile(w, r, doc.Path)
	})
}

func resolveDocument(path string) Document {
	var ext string
	isDirectAccess := strings.Contains(path, ".")
	docPath := "web/" + path

	filename := "index"
	split := strings.Split(path, "/")
	if split[len(split)-1] != "" {
		filename = split[len(split)-1]
	}

	if !isDirectAccess {
		docPath = docPath + "/" + filename + ".html"
		ext = "html"
	} else {
		ext = strings.Split(filename, ".")[1]
	}

	return Document{
		Path:      docPath,
		Extension: ext,
	}

}

// func AddTypeHeader(w http.ResponseWriter, doc Document) string {
// 	return mime.TypeByExtension("." + doc.Extension)
// }
