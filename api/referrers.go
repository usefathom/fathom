package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/referrers
var GetReferrersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Referrers(before, after, getRequestedLimit(r))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
