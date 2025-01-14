package main

import (
	"fmt"
	"html"
	"net/http"

	filerouter "github.com/xvargr/very-fast-website/internal/file-router"
)

func handleWithFileRouter(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "%s", html.EscapeString(request.RequestURI))
}

func main() {
	mux := http.NewServeMux()

	filerouter.Route(mux)

	http.ListenAndServe(":85", mux)
}
