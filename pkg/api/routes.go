package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Routes(webroot string) *mux.Router {
	// register routes
	r := mux.NewRouter()
	r.Handle("/collect", NewCollectHandler()).Methods(http.MethodGet)
	r.Handle("/api/session", LoginHandler).Methods(http.MethodPost)
	r.Handle("/api/session", LogoutHandler).Methods(http.MethodDelete)

	r.Handle("/api/stats/page", Authorize(GetPageStatsHandler)).Methods(http.MethodGet)

	r.Path("/tracker.js").Handler(http.FileServer(http.Dir(webroot + "/js/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(webroot)))
	return r
}
