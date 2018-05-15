package api

import (
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) Routes() *mux.Router {
	// register routes
	r := mux.NewRouter()
	r.Handle("/collect", api.NewCollectHandler()).Methods(http.MethodGet)
	r.Handle("/api/session", HandlerFunc(api.LoginHandler)).Methods(http.MethodPost)
	r.Handle("/api/session", HandlerFunc(api.LogoutHandler)).Methods(http.MethodDelete)

	r.Handle("/api/stats/site/pageviews", api.Authorize(HandlerFunc(api.GetSiteStatsPageviewsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/visitors", api.Authorize(HandlerFunc(api.GetSiteStatsVisitorsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/duration", api.Authorize(HandlerFunc(api.GetSiteStatsDurationHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/bounces", api.Authorize(HandlerFunc(api.GetSiteStatsBouncesHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/realtime", api.Authorize(HandlerFunc(api.GetSiteStatsRealtimeHandler))).Methods(http.MethodGet)

	r.Handle("/api/stats/pages", api.Authorize(HandlerFunc(api.GetPageStatsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/pages/pageviews", api.Authorize(HandlerFunc(api.GetPageStatsPageviewsHandler))).Methods(http.MethodGet)

	r.Handle("/api/stats/referrers", api.Authorize(HandlerFunc(api.GetReferrerStatsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/referrers/pageviews", api.Authorize(HandlerFunc(api.GetReferrerStatsPageviewsHandler))).Methods(http.MethodGet)

	// static assets & 404 handler
	box := packr.NewBox("./../../build")
	r.Path("/tracker.js").Handler(serveFileHandler(&box, "js/tracker.js"))
	r.Path("/").Handler(serveFileHandler(&box, "index.html"))
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(box)))
	r.NotFoundHandler = NotFoundHandler(&box)

	return r
}
