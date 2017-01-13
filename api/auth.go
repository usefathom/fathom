package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dannyvankooten/ana/db"
	"github.com/dannyvankooten/ana/models"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var store = sessions.NewFilesystemStore("./storage/sessions/", []byte(os.Getenv("ANA_SECRET_KEY")))

// URL: POST /api/session
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// check login creds
	var l login
	json.NewDecoder(r.Body).Decode(&l)
	var hashedPassword string
	var u models.User
	stmt, _ := db.Conn.Prepare("SELECT id, email, password FROM users WHERE email = ? LIMIT 1")
	err := stmt.QueryRow(l.Email).Scan(&u.ID, &u.Email, &hashedPassword)

	// compare pwd
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(l.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		respond(w, envelope{Error: "invalid_credentials"})
		return
	}

	// TODO: Replace session filesystem store with DB store.
	session, _ := store.Get(r, "auth")
	session.Values["user_id"] = u.ID
	err = session.Save(r, w)
	checkError(err)

	respond(w, envelope{Data: true})
})

// URL: DELETE /api/session
var LogoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")
	if !session.IsNew {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}

	respond(w, envelope{Data: true})
})

/* middleware */
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "auth")
		userID, ok := session.Values["user_id"]

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// find user
		var u models.User
		stmt, _ := db.Conn.Prepare("SELECT id, email FROM users WHERE id = ? LIMIT 1")
		err := stmt.QueryRow(userID).Scan(&u.ID, &u.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
