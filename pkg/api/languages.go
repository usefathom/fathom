package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/count"
)

// URL: /api/languages
var GetLanguagesHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	before, after := getRequestedPeriods(r)
	results := count.Languages(before, after, getRequestedLimit(r))
	respond(w, envelope{Data: results})
})
