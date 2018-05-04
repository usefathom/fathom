package datastore

// TotalBounces returns the total number of pageviews between the given timestamps
func TotalBounces(before int64, after int64) (int64, error) {
	var total int64

	query := dbx.Rebind(`
		SELECT COALESCE(ROUND(AVG(t.count), 0), 0)
		FROM total_bounced t 
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// TotalUniqueBounces returns the total number of unique pageviews between the given timestamps
func TotalUniqueBounces(before int64, after int64) (int64, error) {
	var total int64

	query := dbx.Rebind(`
		SELECT COALESCE(AVG(t.count_unique), 0)
		FROM total_bounced t 
		WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)
	err := dbx.Get(&total, query, before, after)
	if err != nil {
		return 0, err
	}

	return total, nil
}
