package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)

	// get total
	total := count.Visitors(before, after)

	results := count.Languages(before, after, getRequestedLimit(r), total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
