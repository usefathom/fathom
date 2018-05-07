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

	r.Handle("/api/stats/site/pageviews", Authorize(GetSiteStatsPageviewsHandler)).Methods(http.MethodGet)
	r.Handle("/api/stats/site/visitors", Authorize(GetSiteStatsVisitorsHandler)).Methods(http.MethodGet)
	r.Handle("/api/stats/site/duration", Authorize(GetSiteStatsDurationHandler)).Methods(http.MethodGet)
	r.Handle("/api/stats/site/bounces", Authorize(GetSiteStatsBouncesHandler)).Methods(http.MethodGet)
	r.Handle("/api/stats/site/realtime", Authorize(GetSiteStatsRealtimeHandler)).Methods(http.MethodGet)

	r.Path("/tracker.js").Handler(http.FileServer(http.Dir(webroot + "/js/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(webroot)))
	return r
}
