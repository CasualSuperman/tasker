package main

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/sha256"
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

type handlerFunc func(http.ResponseWriter, *http.Request, db.Database)

var handlers = map[string]handlerFunc{
	"user/login":    userLogin,
	"user/activate": userActivate,
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
	res.Header().Add("Content-Type", "text/plain")
	email := req.Form.Get("email")

	if email != "" {
		users := sess.ExistentCollection("Users")
		user, _ := users.Find(db.Cond{"email": email})

		if user != nil && user.GetBool("activated") {
			hashedPass := user.GetString("password")
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass),
				[]byte(req.Form.Get("password")))

			if err == nil {
				session, _ := store.Get(req, "calendar")
				session.Values["logged-in"] = true
				session.Values["uid"] = user.GetInt("uid")
				session.Save(req, res)
				res.Write([]byte("{\"successful\": true}"))
				return
			}
		}
	}
	res.Write([]byte("{\"successful\": false}"))
}

func userActivate(res http.ResponseWriter, req *http.Request, sess db.Database) {
	res.Header().Add("Content-Type", "text/plain")
	validation := req.Form.Get("validation")
	email := req.Form.Get("email")

	if email != "" {
		users := sess.ExistentCollection("Users")
		user, _ := users.Find(db.Cond{"email": email})

		if user != nil && user.GetBool("activated") {
			uid := user.GetInt("uid")
			hashedPass := user.GetString("password")
			key := makeKey(email, hashedPass, uid)
			hash := hmac.New(sha256.New, encKey)
			result := hash.Sum(key)
			if bytes.Equal(result, []byte(validation)) {
				users.Update(
					db.Set{"activated": true},
					db.Cond{"email": email},
				)
				res.Write([]byte("{\"activated\": true}"))
				return
			}
		}
	}
	res.Write([]byte("{\"activated\": false}"))
}

func makeKey(email, hashedPass string, uid int64) []byte {
	bEmail := []byte(email)
	bHashedPass := []byte(hashedPass)
	key := make([]byte, 0)

	key = append(key, bEmail[:len(bEmail)/2]...)
	key = append(key, bHashedPass[len(bHashedPass)/3:]...)
	key = append(key, byte(uid))
	key = append(key, bHashedPass[:len(bHashedPass)/3]...)
	key = append(key, bEmail[len(bEmail)/2:]...)
	return key
}
