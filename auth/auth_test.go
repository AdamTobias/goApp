package main

import(  
  "bytes"
  "net/http"
  "net/http/httptest"
  "testing"
  "fmt"
  "io/ioutil"
  "net/url"
  "encoding/json"
  "os"
)

// type MockDd struct {}

// function (db MockDb) GetBacon() {  
//     return "bacon"
// }
var ts *httptest.Server

// type User struct {
//   Id string
//   Username string
//   Password string
// }

func TestMain(m *testing.M) {
  ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" && r.URL.Path[1:] == "users" {
      var response []byte
      var returnedUser User
      testUsername := r.URL.Query()["username"][0]
      if testUsername == "shouldPass" {
        returnedUser = User{"3", "shouldPass", "$2a$10$GMHH7c5Bo3tZ3oWhvU3bZ.1ft56Vs9j0QTnhEh2i1bkV1h6UrBKA2"}
      } else if testUsername == "shouldFail" {
        returnedUser = User{"4", "shouldFail", "$2a$10$7LS5ZPi3q57.Fy1niqfhkeWnN19pxw4/wCd3yC2SnsU8LqqZb9yue"}
      } else if testUsername == "DNE" {
        returnedUser = User{"", "", ""}
      }
      response, err := json.Marshal(returnedUser)
      if err != nil {
        fmt.Println("error marshalling the json ", err)
      }
      w.Header().Set("Content-Type", "application/json")
      w.Write(response)
    } else if r.Method == "POST" && r.URL.Path[1:] == "users" {
      var returnedId string
      defer r.Body.Close()   
      body, err := ioutil.ReadAll(r.Body)
      if err != nil {
        errorHandler("reading request body", err)
      }
      var newUser User
      err = json.Unmarshal(body, &newUser)
      if err != nil {
        fmt.Println("Error unmarshalling body", err)
      }
      if newUser.Username == "shouldPass" {
        returnedId = "3"
      } else if newUser.Username == "shouldFail" {
        returnedId = ""
      }
      w.Write([]byte(returnedId))
    }
  }))
  defer ts.Close()
  dbURL = ts.URL

  os.Exit(m.Run())
}

func TestLogin(t *testing.T) {  
  if submitUser("shouldPass", "correctPassword", "login") != http.StatusOK {
    t.Errorf("Correct login didn't return %v", http.StatusUnauthorized)
  }

  if submitUser("shouldFail", "incorrectPassword", "login") != http.StatusUnauthorized {
    t.Errorf("Incorrect password didn't return %v", http.StatusUnauthorized)
  }

  if submitUser("DNE", "irrelevant", "login") != http.StatusUnauthorized {
    t.Errorf("Unknown user didn't return %v", http.StatusUnauthorized)
  }
}

func TestSignup(t * testing.T) {
  if submitUser("shouldPass", "irrelevant", "signup") != http.StatusOK {
    t.Errorf("Valid signup didn't return %v", http.StatusOK)
  } else if submitUser("shouldFail", "irrelevant", "signup") != http.StatusBadRequest {
    t.Errorf("Invalid signup didn't return %v", http.StatusBadRequest)
  }
}

func submitUser(username string, password string, mode string) int {
  var fieldTag string
  if mode == "login" {
    fieldTag = "LI"
  } else if mode == "signup" {
    fieldTag = "SU"
  }
  data := url.Values{}
  data.Set("username" + fieldTag, username)
  data.Set("password" + fieldTag, password)
  req, _ := http.NewRequest("POST", "/" + mode, bytes.NewBufferString(data.Encode()))
  w := httptest.NewRecorder()
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  http.HandlerFunc(reqHandler).ServeHTTP(w, req)
  
  return w.Code
}

    // data := url.Values{}
    // data.Set("name", "foo")
    // data.Add("surname", "bar")

    // u, _ := url.ParseRequestURI(apiUrl)
    // u.Path = resource
    // urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

    // client := &http.Client{}
    // r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
    // r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    // r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

// func TestLogin(t *testing.T) {
//   userDB := `{"username": "dennis", "password": "tree", "id": "4"}`  
//   userJson := `{"username": "dennis", "password": "tree", "id": "4"}`
//   userReader := strings.NewReader(userJson) //Convert string to reader

//   ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     response, err := json.Marshal(user)
//     w.Header().Set("Content-Type", "application/json")
//     w.Write(response)
//   }))
//   defer ts.Close()
//   dbURL = ts.URL
//   req, _ := http.NewRequest("POST", "/login", userReader)
//   w := httptest.NewRecorder()
//   http.HandlerFunc(reqHandler).ServeHTTP(w, req)
//   if w.Code != http.StatusOK {
//       t.Errorf("Home page didn't return %v", http.StatusOK)
//   }
//   fmt.Println(w)
// }


// func main() {

//   res, err := http.Get(ts.URL)
//   if err != nil {
//     log.Fatal(err)
//   }
//   greeting, err := ioutil.ReadAll(res.Body)
//   res.Body.Close()
//   if err != nil {
//     log.Fatal(err)
//   }

//   fmt.Printf("%s", greeting)
// }

// import (
//     "fmt"
//     "io"
//     "net/http"
//     "net/http/httptest"
//     "strings"
//     "testing"
// )

// var (
//     server   *httptest.Server
//     reader   io.Reader
//     usersUrl string
// )

// func init() {

//   server = httptest.NewServer(http.HandleFunc("/", reqHandler)) //Creating new server with the user handlers

//   usersUrl = fmt.Sprintf("%s/users", server.URL) //Grab the address for the API endpoint
//   fmt.Println(usersUrl)
// }

// func TestCreateUser(t *testing.T) {
//     userJson := `{"username": "dennis", "balance": 200}`

//     reader = strings.NewReader(userJson) //Convert string to reader

//     request, err := http.NewRequest("POST", usersUrl, reader) //Create request with JSON body

//     res, err := http.DefaultClient.Do(request)

//     if err != nil {
//         t.Error(err) //Something is wrong while sending request
//     }

//     if res.StatusCode != 201 {
//         t.Errorf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
//     }
// }