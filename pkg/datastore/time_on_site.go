package datastore

// TODO: Store in visitors table? Or other table?
func AvgTimeOnSite(before int64, after int64) (int64, error) {
	query := dbx.Rebind(`
   SELECT ROUND(AVG(time_on_site)) FROM ( 
      SELECT SUM(time_on_page) AS time_on_site 
      FROM pageviews 
      WHERE time_on_page > 0 AND UNIX_TIMESTAMP(timestamp) < ? AND UNIX_TIMESTAMP(timestamp) > ?
      GROUP BY visitor_id 
   ) AS time_on_site_query`)

	var total int64
	err := dbx.Get(&total, query, before, after)
	if err != nil {
		return 0, err
	}

	return total, nil
}
