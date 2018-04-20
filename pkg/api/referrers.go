package api

import (
	"net/http"

	"github.com/dannyvankooten/ana/pkg/count"
)

// URL: /api/referrers
var GetReferrersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Referrers(before, after, getRequestedLimit(r))
	respond(w, envelope{Data: results})
})
