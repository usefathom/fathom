package commands

import (
	"github.com/usefathom/fathom/pkg/count"
	"github.com/usefathom/fathom/pkg/datastore"
)

// Seed creates n database records with dummy data
func Seed(n int) {
	datastore.Seed(n)
}

// Archive processes unarchived data (pageviews to aggregated count tables)
func Archive() {
	count.Archive()
}
