package api

import (
  "net/http"
  "github.com/gorilla/sessions"
  "os"
)

var store = sessions.NewFilesystemStore( "./storage/sessions/", []byte(os.Getenv("ANA_SECRET_KEY")))

// URL: POST /api/session
var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth")
    session.Values["user"] = "Danny"
    err := session.Save(r, w)
    checkError(err)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("true"))
})

// URL: DELETE /api/session
var Logout = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth")
    if ! session.IsNew  {
      session.Options.MaxAge = -1
      session.Save(r, w)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("true"))
})

/* middleware */
func Authorize(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth")

    if _, ok := session.Values["user"]; !ok  {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }

    next.ServeHTTP(w, r)
  })
}
