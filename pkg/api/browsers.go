package api

import (
	"github.com/usefathom/fathom/pkg/count"
	"net/http"
)

// URL: /api/browsers
var GetBrowsersHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results, err := count.Browsers(before, after, getRequestedLimit(r))
	if err != nil {
		return err
	}
	return respond(w, envelope{Data: results})
})
