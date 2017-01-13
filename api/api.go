package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultPeriod = 7
const defaultLimit = 10

type envelope struct {
	Data  interface{}
	Error interface{}
}

func respond(w http.ResponseWriter, d interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(d)
}

// log fatal errors
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
