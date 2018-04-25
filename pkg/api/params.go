package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultPeriod = 7
const defaultLimit = 10

func getRequestedLimit(r *http.Request) int {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit == 0 {
		limit = 10
	}

	return limit
}

func getRequestedPeriods(r *http.Request) (int64, int64) {
	var before, after int64
	var err error

	before, err = strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	if err != nil || before == 0 {
		before = time.Now().Unix()
	}

	after, err = strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
	if err != nil || before == 0 {
		after = time.Now().AddDate(0, 0, -7).Unix()
	}

	return before, after
}

func parseMajorMinor(v string) string {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) > 1 {
		v = parts[0] + "." + parts[1]
	}
	return v
}

func getRequestIp(r *http.Request) string {
	ipAddress := r.RemoteAddr

	headerForwardedFor := r.Header.Get("X-Forwarded-For")
	if headerForwardedFor != "" {
		ipAddress = headerForwardedFor
	}

	return ipAddress
}
