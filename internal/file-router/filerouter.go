package filerouter

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xvargr/very-fast-website/internal/logger"
	"github.com/xvargr/very-fast-website/internal/vdoc"
)

type Document struct {
	url         string
	Layouts     []string
	VirtualDoc  *vdoc.VirtualDocument
	isHxRequest bool
}

const (
	basePath         = "web/"
	notFoundFilepath = "web/404.html"
)

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Console(logger.SeverityNormal, fmt.Sprintf("%s %s", r.Method, r.URL.Path))

		if isAssetRequest(r) {
			serveAsset(w, r)
		} else {
			doc := makeDocument(r)
			doc.serve(w, r)
		}
	})
}

func isAssetRequest(r *http.Request) bool {
	return strings.Contains(r.URL.Path, ".")
}

func serveAsset(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, basePath+r.URL.Path)
}

func makeDocument(r *http.Request) Document {
	doc := Document{
		url:         r.URL.Path,
		VirtualDoc:  vdoc.NewVirtualDocument(),
		Layouts:     evaluateLayouts(r),
		isHxRequest: r.Header.Get("Hx-Request") == "true",
	}

	return doc
}

func (doc *Document) resolveContentPath() string {
	parts := strings.Split(strings.Trim(doc.url, "/"), "/")

	curPath := basePath
	var filename string
	if doc.url == "/" {
		filename = "index.html"
	} else {
		filename = parts[len(parts)-1] + ".html"
	}

	for idx, part := range parts {
		if part != "" {
			curPath = curPath + part + "/"
		}

		if idx == len(parts)-2 {
			_, err := os.Stat(curPath + filename)
			if err == nil {
				return curPath + filename
			}
		}
		if idx == len(parts)-1 {
			_, err := os.Stat(curPath + filename)
			if err == nil {
				return curPath + filename
			}
		}
	}

	logger.Console(logger.SeverityError, fmt.Sprintf("ERR 404 : %s", doc.url))
	return notFoundFilepath
}

func evaluateLayouts(r *http.Request) (layouts []string) {
	path := strings.TrimRight(r.URL.Path, "/")
	pathSplit := strings.Split(path, "/")
	filepath := basePath
	isHxRequest := r.Header.Get("Hx-Request") == "true"

	if isHxRequest && path != "/" {
		pathSplit = pathSplit[:len(pathSplit)-1]
	}

	for _, dir := range pathSplit {
		filepath = filepath + dir + "/"
		dirContent, _ := os.ReadDir(filepath)
		for _, file := range dirContent {
			if file.Name() == "_layout.html" {
				layouts = append(layouts, filepath+file.Name())
			}
		}
	}

	if isHxRequest && len(layouts) > 1 {
		layouts = layouts[:len(layouts)-1]
	}

	return layouts
}

func (doc *Document) compile() {
	for _, layout := range doc.Layouts {
		layoutContent, _ := os.ReadFile(layout)
		extr := vdoc.Extract(layoutContent)
		doc.VirtualDoc.Merge(extr)
	}

	content, _ := os.ReadFile(doc.resolveContentPath())

	extr := vdoc.Extract(content)
	doc.VirtualDoc.Merge(extr)
}

func (doc *Document) serve(w http.ResponseWriter, r *http.Request) {
	doc.compile()
	w.Write([]byte(doc.VirtualDoc.RenderHtml()))
}
