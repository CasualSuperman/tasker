package main

import (
	"github.com/gorilla/sessions"
	"github.com/gosexy/db"
	"net/http"
)

// Set up secure cookie storage. This byte string is a secret key used to
// authenticate a cookie.
var store = sessions.NewCookieStore(cookieStoreKeys...)

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
