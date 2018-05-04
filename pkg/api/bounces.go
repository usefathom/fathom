package api

import (
	"net/http"

	"github.com/usefathom/fathom/pkg/datastore"
)

// URL: /api/bounces/count
var GetBouncesCountHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	before, after := getRequestedPeriods(r)
	result, err := datastore.TotalBounces(before, after)
	if err != nil {
		return err
	}

	return respond(w, envelope{Data: result})
})
