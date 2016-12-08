package commands

import(
  "flag"
)

var runCreateUserCommand bool
var runStartServerCommand bool
var emailArg string
var passwordArg string

func Parse() {
  // parse commands
  flag.BoolVar(&runCreateUserCommand, "create_user", false, "Create a new user")
  flag.BoolVar(&runStartServerCommand, "start_server", true, "Start the API web server")
  flag.StringVar(&emailArg, "email", "", "Email address")
  flag.StringVar(&passwordArg, "password", "", "Password")
  flag.Parse()
}

func Run() {
  if runCreateUserCommand {
    CreateUser()
  }

  if runStartServerCommand {
    StartServer()
  }
}
