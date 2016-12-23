package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/browsers
var GetBrowsersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)

	// get total
	total := count.Visitors(before, after)

	// get rows
	results := count.Browsers(before, after, getRequestedLimit(r), total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
