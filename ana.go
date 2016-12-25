package main

import (
	"log"

	"github.com/dannyvankooten/ana/commands"
	"github.com/dannyvankooten/ana/db"
	"github.com/joho/godotenv"
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

	// parse & run cli commands
	commands.Parse()
	commands.Run()
}
