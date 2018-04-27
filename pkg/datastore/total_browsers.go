package datastore

// TotalUniqueBrowsers returns the total # of unique browsers between two given timestamps
func TotalUniqueBrowsers(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
	    SELECT
	      SUM(IFNULL(t.count_unique, 0))
	    FROM total_browser_names t
	    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	return total, err
}
