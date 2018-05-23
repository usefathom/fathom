package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func register(c *cli.Context) error {
	email := c.String("email")
	if email == "" {
		return errors.New("invalid args: missing email address")
	}

	password := c.String("password")
	if password == "" {
		return errors.New("invalid args: missing password")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &models.User{
		Email:    email,
		Password: string(hash),
	}
	err := app.database.SaveUser(user)

	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	log.Infof("created user %s", user.Email)
	return nil
}
