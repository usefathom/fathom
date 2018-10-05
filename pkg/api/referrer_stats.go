package api

import (
	"net/http"
)

// URL: /api/stats/referrer
func (api *API) GetReferrerStatsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedReferrerStats(params.SiteID, params.StartDate, params.EndDate, params.Limit)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

// URL: /api/stats/referrer/pageviews
func (api *API) GetReferrerStatsPageviewsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedReferrerStatsPageviews(params.SiteID, params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}
