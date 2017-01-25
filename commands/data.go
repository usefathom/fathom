package commands

import (
	"github.com/dannyvankooten/ana/count"
	"github.com/dannyvankooten/ana/datastore"
)

// Seed creates n database records with dummy data
func Seed(n int) {
	datastore.Seed(n)
}

// Archive processes unarchived data (pageviews to aggeegated count tables)
func Archive() {
	count.Archive()
}
