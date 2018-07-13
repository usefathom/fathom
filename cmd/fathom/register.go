package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var registerCmd = cli.Command{
	Name:    "register",
	Aliases: []string{"r"},
	Usage:   "register a new admin user",
	Action:  register,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "email, e",
			Usage: "user email",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "user password",
		},
	},
}

func register(c *cli.Context) error {
	email := c.String("email")
	if email == "" {
		return errors.New("Invalid arguments: missing email address")
	}

	password := c.String("password")
	if password == "" {
		return errors.New("Invalid arguments: missing password")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := &models.User{
		Email:    email,
		Password: string(hash),
	}
	err := app.database.SaveUser(user)

	if err != nil {
		return fmt.Errorf("Error creating user: %s", err)
	}

	log.Infof("Created user %s", user.Email)
	return nil
}
