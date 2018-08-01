package main

import (
	"os"

	log "github.com/sirupsen/logrus"
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
		serverCmd,
		registerCmd,
		statsCmd,
	}

	if len(os.Args) < 2 || os.Args[1] != "--version" {
		log.Printf("%s %s", app.Name, app.Version)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func before(c *cli.Context) error {
	configFile := c.String("config")
	config.LoadEnv(configFile)
	app.config = config.Parse()
	app.database = datastore.New(app.config.Database)
	return nil
}

func after(c *cli.Context) error {
	err := app.database.Close()
	return err
}
