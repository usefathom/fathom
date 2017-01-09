package commands

import (
	"flag"

	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
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
	portArg               int
)

// Parse CLI arguments
func Parse() {
	// parse commands
	flag.BoolVar(&runCreateUserCommand, "create_user", false, "Create a new user")
	flag.BoolVar(&runDeleteUserCommand, "delete_user", false, "Deletes a user")
	flag.BoolVar(&runStartServerCommand, "start_server", true, "Start the API web server, listen on -port")
	flag.BoolVar(&runSeedDataCommand, "db_seed", false, "Seed the database -n times")
	flag.BoolVar(&runArchiveDataCommand, "archive_data", false, "Aggregates data into daily totals")
	flag.StringVar(&emailArg, "email", "", "Email address")
	flag.StringVar(&passwordArg, "password", "", "Password")
	flag.IntVar(&idArg, "id", 0, "Object ID")
	flag.IntVar(&nArg, "n", 0, "Number")
	flag.IntVar(&portArg, "port", 8080, "Port")
	flag.Parse()
}

// Run parsed CLI command. Defaults to starting the HTTP server.
func Run() {
	if runCreateUserCommand {
		createUser()
	} else if runDeleteUserCommand {
		deleteUser()
	} else if runSeedDataCommand {
		db.Seed(nArg)
	} else if runArchiveDataCommand {
		count.Archive()
	} else if runStartServerCommand {
		startServer(portArg)
	}
}
