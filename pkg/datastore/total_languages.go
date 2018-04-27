package datastore

// TotalUniqueLanguages returns the total # of unique browser languages between two given timestamps
func TotalUniqueLanguages(before int64, after int64) (int, error) {
	var total int

	query := dbx.Rebind(`
    SELECT
      IFNULL( SUM(t.count_unique), 0 )
    FROM total_browser_languages t
    WHERE UNIX_TIMESTAMP(t.date) <= ? AND UNIX_TIMESTAMP(t.date) >= ?`)

	err := dbx.Get(&total, query, before, after)
	return total, err
}
