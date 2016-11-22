package api

import (
  "net/http"
  "log"
)

func RegisterRoutes() {
  http.HandleFunc("/api/visits/count/day", GetVisitsDayCountHandler)
  http.HandleFunc("/api/visits/count/realtime", GetVisitsRealtimeCount)
  http.HandleFunc("/api/visits", GetVisitsHandler)
  http.HandleFunc("/api/pageviews", GetPageviewsHandler)
}

// log fatal errors
func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
