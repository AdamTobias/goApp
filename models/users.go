package users

import (
  "github.com/db-docker/common"
  "strconv"
  "fmt"
)

func GetUser(username string) common.User {
  var retrievedUser common.User;
  queryString := "SELECT id, username, password FROM users WHERE username = ?"
  rows, err := common.DBCon.Query(queryString, username)
  
  if err != nil {
    fmt.Println("querying DB", err)
  }

  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&retrievedUser.Id, &retrievedUser.Username, &retrievedUser.Password)
    if err != nil {
      fmt.Println("scanning rows", err)
    }
  }
  return retrievedUser;
}

func AddUser (username string, password string) []byte {
  stmt, err := common.DBCon.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
  if err != nil {
    fmt.Println("preparing to write to DB", err)
  }
  user, err := stmt.Exec(username, password)
  if err != nil {
    fmt.Println("error writing to DB", err)
  }
  insertId, _ := user.LastInsertId()
  return []byte(strconv.Itoa(int(insertId)))
}