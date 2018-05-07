package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Params struct {
	Limit     int
	StartDate time.Time
	EndDate   time.Time
}

func GetRequestParams(r *http.Request) *Params {
	params := &Params{
		Limit:     20,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, -7),
	}

	q := r.URL.Query()
	if after, err := strconv.ParseInt(q.Get("after"), 10, 64); err == nil && after > 0 {
		params.StartDate = time.Unix(after, 0)
	}

	if before, err := strconv.ParseInt(q.Get("before"), 10, 64); err == nil && before > 0 {
		params.EndDate = time.Unix(before, 0)
	}

	if limit, err := strconv.Atoi(q.Get("limit")); err == nil && limit > 0 {
		params.Limit = limit
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
