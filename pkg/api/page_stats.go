package api

import (
	"net/http"
)

// URL: /api/stats/page
func (api *API) GetPageStatsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedPageStats(params.SiteID, params.StartDate, params.EndDate, params.Limit)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

func (api *API) GetPageStatsPageviewsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedPageStatsPageviews(params.SiteID, params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}
