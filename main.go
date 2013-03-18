package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"net/http"
)

func main() {
	// This connects to a local mysql server.
	db, e := sql.Open("mysql", DB_USER+":"+DB_PASS+"@/"+DB_NAME+"?charset=utf8")

	// If it didn't work, quit and tell us about it.
	if e != nil {
		panic(e)
	}
	// Otherwise, close the connection when we exit.
	defer db.Close()

	// Register the api server hooks on our HTTP server
	runApiServer(db)
	// Also register the static content handler
	runStaticContentServer()

	// Actually serve the content
	http.ListenAndServe("0.0.0.0:1802", nil)
}
