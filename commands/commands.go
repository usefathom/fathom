package commands

import(
  "flag"
)

var runCreateUserCommand bool
var runDeleteUserCommand bool
var runStartServerCommand bool
var runSeedDatabaseCommand bool
var idArg int
var emailArg string
var passwordArg string
var nArg int

func Parse() {
  // parse commands
  flag.BoolVar(&runCreateUserCommand, "create_user", false, "Create a new user")
  flag.BoolVar(&runDeleteUserCommand, "delete_user", false, "Deletes a user")
  flag.BoolVar(&runStartServerCommand, "start_server", true, "Start the API web server")
  flag.BoolVar(&runSeedDatabaseCommand, "seed_database", false, "Seed the database -n times")
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
  } else if runSeedDatabaseCommand {
    seedDatabase()
  } else if runStartServerCommand {
    startServer()
  }
}
