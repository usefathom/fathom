package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/usefathom/fathom/pkg/api"
)

func server(c *cli.Context) error {
	var h http.Handler
	h = api.Routes()

	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		h = handlers.LoggingHandler(log.StandardLogger().Writer(), h)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	addr := c.String("addr")
	log.Infof("Server listening on %s", addr)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Errorln(err)
	}
	return nil
}
