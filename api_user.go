package main

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/gosexy/db"
	"net/http"
)

// Login a user by checking their email/password against the email and bcrypt'd
// password in the database. If it is successful, the user gets a session.
func userLogin(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	resp := defaultUserResponse()

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
				resp.Succeed()
			}
		}
	}
	return resp
}

// Activate a user using HMAC authentication. HMAC generated using sha256, a
// secret key, the user's uid, email, and bcrypt'd password. If it is
// successful, we log the user in.
func userActivate(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	resp := defaultUserResponse()

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
					db.Cond{"uid": uid},
				)
				session, _ := store.Get(req, "calendar")
				session.Values["logged-in"] = true
				session.Values["uid"] = uid
				session.Save(req, res)
				resp.Succeed()
			}
		}
	}
	return resp
}

// Using these components of a user in our db, generate a consistent string
// while not revealing important data.
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
