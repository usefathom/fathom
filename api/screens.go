package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/screen-resolutions
var GetScreenResolutionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Screens(before, after, getRequestedLimit(r))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
