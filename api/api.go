package api

import (
  "net/http"
)

func RegisterRoutes() {
  http.HandleFunc("/api/visits", GetVisitsHandler)
}
