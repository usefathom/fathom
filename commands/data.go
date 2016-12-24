package commands

import (
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
)

func seedData() {
	db.Seed(nArg)
}

func archiveData() {
	count.CreateVisitorArchives()
	count.CreatePageviewArchives()
	count.CreateScreenArchives()
	count.CreateLanguageArchives()
	count.CreateBrowserArchives()
	count.CreateReferrerArchives()
}
