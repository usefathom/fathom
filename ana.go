package main

import (
  "net/http"

  "github.com/dannyvankooten/ana/core"
  "github.com/dannyvankooten/ana/api"
)


func main() {
    db := core.SetupDatabaseConnection()
    defer db.Close()

    // register routes
    api.RegisterRoutes()
    http.HandleFunc("/collect", api.CollectHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, "./static/" + r.URL.Path[1:])
    })
    http.ListenAndServe(":8080", nil)
}
