package main

import (
	"net/http"

	filerouter "github.com/xvargr/very-fast-website/internal/file-router"
)

func main() {
	mux := http.NewServeMux()

	filerouter.Route(mux)

	http.ListenAndServe(":85", mux)
}
