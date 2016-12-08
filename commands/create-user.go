package commands

import(
  "github.com/dannyvankooten/ana/db"
  "golang.org/x/crypto/bcrypt"
  "log"
)

func CreateUser() {
  if emailArg == "" || passwordArg == "" {
    log.Fatal("Please supply -email and -password values")
  }

  stmt2, _ := db.Conn.Prepare("INSERT INTO users(email, password) VALUES(?, ?)")
  hash, _ := bcrypt.GenerateFromPassword([]byte(passwordArg), 10)
  stmt2.Exec(emailArg, hash)
}
