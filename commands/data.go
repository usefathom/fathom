package commands

import (
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
)

func seedData() {
	db.Seed(nArg)
}

func archiveData() {
	count.CreatePageviewArchives()
	count.CreateVisitorArchives()
	count.CreatePageviewArchivesPerPage()
	count.CreateScreenArchives()
	count.CreateLanguageArchives()
	count.CreateBrowserArchives()
	count.CreateReferrerArchives()
}
