package datastore

// TotalReferrers returns the total # of referrers between two given timestamps
func TotalReferrers(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
	    SELECT
	      IFNULL( SUM(t.count), 0 )
	    FROM total_referrers t
	    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)

	err := dbx.Get(&total, query, before, after)
	return total, err
}
