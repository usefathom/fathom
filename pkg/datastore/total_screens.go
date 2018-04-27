package datastore

// TotalUniqueScreens returns the total # of screens between two given timestamps
func TotalUniqueScreens(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
	    SELECT
	    	IFNULL( SUM(t.count_unique), 0 )
	    FROM total_screens t
	    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	return total, err
}
