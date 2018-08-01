package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/models"
)

var registerCmd = cli.Command{
	Name:   "register",
	Usage:  "register a new admin user",
	Action: register,
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
		return errors.New("Invalid arguments: missing email")
	}

	password := c.String("password")
	if password == "" {
		return errors.New("Invalid arguments: missing password")
	}

	user := models.NewUser(email, password)
	err := app.database.SaveUser(&user)

	if err != nil {
		return fmt.Errorf("Error creating user: %s", err)
	}

	log.Infof("Created user %s", user.Email)
	return nil
}
