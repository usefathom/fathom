package api

import (
	"github.com/dannyvankooten/ana/pkg/count"
	"net/http"
)

// URL: /api/browsers
var GetBrowsersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Browsers(before, after, getRequestedLimit(r))
	respond(w, envelope{Data: results})
})
