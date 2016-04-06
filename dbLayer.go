package main

import (
  "net/http"
  "io/ioutil"
  // "strings"
  // "net/url"
  // "strconv"
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // "golang.org/x/crypto/bcrypt"
  // "/Users/Adam/HIR/gocode/src/github.com/dgrijalva/jwt-go"
  // "time"
  "encoding/json"
  "github.com/db-docker/myDB"
  "github.com/db-docker/models"
)

type User struct {
  Id string
  Username string
  Password string
}

// initializes the database connection
func dbInit() {
  var err error
  myDB.DBCon, err = sql.Open("mysql", "root:rodam@/gopractice")
  if err != nil {
    fmt.Println("error connecting to db", err)
    return
  }
  fmt.Println("Database connection established")
}


func main () {
  http.HandleFunc("/", reqHandler)
  fmt.Println("Listening on 8000")
  dbInit()
  defer myDB.DBCon.Close()
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
  defer r.Body.Close()   
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    errorHandler("reading request body", err)
  }

  if r.URL.Path[1:] == "user" {
    var newUser User
    err := json.Unmarshal(body, &newUser)
    if err != nil {
      fmt.Println("Error unmarshalling body", err)
    }
    id := users.AddUser(newUser.Username, newUser.Password)
    w.Write(id)
  }
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path[1:] == "users" {
    queries := r.URL.Query()
    if queries["username"] == nil {
      w.Write([]byte("Query incorrectly specified"))
      return
    }
    user := users.GetUser(queries["username"][0])

    response, err := json.Marshal(user)
    if err != nil {
      fmt.Println("error marshalling the json ", err)
    }
    fmt.Println("response is", response)
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
  }
}