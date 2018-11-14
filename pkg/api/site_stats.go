package api

import (
	"net/http"
)

// URL: /api/sites/{id:[0-9]+}/stats/site/agg
func (api *API) GetAggregatedSiteStatsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetAggregatedSiteStats(params.SiteID, params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

// URL: /api/sites/{id:[0-9]+}/stats/site/realtime
func (api *API) GetSiteStatsRealtimeHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.GetRealtimeVisitorCount(params.SiteID)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}

// URL: /api/sites/{id:[0-9]+}/stats/site
func (api *API) GetSiteStatsHandler(w http.ResponseWriter, r *http.Request) error {
	params := GetRequestParams(r)
	result, err := api.database.SelectSiteStats(params.SiteID, params.StartDate, params.EndDate)
	if err != nil {
		return err
	}
	return respond(w, http.StatusOK, envelope{Data: result})
}
