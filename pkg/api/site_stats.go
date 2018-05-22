package api

import (
	"net/http"
)

// URL: /api/stats/site/pageviews
func (api *API) GetSiteStatsPageviewsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetTotalSiteViews(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
}

// URL: /api/stats/site/visitors
func (api *API) GetSiteStatsVisitorsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetTotalSiteVisitors(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
}

// URL: /api/stats/site/duration
func (api *API) GetSiteStatsDurationHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAverageSiteDuration(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
}

// URL: /api/stats/site/bounces
func (api *API) GetSiteStatsBouncesHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAverageSiteBounceRate(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
}

// URL: /api/stats/site/realtime
func (api *API) GetSiteStatsRealtimeHandler(w http.ResponseWriter, r *http.Request) error {
	result, err := api.database.GetRealtimeVisitorCount()
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
}

// URL: /api/stats/site/groupby/day
func (api *API) GetSiteStatsPerDayHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetSiteStatsPerDay(params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: result})
}
