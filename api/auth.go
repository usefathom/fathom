package api

import (
  "net/http"
  "github.com/gorilla/sessions"
)

var store = sessions.NewFilesystemStore( "./storage/sessions/", []byte("something-very-secret"))

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
    session.Options.MaxAge = -1
    err := session.Save(r, w)
    checkError(err)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("true"))
})

/* middleware */
func Authorize(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "auth")
    checkError(err)

    if user, ok := session.Values["user"]; !ok  {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }

    next.ServeHTTP(w, r)
  })
}
