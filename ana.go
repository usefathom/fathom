package main

import (
  "net/http"

  "github.com/dannyvankooten/ana/core"
  "github.com/dannyvankooten/ana/api"
  "github.com/gorilla/mux"
)

// TODO: Authentication.

func main() {
    db := core.SetupDatabaseConnection()
    defer db.Close()

    r := mux.NewRouter()

    // register routes
    r.HandleFunc("/collect", api.CollectHandler).Methods("GET")
    r.HandleFunc("/api/visits/count/day", api.GetVisitsDayCountHandler).Methods("GET")
    r.HandleFunc("/api/visits/count/realtime", api.GetVisitsRealtimeCount).Methods("GET")
    r.HandleFunc("/api/visits", api.GetVisitsHandler).Methods("GET")
    r.HandleFunc("/api/pageviews", api.GetPageviewsHandler).Methods("GET")
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    r.Handle("/", http.FileServer(http.Dir("./views/")))

    // r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //   http.ServeFile(w, r, "./static/" + r.URL.Path[1:])
    // })

    http.ListenAndServe(":8080", r)
}
