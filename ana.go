package main

import (
	"github.com/dannyvankooten/ana/pkg/commands"
	"github.com/dannyvankooten/ana/pkg/count"
	"github.com/dannyvankooten/ana/pkg/datastore"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	app              = kingpin.New("ana", "Open-source web analytics.")
	register         = app.Command("register", "Register a new user.")
	registerEmail    = register.Arg("email", "Email for user.").Required().String()
	registerPassword = register.Arg("password", "Password for user.").Required().String()
	server           = app.Command("server", "Start webserver.").Default()
	serverPort       = server.Flag("port", "Port to listen on.").Default("8080").Int()
	archive          = app.Command("archive", "Process unarchived data.")
	seed             = app.Command("seed", "Seed the database.")
	seedN            = seed.Flag("n", "Number of records to seed.").Int()
)

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// setup database connection
	db := datastore.Init()
	defer db.Close()

	// setup cron to run count.Archive every hour
	c := cron.New()
	c.AddFunc("@hourly", count.Archive)
	c.Start()

	// parse & run cli commands
	app.Version("1.0")
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "register":
		commands.Register(*registerEmail, *registerPassword)

	case "server":
		commands.Server(*serverPort)

	case "archive":
		commands.Archive()

	case "seed":
		commands.Seed(*seedN)
	}

}
