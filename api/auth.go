package api

import (
  "net/http"
  "github.com/gorilla/sessions"
  "github.com/dannyvankooten/ana/db"
  "github.com/dannyvankooten/ana/models"
  "golang.org/x/crypto/bcrypt"
  "os"
  "encoding/json"
)

type Login struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

var store = sessions.NewFilesystemStore("./storage/sessions/", []byte(os.Getenv("ANA_SECRET_KEY")))

// URL: POST /api/session
var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  var l Login
  json.NewDecoder(r.Body).Decode(&l)

  // check login creds
  var hashedPassword string
  var u models.User
  stmt, _ := db.Conn.Prepare("SELECT id, email, password FROM users WHERE email = ? LIMIT 1")
  err := stmt.QueryRow(l.Email).Scan(&u.ID, &u.Email, &hashedPassword)

  if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(l.Password)) != nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  // TODO: Replace session filesystem store with DB store.
  session, _ := store.Get(r, "auth")
  session.Values["user_id"] = u.ID
  err = session.Save(r, w)
  checkError(err)

  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte("true"))
})

// URL: DELETE /api/session
var LogoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth")
    if ! session.IsNew  {
      session.Options.MaxAge = -1
      session.Save(r, w)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("true"))
})

/* middleware */
func Authorize(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "auth")
    userID, ok := session.Values["user_id"];

    if !ok  {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }

    // find user
    var u models.User
    stmt, _ := db.Conn.Prepare("SELECT id, email FROM users WHERE id = ? LIMIT 1")
    err := stmt.QueryRow(userID).Scan(&u.ID, &u.Email)
    if err != nil {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }

    next.ServeHTTP(w, r)
  })
}
