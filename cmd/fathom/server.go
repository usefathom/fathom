package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/gorilla/handlers"
	"github.com/usefathom/fathom/pkg/api"
)

func server(c *cli.Context) error {
	var h http.Handler
	h = api.Routes()

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
	log.Infof("Server listening on %s", addr)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Errorln(err)
	}
	return nil
}
