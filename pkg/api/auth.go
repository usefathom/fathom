package api

import (
	"context"
	"encoding/json"
	"net/http"

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

// URL: POST /api/session
func (api *API) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	// check login creds
	var l login
	json.NewDecoder(r.Body).Decode(&l)

	// find user with given email
	u, err := api.database.GetUserByEmail(l.Email)
	if err != nil && err != datastore.ErrNoResults {
		return err
	}

	// compare pwd
	if err == datastore.ErrNoResults || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return respond(w, envelope{Error: "invalid_credentials"})
	}

	session, _ := api.sessions.Get(r, "auth")
	session.Values["user_id"] = u.ID
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: true})
}

// URL: DELETE /api/session
func (api *API) LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	session, _ := api.sessions.Get(r, "auth")
	if !session.IsNew {
		session.Options.MaxAge = -1
		err := session.Save(r, w)
		if err != nil {
			return err
		}
	}

	return respond(w, envelope{Data: true})
}

/* middleware */
func (api *API) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := api.sessions.Get(r, "auth")
		if err != nil {
			return
		}

		userID, ok := session.Values["user_id"]

		if session.IsNew || !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// find user
		u, err := api.database.GetUser(userID.(int64))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
