package main

import (
  "net/http"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "golang.org/x/crypto/bcrypt"
  "github.com/dgrijalva/jwt-go"
  "time"
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
  d.con = dbConnection
}

//create a global db struct for anyone to use
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

  if r.Method == "POST" {
    postHandler(w, r)
  } else if r.Method == "GET" {
    getHandler(w, r)
  }

}

func errorHandler(task string, err error) {
  fmt.Printf("Got an error trying to %s: %v\n", task, err)
}

func login(w http.ResponseWriter, r *http.Request) {
  var id int64
  var hashedPass []byte
  username := r.FormValue("usernameLI")
  password := r.FormValue("passwordLI")
  rows, err := db.con.Query("select id, password from users where username = ?", 
    username)

  if err != nil {
    errorHandler("querying DB", err)
  }

  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&id, &hashedPass)
    if err != nil {
      errorHandler("scanning rows", err)
    }
  }
  err = rows.Err()
  if err != nil {
    errorHandler("scanned rows are messed up", err)
  }
  if bcrypt.CompareHashAndPassword(hashedPass, []byte(password)) == nil {
    jwt := createJWT(id)
    sendJWT(w, jwt)
  } else {
    fmt.Println("user does not exist")
  }
}

func signup (w http.ResponseWriter, r *http.Request) {
  username := r.FormValue("usernameSU")
  password := r.FormValue("passwordSU")
  rows, err := db.con.Query("select id from users where username = ?", username)
  
  if err != nil {
    errorHandler("querying DB", err)
  }

  defer rows.Close()
  for rows.Next() {
    return
  }
  
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    errorHandler("generating hash", err)
  }

  stmt, err := db.con.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
  if err != nil {
    errorHandler("preparing to write to DB", err)
  }
  user, err := stmt.Exec(username, hashedPass)
  if err != nil {
    errorHandler("writing to DB", err)
  }
  insertId, _ := user.LastInsertId()
  jwt := createJWT(insertId)
  sendJWT(w, jwt)
}

func sendJWT (w http.ResponseWriter, jwt string) {
  c := http.Cookie{Name: "adamJWT", Value: jwt, MaxAge: 0, HttpOnly: true}
  http.SetCookie(w, &c)
  w.Write([]byte("Successful login!"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  if path == "/login" {
    login(w, r)
  } else if path == "/signup" {
    signup(w, r)
  }
}

func createJWT(userId int64) string {
  token := jwt.New(jwt.SigningMethodHS256)
  token.Claims["userId"] = userId
  token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
  tokenString, err := token.SignedString([]byte("REPLACE_WITH_A_SECRET"))
  if err != nil {
    errorHandler("signing token", err)
  }
  return tokenString
}

func getHandler(w http.ResponseWriter, r *http.Request) {
}