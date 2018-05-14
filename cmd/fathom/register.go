package main

import (
	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func register(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(c.String("password")), 10)
	user := &models.User{
		Email:    c.String("email"),
		Password: string(hash),
	}
	err := datastore.SaveUser(user)

	if err != nil {
		log.Printf("Error creating user: %s", err)
	} else {
		log.Printf("User %s #%d created.\n", user.Email, user.ID)
	}
	return nil
}
