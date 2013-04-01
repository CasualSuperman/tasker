package main

import (
	"github.com/gorilla/sessions"
	"github.com/gosexy/db"
	"net/http"
)

var encKey = []byte{0xe3, 0x23, 0x7f, 0x14,
	0x15, 0x16, 0x17, 0xef,
	0x18, 0x19, 0x1a, 0x1b}

// Set up secure cookie storage. This byte string is a secret key used to
// authenticate a cookie.
var store = sessions.NewCookieStore([]byte{0x65, 0x23, 0x51, 0x53, 0x6e, 0x4b,
	0x65, 0x34, 0x33, 0x39, 0x55, 0xff,
	0x3e, 0xe4, 0x77, 0x20, 0x00, 0xe1})

type apiResponse interface {
	Json() []byte
	Code() int
	Type() string
}

type handlerFunc func(http.ResponseWriter, *http.Request, db.Database) apiResponse

// A map of url handlers
var handlers = map[string]handlerFunc{
	"user/login":    userLogin,
	"user/activate": userActivate,
}

// This runs our API server. We take a database connection so we could
// theoretically run multiple API servers at different locations with different
// database conenctions.
func runApiServer(sess db.Database) {
	// If we get a call on /api, check the path and see if we have a handler
	// for it.
	http.HandleFunc("/api/", func(res http.ResponseWriter, req *http.Request) {
		handler, ok := handlers[req.URL.Path[len("/api/"):]]
		if ok {
			resp := handler(res, req, sess)
			res.Header().Add("content-type", resp.Type())
			res.WriteHeader(resp.Code())
			res.Write(resp.Json())
		} else {
			res.WriteHeader(http.StatusNotImplemented)
		}
	})
}
