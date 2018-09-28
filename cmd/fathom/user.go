package main

import (
	"errors"
	"fmt"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
	"golang.org/x/crypto/ssh/terminal"
)

var userCmd = cli.Command{
	Name:   "user",
	Usage:  "manage registered admin users",
	Action: userAdd,
	Subcommands: []cli.Command{
		cli.Command{
			Name:    "add",
			Aliases: []string{"register"},
			Action:  userAdd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "email, e",
					Usage: "user email",
				},
				cli.BoolFlag{
					Name:  "skip-bcrypt",
					Usage: "store password string as-is, skipping bcrypt",
				},
			},
		},
		cli.Command{
			Name:   "delete",
			Action: userDelete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "email, e",
					Usage: "user email",
				},
			},
		},
	},
}

func userAdd(c *cli.Context) error {
	email := c.String("email")
	if email == "" {
		return errors.New("Invalid arguments: missing email")
	}

	// Read and validate password
	fmt.Print("\nPassword: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n\n")
	if err != nil {
		return fmt.Errorf("Error reading password: %s", err)
	}

	password := string(bytePassword)
	if len(password) < 8 {
		return errors.New("Invalid password: must be at least 8 characters")
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

func userDelete(c *cli.Context) error {
	email := c.String("email")
	if email == "" {
		return errors.New("Invalid arguments: missing email")
	}

	user, err := app.database.GetUserByEmail(email)
	if err != nil {
		if err == datastore.ErrNoResults {
			return fmt.Errorf("No user with email %s", email)
		}

		return err
	}

	if err := app.database.DeleteUser(user); err != nil {
		return err
	}

	log.Infof("Deleted user %s", user.Email)

	return nil
}
