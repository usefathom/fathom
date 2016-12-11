package api

import (
	"encoding/json"
	"github.com/dannyvankooten/ana/count"
	"net/http"
)

// URL: /api/screen-resolutions
var GetScreenResolutionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)

	// get total
	total := count.Visitors(before, after)

	// get rows
	results := count.Custom(`
    SELECT
    v.screen_resolution,
    COUNT(DISTINCT(pv.visitor_id)) AS count
    FROM pageviews pv
    LEFT JOIN visitors v ON v.id = pv.visitor_id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    GROUP BY v.screen_resolution
    ORDER BY count DESC
    LIMIT ?`, before, after, getRequestedLimit(r), total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
