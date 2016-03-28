package main

import (
  "net/http"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type database struct {
  con *sql.DB
}

// user:password@/dbname
func (d *database) dbInit() {
  // db, err := sql.Open("mysql", "root:rodam@tcp(127.0.0.1:5587)/gopractice")
  dbConnection, err := sql.Open("mysql", "root:rodam@/gopractice")
  if err != nil {
    fmt.Println("error connecting to db")
    return
  }
  d.con = dbConnection
}

var db database

func main () {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.HandleFunc("/", reqHandler)
  fmt.Println("Listening on 8080")
  db.dbInit()
  defer db.con.Close()
  http.ListenAndServe(":8080", nil)
}

func reqHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Receiving %v request from %v\n", r.Method, r.URL)
  // stmt, err := db.con.Prepare("INSERT INTO users(username) VALUES(?)")
  // if err != nil {
  //   fmt.Printf("error preparing thing %v", err)
  // }
  // _, err = stmt.Exec("Dolly")
  // if err != nil {
  //   fmt.Println("error executing thing")
  // }

  if r.Method == "POST" {
    postHandler(w, r)
  } else if r.Method == "GET" {
    getHandler(w, r)
  }

}

func postHandler(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  if path == "/login" {
    //login stuff:
    //check if user + pw is in DB
      //if so, make them a JWT...
      //if not, tell them they failed

    var id int
    rows, err := db.con.Query("select id from users where username = ? and password = ?", 
      r.FormValue("usernameLI"), r.FormValue("passwordLI"))
    
    fmt.Printf("username = %s AND password = %s", r.FormValue("usernameLI"), r.FormValue("passwordLI"))

    if err != nil {
      fmt.Println("error! 14")
    }

    defer rows.Close()
    for rows.Next() {
      fmt.Println("flag 12")
      err := rows.Scan(&id)
      if err != nil {
        fmt.Println("error! 15")
      }
    }
    err = rows.Err()
    if err != nil {
      fmt.Println("error! 16")
    }
    if id != 0 {
      fmt.Println("user exists")
      //TODO: user exists stuff (JWT ?)
    } else {
      fmt.Println("user does not exist")
      //TODO: user does not exist stuff (error msg ?)
      //should I be using AJAX here for the POST instead of non AJAX HTTP?
    }

  } else if path == "/signup" {
    //signup stuff
    stmt, err := db.con.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
    if err != nil {
      fmt.Printf("error preparing thing %v", err)
    }
    _, err = stmt.Exec(r.FormValue("usernameSU"), r.FormValue("passwordSU"))
    if err != nil {
      //this currently fires when the username is already taken -- maybe we should do a lookup first to verify
      //the username is available.  Currently I can't tell the difference between "the execution failed" and 
      //"username taken"
      fmt.Println("error executing thing")
      //when user is rejected, should get some feedback and no redirection.
      //should likely be using AJAX request here
    }
    //need to do something when a user is authenticated via signup -- JWT and redirect?
  }
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  
}

// func dbInit () {
//   db, err := sql.Open("mysql", "root:rodam@tcp(127.0.0.1:5587)/users")
//   if err != nil {
//     fmt.Println("error connecting to db")
//     return
//   }
//   dbCon(db)
// }

// func dbCon (db... *sql.DB) {
//   if len(db) == 1 {

//   }
// }