package commands

import(
  "github.com/dannyvankooten/ana/db"
)

func seedData() {
  db.Seed(nArg)
}

func archiveData() {
  db.CreateArchives()
}
