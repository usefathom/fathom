package cli

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/gorilla/handlers"
	"github.com/usefathom/fathom/pkg/api"
	"golang.org/x/crypto/acme/autocert"
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
			EnvVar: "FATHOM_LETS_ENCRYPT",
			Name:   "lets-encrypt",
		},

		cli.StringFlag{
			EnvVar: "FATHOM_HOSTNAME",
			Name:   "hostname",
			Usage:  "domain when using --lets-encrypt",
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

	// start server without letsencrypt / tls enabled
	if !c.Bool("lets-encrypt") {
		// start listening
		server := &http.Server{
			Addr:         addr,
			Handler:      h,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		log.Infof("Server is now listening on %s", server.Addr)
		log.Fatal(server.ListenAndServe())
		return nil
	}

	// start server with autocert (letsencrypt)
	hostname := c.String("hostname")
	log.Infof("Server is now listening on %s:443", hostname)
	log.Fatal(http.Serve(autocert.NewListener(hostname), h))
	return nil
}
