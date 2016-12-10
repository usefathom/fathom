package models

import(
  "database/sql"
)

type Page struct {
  ID int64
  Hostname string
  Path string
  Title string
}

func (p *Page) Save(conn *sql.DB) error {
    // prepare statement for inserting data
    stmt, err := conn.Prepare(`INSERT INTO pages(
      hostname,
      path,
      title
      ) VALUES( ?, ?, ? )`)
    if err != nil {
        return err
    }
    defer stmt.Close()

    result, err := stmt.Exec(p.Hostname, p.Path, p.Title)
    p.ID, _ = result.LastInsertId()

    return err
  }
