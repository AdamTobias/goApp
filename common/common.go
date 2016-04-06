package common

import "database/sql"

var (
    // DBCon is the connection handle
    // for the database
    DBCon *sql.DB
)

type User struct {
  Id string
  Username string
  Password string
}