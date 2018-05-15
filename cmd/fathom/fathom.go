package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/datastore"
)

var db *sqlx.DB
var config *Config

func main() {
	app := cli.NewApp()
	app.Name = "Fathom"
	app.Usage = "simple & transparent website analytics"
	app.Version = "1.0.0"
	app.HelpName = "fathom"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: ".env",
			Usage: "Load configuration from `FILE`",
		},
	}
	app.Before = before
	app.After = after
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the fathom web server",
			Action:  server,
			Flags: []cli.Flag{
				cli.StringFlag{
					EnvVar: "FATHOM_SERVER_ADDR",
					Name:   "addr,port",
					Usage:  "server address",
					Value:  ":8080",
				},
				cli.BoolFlag{
					EnvVar: "FATHOM_DEBUG",
					Name:   "debug, d",
				},
			},
		},
		{
			Name:      "register",
			Aliases:   []string{"r"},
			Usage:     "register a new admin user",
			ArgsUsage: "<email> <password>",
			Action:    register,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func before(c *cli.Context) error {
	config = parseConfig(c.String("config"))
	db = datastore.Init(config.Database)
	return nil
}

func after(c *cli.Context) error {
	db.Close()
	return nil
}
