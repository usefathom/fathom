package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/gorilla/handlers"
	"github.com/usefathom/fathom/pkg/api"
)

var serverCmd = cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "start the fathom web server",
	Action:  server,
	Flags: []cli.Flag{
		cli.StringFlag{
			EnvVar: "FATHOM_SERVER_ADDR,PORT",
			Name:   "addr,port",
			Usage:  "server address",
			Value:  ":8080",
		},
		cli.BoolFlag{
			EnvVar: "FATHOM_DEBUG",
			Name:   "debug, d",
		},
	},
}

func server(c *cli.Context) error {
	var h http.Handler
	a := api.New(app.database, app.config.Secret)
	h = a.Routes()

	// set debug log level if --debug was passed
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		h = handlers.LoggingHandler(log.StandardLogger().Writer(), h)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	// if addr looks like a number, prefix with :
	addr := c.String("addr")
	if _, err := strconv.Atoi(addr); err == nil {
		addr = ":" + addr
	}

	// start listening
	log.Printf("Server is now listening on %s", addr)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Errorln(err)
	}
	return nil
}
