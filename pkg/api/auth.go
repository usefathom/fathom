package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/usefathom/fathom/pkg/datastore"
	"golang.org/x/crypto/bcrypt"
)

type key int

const (
	userKey key = 0
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("ANA_SECRET_KEY")))

// URL: POST /api/session
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// check login creds
	var l login
	json.NewDecoder(r.Body).Decode(&l)

	u, err := datastore.GetUserByEmail(l.Email)

	// compare pwd
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(l.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		respond(w, envelope{Error: "invalid_credentials"})
		return
	}

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
		u, err := datastore.GetUser(userID.(int64))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
