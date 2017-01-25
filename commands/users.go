package commands

import (
	"github.com/dannyvankooten/ana/datastore"
	"github.com/dannyvankooten/ana/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Register creates a new user with the given email & password
func Register(email string, password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := models.User{
		Email:    email,
		Password: string(hash),
	}
	user.Save(datastore.DB)

	log.Printf("User %s #%d created.\n", email, user.ID)
}
