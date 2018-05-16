package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/config"
	"github.com/usefathom/fathom/pkg/datastore"
)

type App struct {
	*cli.App
	database datastore.Datastore
	config   *config.Config
}

var app *App

func main() {
	// force all times in UTC, regardless of server timezone
	os.Setenv("TZ", "")

	// setup CLI app
	app = &App{cli.NewApp(), nil, nil}
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
	app.config = config.Parse(c.String("config"))
	app.database = datastore.New(app.config.Database)
	return nil
}

func after(c *cli.Context) error {
	app.database.Close()
	return nil
}
