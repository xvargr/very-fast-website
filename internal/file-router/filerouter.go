package filerouter

import (
	"fmt"
	"html"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"
)

type Document struct {
	Path            string
	Layouts         []string
	CompiledContent string
	Extension       string
}

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestUri := html.EscapeString(r.RequestURI)
		fmt.Println(time.Now().Format(time.DateTime), r.Method, requestUri)

		doc := resolveDocument(requestUri)
		fmt.Println(doc)
		// doc.Compile()

		// AddTypeHeader(w, doc)
		// http.ServeFile(w, r, doc.Path)
		doc.Serve(w, r)
	})
}

func resolveDocument(path string) Document {
	var ext string
	isDirectAccess := strings.Contains(path, ".")
	docPath := "web" + strings.TrimRight(path, "/")
	docPath = strings.Replace(docPath, "..", "", -1)
	fmt.Println("DocPath", docPath)

	filename := "index"
	split := strings.Split(docPath, "/")
	fmt.Println("Split", split, len(split), split[len(split)-1])
	if split[len(split)-1] != "" {
		filename = split[len(split)-1]
	}

	if !isDirectAccess {
		docPath = docPath + "/" + filename + ".html"
		ext = "html"
	} else {
		ext = strings.Split(filename, ".")[1]
	}

	fmt.Println("filename", filename, "ext", ext)

	return Document{
		Path:      docPath,
		Layouts:   evaluateLayouts(path),
		Extension: ext,
	}
}

func evaluateLayouts(path string) (layouts []string) {
	fpath := "web"
	for _, dir := range strings.Split(path, "/") {
		fpath = fpath + dir + "/"

		dirContent, _ := os.ReadDir(fpath)
		for _, file := range dirContent {
			if file.Name() == "_layout.html" {
				layouts = append(layouts, fpath+file.Name())
			}
		}
	}

	return layouts
}

// func (doc Document) String() string {
// 	return fmt.Sprintf("Document: Path: %s, Layouts: %v, Extension: %s,", doc.Path, doc.Layouts, doc.Extension)
// }

func (doc *Document) Compile() {
	fmt.Println("Compiling", doc)
	for _, layout := range doc.Layouts {
		fmt.Println("Layout", layout)
		rawContent, _ := os.ReadFile(layout)
		// content := string(rawContent)

		if doc.CompiledContent == "" {
			doc.CompiledContent = string(rawContent)
		} else {
			doc.CompiledContent = strings.Replace(doc.CompiledContent, "{{content}}", strings.TrimSpace(string(rawContent)), -1)
		}
	}

	fmt.Println("Compiled", doc.CompiledContent)
}

func (doc *Document) IsDirectAccess() bool {
	return strings.Contains(doc.Path, ".")
}

func (doc *Document) AddTypeHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", mime.TypeByExtension("."+doc.Extension))
}

func (doc *Document) Serve(w http.ResponseWriter, r *http.Request) {
	if doc.IsDirectAccess() {
		http.ServeFile(w, r, doc.Path)
		return
	}

	doc.AddTypeHeader(w)
	w.Write([]byte(doc.CompiledContent))
}

// func (doc *Document) addContent() {
// 	content, _ := os.ReadFile(doc.Path)
// }

// func AddTypeHeader(w http.ResponseWriter, doc Document) string {
// 	return mime.TypeByExtension("." + doc.Extension)
// }
