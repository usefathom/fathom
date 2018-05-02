package commands

import (
	"github.com/usefathom/fathom/pkg/count"
)

// Archive processes unarchived data (pageviews to aggregated count tables)
func Archive() {
	count.Archive()
}
