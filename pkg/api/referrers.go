package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
)

// URL: /api/referrers
var GetReferrersHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results := count.Referrers(before, after, getRequestedLimit(r))
	return respond(w, envelope{Data: results})
})
