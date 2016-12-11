package commands

import (
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
)

func seedData() {
	db.Seed(nArg)
}

func archiveData() {
	count.CreateArchives()
}
