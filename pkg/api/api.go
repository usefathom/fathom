package api

import (
	"github.com/gorilla/sessions"
	"github.com/usefathom/fathom/pkg/datastore"
)

type API struct {
	database        datastore.Datastore
	sessions        sessions.Store
	publicDashboard bool
}

// New instantiates a new API object
func New(db datastore.Datastore, secret string, publicDashboard bool) *API {
	return &API{
		database:        db,
		sessions:        sessions.NewCookieStore([]byte(secret)),
		publicDashboard: publicDashboard,
	}
}
