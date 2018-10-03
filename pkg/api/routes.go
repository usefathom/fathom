package api

import (
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
)

func (api *API) Routes() *mux.Router {
	// register routes
	r := mux.NewRouter()
	r.Handle("/collect", NewCollector(api.database)).Methods(http.MethodGet)

	r.Handle("/api/sites", HandlerFunc(api.GetSitesHandler)).Methods(http.MethodGet)
	r.Handle("/api/sites", HandlerFunc(api.SaveSiteHandler)).Methods(http.MethodPost)
	r.Handle("/api/sites/{id:[0-9]+}", HandlerFunc(api.SaveSiteHandler)).Methods(http.MethodPost)
	r.Handle("/api/sites/{id:[0-9]+}", HandlerFunc(api.DeleteSiteHandler)).Methods(http.MethodDelete)

	r.Handle("/api/session", HandlerFunc(api.GetSession)).Methods(http.MethodGet)
	r.Handle("/api/session", HandlerFunc(api.CreateSession)).Methods(http.MethodPost)
	r.Handle("/api/session", HandlerFunc(api.DeleteSession)).Methods(http.MethodDelete)

	r.Handle("/api/stats/site", api.Authorize(HandlerFunc(api.GetSiteStatsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/groupby/day", api.Authorize(HandlerFunc(api.GetSiteStatsPerDayHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/pageviews", api.Authorize(HandlerFunc(api.GetSiteStatsPageviewsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/visitors", api.Authorize(HandlerFunc(api.GetSiteStatsVisitorsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/duration", api.Authorize(HandlerFunc(api.GetSiteStatsDurationHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/bounces", api.Authorize(HandlerFunc(api.GetSiteStatsBouncesHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/site/realtime", api.Authorize(HandlerFunc(api.GetSiteStatsRealtimeHandler))).Methods(http.MethodGet)

	r.Handle("/api/stats/pages", api.Authorize(HandlerFunc(api.GetPageStatsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/pages/pageviews", api.Authorize(HandlerFunc(api.GetPageStatsPageviewsHandler))).Methods(http.MethodGet)

	r.Handle("/api/stats/referrers", api.Authorize(HandlerFunc(api.GetReferrerStatsHandler))).Methods(http.MethodGet)
	r.Handle("/api/stats/referrers/pageviews", api.Authorize(HandlerFunc(api.GetReferrerStatsPageviewsHandler))).Methods(http.MethodGet)

	r.Handle("/health", HandlerFunc(api.Health)).Methods(http.MethodGet)

	// static assets & 404 handler
	box := packr.NewBox("./../../assets/build")
	r.Path("/tracker.js").Handler(serveTrackerFile(&box))
	r.Path("/").Handler(serveFileHandler(&box, "index.html"))
	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(box)))
	r.NotFoundHandler = NotFoundHandler(&box)

	return r
}

func serveTrackerFile(box *packr.Box) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Tk", "N")
		next := serveFile(box, "js/tracker.js")
		next.ServeHTTP(w, r)
	})
}
