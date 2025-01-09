package main

import (
	"fmt"
	"html"
	"net/http"

	frouter "github.com/xvargr/very-fast-website/internal/file-router"
)

func handleWithFileRouter(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "%s", html.EscapeString(request.RequestURI))
}

func main() {
	mux := http.NewServeMux()

	frouter.Route(mux)

	http.ListenAndServe(":8080", mux)
}