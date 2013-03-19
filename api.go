package main

import (
	"github.com/gorilla/sessions"
	"github.com/gosexy/db"
	"net/http"
)

// Set up secure cookie storage. This byte string is a secret key used to
// authenticate a cookie.
var store = sessions.NewCookieStore([]byte{0x65, 0x23, 0x51, 0x53, 0x6e, 0x4b,
	0x65, 0x34, 0x33, 0x39, 0x55, 0xff,
	0x3e, 0xe4, 0x77, 0x20, 0x00, 0xe1})

type handlerFunc func(http.ResponseWriter, *http.Request, db.Database)

var handlers = map[string]handlerFunc{
	"user/login": userLogin,
}

// This runs our API server. We take a database connection so we could
// theoretically run multiple API servers at different locations with different
// database conenctions.
func runApiServer(sess db.Database) {
	// This just uses an anonymous function for now to show that it works.
	http.HandleFunc("/api/", func(res http.ResponseWriter, req *http.Request) {
		handler, ok := handlers[req.URL.Path[len("/api/"):]]
		if ok {
			handler(res, req, sess)
		} else {
			res.WriteHeader(http.StatusNotImplemented)
		}
	})
}

func userLogin(res http.ResponseWriter, req *http.Request, sess db.Database) {
	res.Write([]byte("Trying to login user."))
}
