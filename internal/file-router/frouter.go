package frouter

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"time"
)

func Route(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format(time.DateTime), r.Method, html.EscapeString(r.RequestURI))

		content, err := os.ReadFile(`web` + strings.ReplaceAll(html.EscapeString(r.RequestURI), "/", "\\"))
		if err != nil {
			fmt.Println("ERR 404", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "%s", content)
	})
}
