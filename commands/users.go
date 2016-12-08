package commands

import(
  "github.com/dannyvankooten/ana/db"
  "golang.org/x/crypto/bcrypt"
  "log"
)

func createUser() {
  if emailArg == "" || passwordArg == "" {
    log.Fatal("Please supply -email and -password values")
  }

  stmt2, _ := db.Conn.Prepare("INSERT INTO users(email, password) VALUES(?, ?)")
  hash, _ := bcrypt.GenerateFromPassword([]byte(passwordArg), 10)
  stmt2.Exec(emailArg, hash)

    log.Printf("User %s created", emailArg)
}

func deleteUser() {
  if emailArg == "" && idArg == 0 {
    log.Fatal("Please supply an -email or -id value")
  }

  stmt2, _ := db.Conn.Prepare("DELETE FROM users WHERE email = ? OR id = ?")
  stmt2.Exec(emailArg, idArg)

  log.Printf("User with email %s or ID %d deleted", emailArg, idArg)
}
