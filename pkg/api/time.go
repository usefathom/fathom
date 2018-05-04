package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
)

// TODO: Come up with more consistent URL names.
// URL: /api/time-on-site/count
var GetTimeOnSiteCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	result, err := datastore.AvgTimeOnSite(before, after)

	if err != nil {
		return err
	}

	return respond(w, envelope{Data: result})
})
