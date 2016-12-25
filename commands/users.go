package commands

import (
	"log"

	"github.com/dannyvankooten/ana/db"
	"github.com/dannyvankooten/ana/models"
	"golang.org/x/crypto/bcrypt"
)

func createUser() {
	if emailArg == "" || passwordArg == "" {
		log.Fatal("Please supply -email and -password values")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(passwordArg), 10)
	user := models.User{
		Email:    emailArg,
		Password: string(hash),
	}
	user.Save(db.Conn)

	log.Printf("User %s #%d created", emailArg, user.ID)
}

func deleteUser() {
	if emailArg == "" && idArg == 0 {
		log.Fatal("Please supply an -email or -id value")
	}

	stmt2, _ := db.Conn.Prepare("DELETE FROM users WHERE email = ? OR id = ?")
	stmt2.Exec(emailArg, idArg)

	log.Printf("User with email %s or ID %d deleted", emailArg, idArg)
}
