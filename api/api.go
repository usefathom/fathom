package api

import (
  "net/http"
)

func RegisterRoutes() {
  http.HandleFunc("/api/visits/count/realtime", GetVisitsRealtimeCount)
  http.HandleFunc("/api/visits", GetVisitsHandler)
}
