package api

import (
	"github.com/gorilla/sessions"
	"github.com/usefathom/fathom/pkg/datastore"
)

type API struct {
	database datastore.Datastore
	sessions sessions.Store
}

// New instantiates a new API object
func New(db datastore.Datastore, secret string) *API {
	return &API{
		database: db,
		sessions: sessions.NewCookieStore([]byte(secret)),
	}
}
