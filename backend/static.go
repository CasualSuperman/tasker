package main

import (
	"net/http"
)

func runStaticContentServer() {
	// Use the built in HTTP fileserver.
	// Hard-code that we want to look in the ./static drectory.
	http.Handle("/", http.FileServer(http.Dir("static")))
}
