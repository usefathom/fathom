package main

import (
	"log"

	"github.com/dannyvankooten/ana/commands"
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func main() {
	log.Println("starting...")

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}

	// setup database connection
	conn, err := db.SetupDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// setup cron to run count.Archive every hour
	c := cron.New()
	c.AddFunc("@hourly", count.Archive)
	c.Start()

	// parse & run cli commands
	commands.Parse()
	commands.Run()
}
