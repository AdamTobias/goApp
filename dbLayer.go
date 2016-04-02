package main

import (
  "net/http"
  // "strings"
  // "net/url"
  "strconv"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // "golang.org/x/crypto/bcrypt"
  // "github.com/dgrijalva/jwt-go"
  // "time"
  "encoding/json"
)

type database struct {
  con *sql.DB
}

// initializes the database connection
func (d *database) dbInit() {
  // db, err := sql.Open("mysql", "root:rodam@tcp(127.0.0.1:5587)/gopractice")
  dbConnection, err := sql.Open("mysql", "root:rodam@/gopractice")
  if err != nil {
    fmt.Println("error connecting to db")
    return
  }
  fmt.Println("Database connection established")
  d.con = dbConnection
}

//create a global db struct for anyone to use
var db database

func main () {
  http.HandleFunc("/", reqHandler)
  fmt.Println("Listening on 8000")
  db.dbInit()
  defer db.con.Close()
  http.ListenAndServe(":8000", nil)
}

func reqHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Receiving %v request from %v\n", r.Method, r.URL)

  if r.Method == "POST" {
    postHandler(w, r)
  } else if r.Method == "GET" {
    getHandler(w, r)
  }

}

func errorHandler(task string, err error) {
  fmt.Printf("Got an error trying to %s: %v\n", task, err)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path[1:] == "users" {
    getUser(w, r)
  }
}

func getUser (w http.ResponseWriter, r *http.Request) {
  queries := r.URL.Query()
  queryString := "SELECT id, password FROM users WHERE username = ?"
  if queries["username"] == nil {
    w.Write([]byte("Query incorrectly specified"))
    return
  }

  rows, err := db.con.Query(queryString, queries["username"][0])
  
  if err != nil {
    fmt.Println("querying DB", err)
  }

  defer rows.Close()
  for rows.Next() {
    var id int
    var password string
    err := rows.Scan(&id, &password)
    if err != nil {
      fmt.Println("scanning rows", err)
      w.Write([]byte("Error reading DB"))
    }
    w.Header().Set("Content-Type", "application/json")
    response, err := json.Marshal(map[string]string{"id": strconv.Itoa(id), "password": password})
    if err != nil {
      fmt.Println("error marshalling the json ", err)
    }
    w.Write(response)
  }
  w.Write([]byte("No user found"))
}