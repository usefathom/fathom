package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
	"github.com/usefathom/fathom/pkg/datastore"
)

type pageviews struct {
	Hostname    string
	Path        string
	Count       int
	CountUnique int
}

// URL: /api/pageviews
var GetPageviewsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
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
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(before, after, defaultLimit)
	if err != nil {
		return err
	}
	defer rows.Close()

	results := make([]pageviews, 0)
	for rows.Next() {
		var p pageviews
		err = rows.Scan(&p.Hostname, &p.Path, &p.Count, &p.CountUnique)
		if err != nil {
			return err
		}

		results = append(results, p)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: results})
})

// URL: /api/pageviews/count
var GetPageviewsCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	result := count.Pageviews(before, after)
	return respond(w, envelope{Data: result})
})

// URL: /api/pageviews/group/day
var GetPageviewsPeriodCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results := count.PageviewsPerDay(before, after)
	return respond(w, envelope{Data: results})
})
