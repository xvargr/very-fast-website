package filerouter

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/xvargr/very-fast-website/internal/logger"
	"github.com/xvargr/very-fast-website/internal/vdoc"
)

type Document struct {
	URI         string
	docPath     string
	Layouts     []string
	VirtualDoc  *vdoc.VirtualDocument
	Extension   string
	isHxRequest bool
}

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path
		logger.Console(logger.SeverityNormal, fmt.Sprintf("%s %s", r.Method, uri))

		doc := resolveDocument(r)

		doc.Serve(w, r)
	})
}

func resolveDocument(r *http.Request) Document {
	var ext string
	path := r.URL.Path
	isDirectAccess := strings.Contains(path, ".")
	isHxRequest := r.Header.Get("Hx-Request") == "true"

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
		URI:         path,
		docPath:     docPath,
		VirtualDoc:  vdoc.NewVirtualDocument(),
		Layouts:     evaluateLayouts(r),
		Extension:   ext,
		isHxRequest: isHxRequest,
	}

}

func evaluateLayouts(r *http.Request) (layouts []string) {
	path := strings.TrimRight(r.URL.Path, "/")
	pathSplit := strings.Split(path, "/")
	fpath := "web"
	isHxRequest := r.Header.Get("Hx-Request") == "true"

	if isHxRequest && path != "/" {
		pathSplit = pathSplit[:len(pathSplit)-1]
	}

	for _, dir := range pathSplit {
		fpath = fpath + dir + "/"
		dirContent, _ := os.ReadDir(fpath)
		for _, file := range dirContent {
			if file.Name() == "_layout.html" {
				layouts = append(layouts, fpath+file.Name())
			}
		}
	}

	if isHxRequest && len(layouts) > 1 {
		layouts = layouts[:len(layouts)-1]
	}

	return layouts
}

func evaluateDocPath(r *http.Request) string {

	return ""
}

func (doc *Document) Compile() {
	for _, layout := range doc.Layouts {
		// logger.Console(logger.SeverityDebug, fmt.Sprintf("layout %s", layout))
		layoutContent, _ := os.ReadFile(layout)
		extr := vdoc.Extract(layoutContent)
		doc.VirtualDoc.Merge(extr)
	}

	logger.Console(logger.SeverityDebug, fmt.Sprintf("docPath %s", doc.docPath))
	// contentDirExists, err := os.Stat(doc.docPath)
	// if os.IsNotExist(err) {

	// }

	mainContent, err := os.ReadFile(doc.docPath)
	if err != nil {
		// logger.Console(logger.SeverityError, fmt.Sprintf("failed to read file %s", doc.docPath))
		mainContent, _ = os.ReadFile("web/404.html")
	}

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
		http.ServeFile(w, r, doc.docPath)
		return
	}

	doc.Compile()
	doc.AddTypeHeader(w)
	w.Write([]byte(doc.VirtualDoc.RenderHtml()))
}
