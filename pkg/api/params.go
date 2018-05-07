package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: Move params into Params struct (with defaults)

const defaultPeriod = 7
const defaultLimit = 10

func getRequestedLimit(r *http.Request) int64 {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit == 0 {
		limit = defaultLimit
	}

	return limit
}

func getRequestedDatePeriods(r *http.Request) (time.Time, time.Time) {
	var startDate, endDate time.Time
	var err error

	beforeUnix, err := strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	if err != nil || beforeUnix == 0 {
		endDate = time.Now()
	} else {
		endDate = time.Unix(beforeUnix, 0)
	}

	afterUnix, err := strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
	if err != nil || afterUnix == 0 {
		startDate = endDate.AddDate(0, 0, -defaultPeriod)
	} else {
		startDate = time.Unix(afterUnix, 0)
	}

	return startDate, endDate
}

func parseMajorMinor(v string) string {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) > 1 {
		v = parts[0] + "." + parts[1]
	}
	return v
}
