package commands

import (
	"flag"
)

var runCreateUserCommand bool
var runDeleteUserCommand bool
var runStartServerCommand bool
var runSeedDataCommand bool
var runArchiveDataCommand bool
var idArg int
var emailArg string
var passwordArg string
var nArg int

func Parse() {
	// parse commands
	flag.BoolVar(&runCreateUserCommand, "create_user", false, "Create a new user")
	flag.BoolVar(&runDeleteUserCommand, "delete_user", false, "Deletes a user")
	flag.BoolVar(&runStartServerCommand, "start_server", true, "Start the API web server")
	flag.BoolVar(&runSeedDataCommand, "seed_data", false, "Seed the database -n times")
	flag.BoolVar(&runArchiveDataCommand, "archive_data", false, "Archives data into daily aggregated totals")
	flag.StringVar(&emailArg, "email", "", "Email address")
	flag.StringVar(&passwordArg, "password", "", "Password")
	flag.IntVar(&idArg, "id", 0, "Object ID")
	flag.IntVar(&nArg, "n", 0, "Number")
	flag.Parse()
}

func Run() {
	if runCreateUserCommand {
		createUser()
	} else if runDeleteUserCommand {
		deleteUser()
	} else if runSeedDataCommand {
		seedData()
	} else if runArchiveDataCommand {
		archiveData()
	} else if runStartServerCommand {
		startServer()
	}
}
