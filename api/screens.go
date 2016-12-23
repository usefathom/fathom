package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dannyvankooten/ana/count"
)

// URL: /api/screen-resolutions
var GetScreenResolutionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)

	// get total
	total := count.Visitors(before, after)

	fmt.Printf("Total # of visitors: %d\n", total)

	// get rows
	results := count.Screens(before, after, getRequestedLimit(r), total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
