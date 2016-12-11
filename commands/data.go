package commands

import(
  "github.com/dannyvankooten/ana/db"
  "github.com/dannyvankooten/ana/count"
)

func seedData() {
  db.Seed(nArg)
}

func archiveData() {
  count.CreateArchives()
}
