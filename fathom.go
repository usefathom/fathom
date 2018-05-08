package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/usefathom/fathom/pkg/commands"
	"github.com/usefathom/fathom/pkg/counter"
	"github.com/usefathom/fathom/pkg/datastore"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Database *datastore.Config

	Secret string `required:"true"`
}

var (
	app              = kingpin.New("fathom", "Simple website analytics.")
	register         = app.Command("register", "Register a new user.")
	registerEmail    = register.Arg("email", "Email for user.").Required().String()
	registerPassword = register.Arg("password", "Password for user.").Required().String()
	server           = app.Command("server", "Start webserver.").Default()
	serverPort       = server.Flag("port", "Port to listen on.").Default("8080").Int()
	serverWebRoot    = server.Flag("webroot", "Root directory of static assets").Default("./").String()
	archive          = app.Command("archive", "Process unarchived data.")
)

func main() {
	// load .env file
	var cfg Config
	godotenv.Load()
	err := envconfig.Process("Fathom", &cfg)
	if err != nil {
		log.Fatalf("Error parsing Fathom config from environment: %s", err)
	}

	// setup database connection
	db := datastore.Init(cfg.Database)
	defer db.Close()

	// parse & run cli commands
	app.Version("1.0")
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "register":
		commands.Register(*registerEmail, *registerPassword)

	case "server":
		commands.Server(*serverPort, *serverWebRoot)

	case "archive":
		err := counter.Aggregate()
		if err != nil {
			log.Warn(err)
		}
	}

}
