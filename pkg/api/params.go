package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Params defines the commonly used API parameters
type Params struct {
	SiteID    int64
	Offset    int64
	Limit     int64
	StartDate time.Time
	EndDate   time.Time
}

// GetRequestParams parses the query parameters and returns commonly used API parameters, with defaults
func GetRequestParams(r *http.Request) *Params {
	params := &Params{
		SiteID:    0,
		Limit:     20,
		Offset:    0,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, -7),
	}

	vars := mux.Vars(r)
	if _, ok := vars["id"]; ok {
		if siteID, err := strconv.ParseInt(vars["id"], 10, 64); err == nil {
			params.SiteID = siteID
		}
	}

	q := r.URL.Query()
	if q.Get("after") != "" {
		if after, err := strconv.ParseInt(q.Get("after"), 10, 64); err == nil && after > 0 {
			params.StartDate = time.Unix(after, 0)
		}
	}

	if q.Get("before") != "" {
		if before, err := strconv.ParseInt(q.Get("before"), 10, 64); err == nil && before > 0 {
			params.EndDate = time.Unix(before, 0)
		}
	}

	if q.Get("limit") != "" {
		if limit, err := strconv.ParseInt(q.Get("limit"), 10, 64); err == nil && limit > 0 {
			params.Limit = limit
		}
	}

	if q.Get("offset") != "" {
		if offset, err := strconv.ParseInt(q.Get("offset"), 10, 64); err == nil && offset > 0 {
			params.Offset = offset
		}
	}

	return params
}

func parseMajorMinor(v string) string {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) > 1 {
		v = parts[0] + "." + parts[1]
	}
	return v
}
