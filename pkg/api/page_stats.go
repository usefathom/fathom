package api

import (
	"net/http"
)

// URL: /api/sites/{id:[0-9]+}/stats/pages/agg
func (api *API) GetAggregatedPageStatsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.SelectAggregatedPageStats(params.SiteID, params.StartDate, params.EndDate, params.Limit)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

func (api *API) GetAggregatedPageStatsPageviewsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedPageStatsPageviews(params.SiteID, params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}
