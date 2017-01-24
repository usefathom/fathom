package commands

import (
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/db"
)

// Seed creates n database records with dummy data
func Seed(n int) {
	db.Seed(n)
}

// Archive processes unarchived data (pageviews to aggeegated count tables)
func Archive() {
	count.Archive()
}
