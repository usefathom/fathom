package api

import (
	"encoding/json"
	"net/http"
	"strings"

	gcontext "github.com/gorilla/context"
	"github.com/usefathom/fathom/pkg/datastore"
)

type key int

const (
	userKey key = 0
)

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *login) Sanitize() {
	l.Email = strings.ToLower(strings.TrimSpace(l.Email))
}

// GET /api/session
func (api *API) GetSession(w http.ResponseWriter, r *http.Request) error {
	userCount, err := api.database.CountUsers()
	if err != nil {
		return err
	}

	// if 0 users in database, dashboard is public
	if userCount == 0 {
		return respond(w, envelope{Data: true})
	}

	// if existing session, assume logged-in
	session, _ := api.sessions.Get(r, "auth")
	if !session.IsNew {
		respond(w, envelope{Data: true})
	}

	// otherwise: not logged-in yet
	return respond(w, envelope{Data: false})
}

// URL: POST /api/session
func (api *API) CreateSession(w http.ResponseWriter, r *http.Request) error {
	// check login creds
	var l login
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		return err
	}
	l.Sanitize()

	// find user with given email
	u, err := api.database.GetUserByEmail(l.Email)
	if err != nil && err != datastore.ErrNoResults {
		return err
	}

	// compare pwd
	if err == datastore.ErrNoResults || u.ComparePassword(l.Password) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return respond(w, envelope{Error: "invalid_credentials"})
	}

	// ignore error here as we want a (new) session regardless
	session, _ := api.sessions.Get(r, "auth")
	session.Values["user_id"] = u.ID
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: true})
}

// URL: DELETE /api/session
func (api *API) DeleteSession(w http.ResponseWriter, r *http.Request) error {
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

// Authorize is middleware that aborts the request if unauthorized
func (api *API) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// clear context from request after it is handled
		// see http://www.gorillatoolkit.org/pkg/sessions#overview
		defer gcontext.Clear(r)

		// first count users in datastore
		// if 0, assume dashboard is public
		userCount, err := api.database.CountUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userCount > 0 {
			session, err := api.sessions.Get(r, "auth")
			// an err is returned if cookie has been tampered with, so check that
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, ok := session.Values["user_id"]
			if session.IsNew || !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// validate user ID in session
			if _, err := api.database.GetUser(userID.(int64)); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
