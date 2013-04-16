package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/gosexy/db"
	"net/http"
	"net/smtp"
	"strings"
)

// Give info about the logged in user if they are logged in
func userInfo(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	session, _ := store.Get(req, "calendar")

	if val, ok := session.Values["logged-in"]; ok && val.(bool) {
		uid := int(session.Values["uid"].(int64))
		users := sess.ExistentCollection("Users")
		user, err := users.Find(db.Cond{"uid": uid})

		if user != nil && err == nil {
			res.Header().Add("content-type", "text/plain")
			res.WriteHeader(http.StatusOK)
			data, _ := json.Marshal(user)
			res.Write(data)
			return nil
		}
	}
	return apiUserResponse{false, "", http.StatusOK}
}

// Register a user by checking if their email is in the database. If it isn't,
// the account is created and the activation email is sent.
func userRegister(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	resp := defaultUserResponse()

	email := req.FormValue("email")
	password := req.FormValue("password")

	// Only proceed if we were given both.
	if email != "" && password != "" {
		users := sess.ExistentCollection("Users")
		num, err := users.Count(db.Cond{"email": email})

		// Make sure there are no other users with that email
		if num == 0 && err == nil {
			// Encrypt their password
			hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password),
				bcrypt.DefaultCost)

			// Add them to the database
			_, err = users.Append(db.Item{
				"email":     email,
				"password":  hashedPass,
				"activated": false,
			})

			if err != nil {
				resp.Fail(err)
			} else {
				// Generate the activation email
				url := baseURL + "api/user/activate?email=" + email +
					"&validation=" + makeValidationCode(email, hashedPass)

				// Send the email
				sendEmail(email, "Tasker Account Activation",
					"Welcome to tasker!\n"+
						"\n"+
						"Please activate your email by clicking this link: "+url+"\n"+
						"\n"+
						"Thanks,\n"+
						"The Tasker Team",
				)

				if err == nil {
					resp.Succeed()
					resp.code = http.StatusCreated
				} else {
					// Remove the user if we couldn't send the email.
					users.Remove(
						db.Cond{"email": email},
					)
					resp.Fail(err)
				}
			}
		} else if err == nil {
			resp.Err = "User already exists."
			resp.code = http.StatusPreconditionFailed
		} else {
			resp.Fail(err)
		}
	} else {
		resp.Err = "Registration requires both a username and a password."
		resp.code = http.StatusBadRequest
	}
	return resp
}

// Login a user by checking their email/password against the email and bcrypt'd
// password in the database. If it is successful, the user gets a session.
func userLogin(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	resp := defaultUserResponse()

	email := req.FormValue("email")

	if email != "" {
		users := sess.ExistentCollection("Users")
		user, _ := users.Find(db.Cond{"email": email})

		if user != nil && user.GetBool("activated") {
			hashedPass := user.GetString("password")
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass),
				[]byte(req.FormValue("password")))

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

// Checks if the request has a session associated with it, and deletes login
// info from the session if there is. Always returns success.
func userLogout(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	resp := defaultUserResponse()
	session, isNew := store.Get(req, "calendar")
	if isNew == nil {
		delete(session.Values, "logged-in")
		delete(session.Values, "uid")
		session.Save(req, res)
	}
	resp.Succeed()
	return resp
}

// Activate a user using HMAC authentication. HMAC generated using sha256, a
// secret key, the user's uid, email, and bcrypt'd password. If it is
// successful, we log the user in.
func userActivate(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	resp := defaultUserResponse()

	validation := req.FormValue("validation")
	email := req.FormValue("email")

	if email != "" {
		users := sess.ExistentCollection("Users")
		user, _ := users.Find(db.Cond{"email": email})

		if user != nil && !user.GetBool("activated") {
			uid := user.GetInt("uid")
			hashedPass := user.GetString("password")
			result := makeValidationCode(email, []byte(hashedPass))
			if result == validation {
				users.Update(
					db.Set{"activated": true},
					db.Cond{"uid": uid},
				)
				session, _ := store.Get(req, "calendar")
				session.Values["logged-in"] = true
				session.Values["uid"] = uid
				session.Save(req, res)

				sendEmail(user.GetString("email"), "Tasker Account Activation",
					"Your account has been activated!\n"+
						"Go to "+baseURL+" to start using Tasker.\n"+
						"\n"+
						"Thanks,\n"+
						"The Tasker Team",
				)

				resp.Succeed()
			}
		}
	}
	return resp
}

// Using these components of a user in our db, generate a consistent string
// while not revealing important data.
func makeKey(email, hashedPass string) []byte {
	bEmail := []byte(email)
	bHashedPass := []byte(hashedPass)
	key := make([]byte, 0)

	key = append(key, bEmail[len(bEmail)/2:]...)
	key = append(key, bHashedPass[len(bHashedPass)/3:]...)
	key = append(key, bEmail[:len(bEmail)/2]...)
	key = append(key, bHashedPass[:2*len(bHashedPass)/3]...)
	return key
}

func makeValidationCode(email string, hashedPass []byte) string {
	key := makeKey(email, string(hashedPass))
	hash := hmac.New(sha256.New, encKey)
	hash.Write(key)
	code := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return strings.TrimRight(code, "=")
}

func sendEmail(to string, subject, body string) error {
	err := smtp.SendMail(
		"smtp.gmail.com:25",
		emailAuth,
		"tasker@casualsuperman.com",
		[]string{to},
		[]byte("Subject: "+subject+"\n"+body),
	)
	return err
}
