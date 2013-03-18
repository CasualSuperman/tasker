package main

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"net/http"
)

// Set up secure cookie storage. This byte string is a secret key used to
// authenticate a cookie.
var store = sessions.NewCookieStore([]byte{0x65, 0x23, 0x51, 0x53, 0x6e, 0x4b,
	0x65, 0x34, 0x33, 0x39, 0x55, 0xff,
	0x3e, 0xe4, 0x77, 0x20, 0x00, 0xe1})

// This runs our API server. We take a database connection so we could
// theoretically run multiple API servers at different locations with different
// database conenctions.
func runApiServer(db *sql.DB) {
	// This just uses an anonymous function for now to show that it works.
	http.HandleFunc("/api/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Welcome to the API server. You requested " +
			req.URL.Path))
	})
}
