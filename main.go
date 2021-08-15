package main

import (
  "fmt"
  "net/http"
  "log"
  "encoding/json"
  "github.com/joho/godotenv"
)

func serveFiles(w http.ResponseWriter, r *http.Request) {
    file_name := r.URL.Path

    if file_name == "/" {
      file_name = "/index"
    }

    path := "./static" + file_name + ".html"
    http.ServeFile(w, r, path)
}

func handleRegister(w http.ResponseWriter, r *http.Request){
  var request_login UserLogin

  err := json.NewDecoder(r.Body).Decode(&request_login)

  if err != nil {
    http.Error(w, "Invalid request", 400)
    return
  }

  found, err := checkUserExists(request_login.Username, request_login.Email)

  if err != nil {
    http.Error(w, err.Error(), 500)
  }
  if found == true {
    http.Error(w, "User already exists", 400)
    return
  }

  new_user_id, err := createUser(request_login.Username, request_login.Password, request_login.Email)

  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }
  //create and return token
  token, err := createUserToken(new_user_id)

  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  token_json, err := json.Marshal(AuthRequest{token})

  if err != nil {
    http.Error(w, "Internal server error", 500)
    return
  }

  fmt.Fprintf(w, string(token_json))
}

func handleLogin(w http.ResponseWriter, r *http.Request){
  var request_user UserLogin

  err := json.NewDecoder(r.Body).Decode(&request_user)

  if err != nil {
    http.Error(w, "Invalid request", 400)
    return
  }

  user, err := getUserByUserNamePass(request_user.Username, request_user.Password)

  if err != nil {
    http.Error(w, "Invalid credentials", 403)
    return
  }

  token, err := createUserToken(user.ID.Hex())

  if err != nil {
    http.Error(w, "Error generating access token", 500)
    return
  }

  token_json, err := json.Marshal(AuthRequest{token})

  if err != nil {
    http.Error(w, "Internal server error", 500)
  }

  fmt.Fprintf(w, string(token_json))
}


func authEndpoint(w http.ResponseWriter, r *http.Request){

  user_id, err := validateAuthRequest(r.Body)
  if err != nil {
    http.Error(w, err.Error(), 403)
  }
  //perform actions with user id
  fmt.Println(user_id)
}

func main(){

  err := godotenv.Load()

  if err != nil {
    log.Fatal("Could not load .env")
  }

  port := ":8080"

  http.HandleFunc("/", serveFiles)
  http.HandleFunc("/page", serveFiles)

  http.HandleFunc("/login", handleLogin)
  http.HandleFunc("/register", handleRegister)

  http.HandleFunc("/authEndpoint", authEndpoint)

  fmt.Println("Server is running on ", port)
  log.Fatal(http.ListenAndServe(port, nil))
}
