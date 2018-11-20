package api

import (
	"net/http"
)

func (api *API) GetAggregatedReferrerStatsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.SelectAggregatedReferrerStats(params.SiteID, params.StartDate, params.EndDate, params.Offset, params.Limit)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

func (api *API) GetAggregatedReferrerStatsPageviewsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedReferrerStatsPageviews(params.SiteID, params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}
