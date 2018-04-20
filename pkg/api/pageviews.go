package api

import (
	"net/http"

	"github.com/dannyvankooten/ana/pkg/count"
	"github.com/dannyvankooten/ana/pkg/datastore"
)

type pageviews struct {
	Hostname    string
	Path        string
	Count       int
	CountUnique int
}

// URL: /api/pageviews
var GetPageviewsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)

	stmt, err := datastore.DB.Prepare(`
		SELECT
			p.hostname,
			p.path,
			SUM(t.count) AS count,
			SUM(t.count_unique) AS count_unique
		FROM total_pageviews t
		LEFT JOIN pages p ON p.id = t.page_id
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?
		GROUP BY p.path, p.hostname
		ORDER BY count DESC
		LIMIT ?`)

	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(before, after, defaultLimit)
	checkError(err)
	defer rows.Close()

	results := make([]pageviews, 0)
	for rows.Next() {
		var p pageviews
		err = rows.Scan(&p.Hostname, &p.Path, &p.Count, &p.CountUnique)
		checkError(err)
		results = append(results, p)
	}

	err = rows.Err()
	checkError(err)

	respond(w, envelope{Data: results})
})

// URL: /api/pageviews/count
var GetPageviewsCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	result := count.Pageviews(before, after)
	respond(w, envelope{Data: result})
})

// URL: /api/pageviews/group/day
var GetPageviewsPeriodCountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.PageviewsPerDay(before, after)
	respond(w, envelope{Data: results})
})
