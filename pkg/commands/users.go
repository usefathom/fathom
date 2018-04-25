package commands

import (
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Register creates a new user with the given email & password
func Register(email string, password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &models.User{
		Email:    email,
		Password: string(hash),
	}
	err := datastore.SaveUser(user)

	if err != nil {
		log.Printf("Error creating user: %s", err)
	} else {
		log.Printf("User %s #%d created.\n", email, user.ID)
	}

}
