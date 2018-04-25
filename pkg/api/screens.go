package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
)

// URL: /api/screen-resolutions
var GetScreenResolutionsHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results := count.Screens(before, after, getRequestedLimit(r))
	return respond(w, envelope{Data: results})
})
