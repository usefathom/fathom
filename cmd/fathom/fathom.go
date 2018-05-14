package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/usefathom/fathom/pkg/datastore"
)

func main() {
	cfg := parseConfig()
	db := datastore.Init(cfg.Database)
	defer db.Close()

	app := cli.NewApp()
	app.Name = "Fathom"
	app.Usage = "simple & transparent website analytics"
	app.Version = "1.0.0"
	app.HelpName = "fathom"
	app.Flags = []cli.Flag{}
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
					Name:   "debug",
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
	}
}
