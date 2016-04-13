package auth

import (
  "net/http"
  "fmt"
  "golang.org/x/crypto/bcrypt"
  "github.com/dgrijalva/jwt-go"
  "time"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "os"
)

type User struct {
  Id string
  Username string
  Password string
}

var dbURL = "http://localhost:8000"

func main () {
  envKey := "DB_NAME"
  envVal, _ := os.LookupEnv(envKey) 
  fmt.Println(envKey, envVal)
  envKey = "DB_PORT"
  envVal, _ = os.LookupEnv(envKey)
  fmt.Println(envKey, envVal)
  if env, _ := os.LookupEnv("env"); env == "docker" {
    dbURL = "HTTP" + envVal[3:]
    fmt.Println("env =", env)
  } else {
    dbURL = "http://localhost:8000"
    fmt.Println("env = ", env)
  }

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.HandleFunc("/", reqHandler)
  fmt.Println("Listening on 8080")
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
  if err != nil {
    fmt.Printf("Got an error trying to %s: %v\n", task, err)
  }
}

func login(w http.ResponseWriter, r *http.Request) {
  username := r.FormValue("usernameLI")
  password := r.FormValue("passwordLI")

  resp, err := http.Get(dbURL + "/users?username=" + username)
  errorHandler("making GET request", err)
  
  defer resp.Body.Close()   
  body, err2 := ioutil.ReadAll(resp.Body)
  errorHandler("reading request body", err2)
  
  var retrievedUser User
  err = json.Unmarshal(body, &retrievedUser)
  errorHandler("Unmarshaling response body", err)

  if retrievedUser.Username != "" {
    if bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(password)) == nil {
      jwt := createJWT(retrievedUser.Id)
      sendJWT(w, jwt)
      return
    } 
  } 
  w.Write([]byte("wrong bloke, mate"))

}

func signup (w http.ResponseWriter, r *http.Request) {
  username := r.FormValue("usernameSU")
  password := r.FormValue("passwordSU")
  
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  errorHandler("generating hash", err)

  body, err := json.Marshal(map[string]string{"Username":username, "Password":string(hashedPass)})
  errorHandler("Marshaling body", err)

  resp, err := http.Post(dbURL + "/users", "application/json", bytes.NewReader(body))
  errorHandler("making POST request", err)

  defer resp.Body.Close()   
  insertId, err2 := ioutil.ReadAll(resp.Body)
  errorHandler("reading request body", err2)
  
  id := string(insertId)
  if id != "" {
    jwt := createJWT(id)
    sendJWT(w, jwt) 
  } else {
    w.Write([]byte("username taken mate"))
  }
}

func sendJWT (w http.ResponseWriter, jwt string) {
  c := http.Cookie{Name: "adamJWT", Value: jwt, MaxAge: 0, HttpOnly: true}
  http.SetCookie(w, &c)
  w.Write([]byte("Successful login!"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  if path == "/login" { //login should be a get
    login(w, r)
  } else if path == "/signup" {
    signup(w, r)
  }
}

func createJWT(userId string) string {
  token := jwt.New(jwt.SigningMethodHS256)
  token.Claims["userId"] = userId
  token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
  tokenString, err := token.SignedString([]byte("REPLACE_WITH_A_SECRET"))
  errorHandler("signing token", err)
  
  return tokenString
}

func getHandler(w http.ResponseWriter, r *http.Request) {
}