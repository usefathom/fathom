package main

import (
  "net/http"
  "os"
  "log"
  "github.com/dannyvankooten/ana/db"
  "github.com/dannyvankooten/ana/api"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "github.com/joho/godotenv"
)

func main() {
  // load .env file
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  conn := db.SetupDatabaseConnection()
  defer conn.Close()

  // register routes
  r := mux.NewRouter()
  r.HandleFunc("/collect", api.CollectHandler).Methods("GET")
  r.Handle("/api/session", api.Login).Methods("POST")
  r.Handle("/api/session", api.Logout).Methods("DELETE")
  r.Handle("/api/visits/count", api.Authorize(api.GetVisitsCountHandler)).Methods("GET")
  r.Handle("/api/visits/count/group/{period}", api.Authorize(api.GetVisitsPeriodCountHandler)).Methods("GET")
  r.Handle("/api/visits/count/realtime", api.Authorize(api.GetVisitsRealtimeCountHandler)).Methods("GET")
  r.Handle("/api/visits", api.Authorize(api.GetVisitsHandler)).Methods("GET")
  r.Handle("/api/pageviews/count", api.Authorize(api.GetPageviewsCountHandler)).Methods("GET")
  r.Handle("/api/pageviews/count/group/{period}", api.Authorize(api.GetPageviewsPeriodCountHandler)).Methods("GET")
  r.Handle("/api/pageviews", api.Authorize(api.GetPageviewsHandler)).Methods("GET")
  r.Handle("/api/languages", api.Authorize(api.GetLanguagesHandler)).Methods("GET")
  r.Handle("/api/screen-resolutions", api.Authorize(api.GetScreenResolutionsHandler)).Methods("GET")
  r.Handle("/api/countries", api.Authorize(api.GetCountriesHandler)).Methods("GET")
  r.Handle("/api/browsers", api.Authorize(api.GetBrowsersHandler)).Methods("GET")

  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
  r.Path("/tracker.js").Handler(http.FileServer(http.Dir("./static/js/")))
  r.Handle("/", http.FileServer(http.Dir("./views/")))

  http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
