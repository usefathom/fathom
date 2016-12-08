package commands

import(
  "flag"
)

var runCreateUserCommand bool
var runDeleteUserCommand bool
var runStartServerCommand bool
var idArg int
var emailArg string
var passwordArg string

func Parse() {
  // parse commands
  flag.BoolVar(&runCreateUserCommand, "create_user", false, "Create a new user")
  flag.BoolVar(&runDeleteUserCommand, "delete_user", false, "Deletes a user")
  flag.BoolVar(&runStartServerCommand, "start_server", true, "Start the API web server")
  flag.StringVar(&emailArg, "email", "", "Email address")
  flag.StringVar(&passwordArg, "password", "", "Password")
  flag.IntVar(&idArg, "id", 0, "Object ID")
  flag.Parse()
}

func Run() {
  if runCreateUserCommand {
    createUser()
  } else if runDeleteUserCommand {
    deleteUser()
  } else if runStartServerCommand {
    startServer()
  }
}
