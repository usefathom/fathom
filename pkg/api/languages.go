package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
)

// URL: /api/languages
var GetLanguagesHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	results := count.Languages(before, after, getRequestedLimit(r))
	return respond(w, envelope{Data: results})
})
