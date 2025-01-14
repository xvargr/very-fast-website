package filerouter

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/xvargr/very-fast-website/internal/logger"
	"github.com/xvargr/very-fast-website/internal/vdoc"
	"golang.org/x/net/html"
)

type Document struct {
	URI        string
	Path       string
	Layouts    []string
	VirtualDoc *vdoc.VirtualDocument
	Extension  string
}

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestUri := html.EscapeString(r.RequestURI)
		logger.Console(logger.SeverityNormal, fmt.Sprintf("%s %s", r.Method, requestUri))

		doc := resolveDocument(requestUri)

		doc.Serve(w, r)
	})
}

func resolveDocument(path string) Document {
	var ext string
	isDirectAccess := strings.Contains(path, ".")
	docPath := "web" + strings.TrimRight(path, "/")
	docPath = strings.Replace(docPath, "..", "", -1)

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
		URI:       path,
		Path:      docPath,
		Layouts:   evaluateLayouts(path),
		Extension: ext,
	}

}

func evaluateLayouts(path string) (layouts []string) {
	fpath := "web"
	for idx, dir := range strings.Split(path, "/") {
		if path == "/" && idx == 1 {
			continue
		}

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

func (doc *Document) Compile() {
	for _, layout := range doc.Layouts {
		// fmt.Println("layout", layout)
		layoutContent, _ := os.ReadFile(layout)
		extr := vdoc.Extract(layoutContent)

		if doc.VirtualDoc == nil {
			doc.VirtualDoc = vdoc.NewVirtualDocument()
		}

		doc.VirtualDoc.Merge(extr)
	}

	mainContent, err := os.ReadFile(doc.Path)
	if err != nil {
		logger.Console(logger.SeverityError, fmt.Sprintf("failed to read file %s", doc.Path))
		mainContent, _ = os.ReadFile("web/404.html")
	}

	// fmt.Println("main", doc.Path)
	extr := vdoc.Extract(mainContent)
	doc.VirtualDoc.Merge(extr)
}

func (doc *Document) IsDirectAccess() bool {
	return strings.Contains(doc.URI, ".")
}

func (doc *Document) AddTypeHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", mime.TypeByExtension("."+doc.Extension))
}

func (doc *Document) Serve(w http.ResponseWriter, r *http.Request) {
	if doc.IsDirectAccess() {
		http.ServeFile(w, r, doc.Path)
		return
	}

	doc.Compile()
	doc.AddTypeHeader(w)
	// doc.VirtualDoc.Render(&w)
	// html.Render(w, doc.VirtualDoc.RenderHtml())
	w.Write([]byte(doc.VirtualDoc.RenderHtml()))
}
