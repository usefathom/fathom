package api

import (
	"encoding/json"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Languages(before, after, getRequestedLimit(r))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
