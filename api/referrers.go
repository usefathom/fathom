package api

import (
	"encoding/json"
	"github.com/dannyvankooten/ana/count"
	"net/http"
)

// URL: /api/referrers
var GetReferrersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)

	// get total
	total := count.Visitors(before, after)

	// get rows
	results := count.Custom(`
    SELECT
    pv.referrer_url,
    COUNT(DISTINCT(pv.visitor_id)) AS count
    FROM pageviews pv
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    AND pv.referrer_url IS NOT NULL
    AND pv.referrer_url != ""
    GROUP BY pv.referrer_url
    ORDER BY count DESC
    LIMIT ?`, before, after, getRequestedLimit(r), total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
