package api

import (
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"net/http"
)

func Routes() *mux.Router {
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

	r.Handle("/api/stats/pages", Authorize(GetPageStatsHandler)).Methods(http.MethodGet)
	r.Handle("/api/stats/pages/pageviews", Authorize(GetPageStatsPageviewsHandler)).Methods(http.MethodGet)

	r.Handle("/api/stats/referrers", Authorize(GetReferrerStatsHandler)).Methods(http.MethodGet)
	r.Handle("/api/stats/referrers/pageviews", Authorize(GetReferrerStatsPageviewsHandler)).Methods(http.MethodGet)

	// static assets & 404 handler
	box := packr.NewBox("./../../build")
	r.Path("/tracker.js").Handler(serveFileHandler(&box, "js/tracker.js"))
	r.Path("/").Handler(serveFileHandler(&box, "index.html"))
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(box)))
	r.NotFoundHandler = NotFoundHandler(&box)

	return r
}
