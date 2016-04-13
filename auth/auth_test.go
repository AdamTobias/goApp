package auth

import(  
    "net/http"
    "net/http/httptest"
    "testing"
    "fmt"
    "strings"
)

// type MockDd struct {}

// function (db MockDb) GetBacon() {  
//     return "bacon"
// }

func TestHome(t *testing.T) {  
    // mockDb := MockDb{}
    // homeHandle := homeHandler(mockDb)
  userJson := `{"username": "dennis", "password": "tree"}`
  reader := strings.NewReader(userJson) //Convert string to reader
  // request, err := http.NewRequest("POST", usersUrl, reader) //Create request with JSON body

  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("hello moto")
    fmt.Fprintln(w, "Hello, client")
  }))
  defer ts.Close()
  // ts.URL = "http://localhost:8000"
  dbURL = ts.URL
  fmt.Println(ts.URL)
  req, _ := http.NewRequest("POST", "/signup", reader)
  w := httptest.NewRecorder()
  http.HandlerFunc(reqHandler).ServeHTTP(w, req)
  if w.Code != http.StatusOK {
      t.Errorf("Home page didn't return %v", http.StatusOK)
  }
  fmt.Println(w)
}


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