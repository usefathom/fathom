package commands

import (
	"flag"
)

var (
	runCreateUserCommand  bool
	runDeleteUserCommand  bool
	runStartServerCommand bool
	runSeedDataCommand    bool
	runArchiveDataCommand bool
	idArg                 int
	emailArg              string
	passwordArg           string
	nArg                  int
)

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
