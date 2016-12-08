package models

import (
  "database/sql"
)

type Site struct {
  ID int64
  Url string
}

func (s *Site) Save(conn *sql.DB) error {
    // prepare statement for inserting data
    stmt, err := conn.Prepare(`INSERT INTO sites(
        url
      ) VALUES(?)`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    result, err := stmt.Exec(s.Url)
    s.ID, _ = result.LastInsertId()

    return err
}
