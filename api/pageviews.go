package api

import (
	"encoding/json"
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
	"github.com/dannyvankooten/ana/models"
	"net/http"
)

// URL: /api/pageviews
var GetPageviewsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	stmt, err := db.Conn.Prepare(`SELECT
      p.hostname,
      p.path,
      COUNT(*) AS pageviews,
      COUNT(DISTINCT(pv.visitor_id)) AS pageviews_unique
    FROM pageviews pv
    LEFT JOIN pages p ON pv.page_id = p.id
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    GROUP BY p.path, p.hostname
    ORDER BY pageviews DESC
    LIMIT ?`)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, defaultLimit)
	checkError(err)
	defer rows.Close()

	results := make([]models.Pageviews, 0)
	for rows.Next() {
		var p models.Pageviews
		err = rows.Scan(&p.Hostname, &p.Path, &p.Count, &p.CountUnique)
		checkError(err)
		results = append(results, p)
	}

	err = rows.Err()
	checkError(err)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})

// URL: /api/pageviews/count
var GetPageviewsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	result := count.Pageviews(before, after)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
})

// URL: /api/pageviews/group/day
var GetPageviewsPeriodCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.PageviewsPerDay(before, after)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
})
