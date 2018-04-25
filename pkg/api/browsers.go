package api

import (
	"github.com/usefathom/fathom/pkg/count"
	"net/http"
)

// URL: /api/browsers
var GetBrowsersHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results := count.Browsers(before, after, getRequestedLimit(r))
	return respond(w, envelope{Data: results})
})
