package main

import (
	"log"

	"github.com/dannyvankooten/ana/commands"
	"github.com/dannyvankooten/ana/db"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// setup database connection
	conn := db.SetupDatabaseConnection()
	defer conn.Close()

	// parse & run cli commands
	commands.Parse()
	commands.Run()
}
