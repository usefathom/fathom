package commands

import(
  "github.com/dannyvankooten/ana/db"
)

func seedDatabase() {
  db.Seed(nArg)
}
