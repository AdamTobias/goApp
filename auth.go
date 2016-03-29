package main

import (
  "net/http"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "golang.org/x/crypto/bcrypt"
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
    var hashedPass []byte
    username := r.FormValue("usernameLI")
    password := r.FormValue("passwordLI")
    rows, err := db.con.Query("select id, password from users where username = ?", 
      username)
  
    if err != nil {
      fmt.Println("error! 14")
    }

    defer rows.Close()
    for rows.Next() {
      fmt.Println("flag 12")
      err := rows.Scan(&id, &hashedPass)
      if err != nil {
        fmt.Println("error! 15")
      }
    }
    err = rows.Err()
    if err != nil {
      fmt.Println("error! 16")
    }
    if bcrypt.CompareHashAndPassword(hashedPass, []byte(password)) == nil {
      fmt.Printf("user exists: %d\n", id)
      //TODO: user exists stuff (JWT ?)
    } else {
      fmt.Println("user does not exist")
      //TODO: user does not exist stuff (error msg ?)
      //should I be using AJAX here for the POST instead of non AJAX HTTP?
    }

  } else if path == "/signup" {
    //signup stuff
    username := r.FormValue("usernameSU")
    password := r.FormValue("passwordSU")
    rows, err := db.con.Query("select id from users where username = ?", username)
    
    if err != nil {
      fmt.Println("error! 114")
    }

    defer rows.Close()
    for rows.Next() {
      // USERNAME TAKEN
      fmt.Fprint(w, "username taken")
      return
    }
    
    // USERNAME NOT TAKEN

    hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
      fmt.Println("error hashing password")
    }

    stmt, err := db.con.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
    if err != nil {
      fmt.Printf("error preparing thing %v", err)
    }
    _, err = stmt.Exec(username, hashedPass)
    if err != nil {
      fmt.Println("error executing thing")
    }
    //need to do something when a user is authenticated via signup -- JWT and redirect?
    http.Redirect(w, r, "/", 301)
    
  }
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  
}