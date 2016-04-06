package users

import (
  // "strings"
  // "net/url"
  "github.com/db-docker/myDB"
  "strconv"
  "fmt"
  // "database/sql"
  // "github.com/go-sql-driver/mysql"
  // "golang.org/x/crypto/bcrypt"
  // "github.com/dgrijalva/jwt-go"
  // "time"
)

type User struct {
  Id string
  Username string
  Password string
}

func GetUser(username string) User {
  var retrievedUser User;
  queryString := "SELECT id, username, password FROM users WHERE username = ?"
  rows, err := myDB.DBCon.Query(queryString, username)
  
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
  stmt, err := myDB.DBCon.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
  if err != nil {
    fmt.Println("preparing to write to DB", err)
  }
  user, err := stmt.Exec(username, password)
  if err != nil {
    // if mysqlError, ok := err.(*mysql.MySQLError); ok {
    //   if mysqlError.Number == 1062 {
    //     return nil
    //   } else {
    //     fmt.Println("error writing to DB", err)
    //   }
    // }
    fmt.Println("error writing to DB", err)
  }
  insertId, _ := user.LastInsertId()
  return []byte(strconv.Itoa(int(insertId)))
}