package main

import (
	"github.com/gosexy/db"
	_ "github.com/gosexy/db/sqlite"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// This connects to a local mysql server.
	sess, err := db.Open("sqlite", dbSettings)

	if err != nil {
		// If it didn't work, quit and tell us about it.
		panic(err)
	}

	// Otherwise, close the connection when we exit.
	defer sess.Close()

	// Register the api server hooks on our HTTP server
	runApiServer(sess)
	// Also register the static content handler
	runStaticContentServer()

	// Actually serve the content
	http.ListenAndServe("0.0.0.0:1802", nil)
}
