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
		cli.BoolFlag{
			Name:  "skip-bcrypt",
			Usage: "store password string as is, skipping bcrypt",
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

	// set password manually if --skip-bcrypt was given
	// this is used to supply an already encrypted password string
	if c.Bool("skip-bcrypt") {
		user.Password = password
	}

	if err := app.database.SaveUser(&user); err != nil {
		return fmt.Errorf("Error creating user: %s", err)
	}

	log.Infof("Created user %s", user.Email)
	return nil
}
