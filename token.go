package main

import (
  "encoding/json"
  "os"
  "fmt"
  "github.com/golang-jwt/jwt"
  "time"
  "errors"
  "io"
)

type UserLogin struct {
  Username string `json:"username"`
  Password string `json:"password"`
  Email string `json:"email"`
}

type AuthRequest struct {
  Token string `json:"token"`
}

func createUserToken(id string)(string, error){
  token_secret := os.Getenv("TOKEN_SECRET")
  claims := jwt.MapClaims{}
  claims["id"] = id
  claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
  token_str := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id": id,
    "exp": time.Now().Add(time.Minute * 30).Unix(),
  })
  token, err := token_str.SignedString([]byte(token_secret))
  return token, err
}

func isTokenExpired(expiry float64)bool{
  now := time.Now().Unix()
  int_exp := int64(expiry)
  if now < int_exp{
    return false
  }
  return true
}

func decodeToken(token_str string)(string, error){
  token_secret := os.Getenv("TOKEN_SECRET")
  token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {

     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, errors.New("Invalid access token")
     }
     return []byte(token_secret), nil
  })

  if err != nil {
    return "", err
  }

  claims, ok := token.Claims.(jwt.MapClaims)
  if ok != true {
    return "", errors.New("Invalid access token")
  }

  expiry := claims["exp"].(float64)
  expired := isTokenExpired(expiry)
  if expired == true {
    return "", errors.New("Token is expired")
  }

  id := fmt.Sprintf("%v", claims["id"])
  return id, nil
}

func validateAuthRequest(body io.Reader)(string, error){
  decoder := json.NewDecoder(body)
  var request_body AuthRequest

  err := decoder.Decode(&request_body)
  if err != nil {
    return "", errors.New("Bad request")
  }
  user_id, err := decodeToken(request_body.Token)
  if err != nil {
    return "", errors.New("Invalid access token")
  }
  //perform actions with user id
  return user_id, nil

}
