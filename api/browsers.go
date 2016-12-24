package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/browsers
var GetBrowsersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Browsers(before, after, getRequestedLimit(r))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
