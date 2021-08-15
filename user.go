package main

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "os"
  "errors"
  "fmt"
  "crypto/sha1"
)
type User struct {
  Username string `bson:"username" json:"username"`
  Password string `bson:"password" json:"password"`
  Email string `bson:"email" json:"email"`
  ID primitive.ObjectID   `json:"id" bson:"_id"`
}

func hashPassword(password string)string{
  hash := sha1.New()
  hash.Write([]byte(password))

  byte_hash := hash.Sum(nil)
  fmt.Println(string(byte_hash))
  return string(byte_hash)
}

func getUsersCollection(mongo_client *mongo.Client)(*mongo.Collection){
  DB := os.Getenv("DEV_DB")
  USERS := os.Getenv("DEV_USER")
  collection := mongo_client.Database(DB).Collection(USERS)
  return collection
}

func getUserByUserNamePass(username string, password string)(User, error){
  var user User

  mongo_client, ctx, cancel := createMongoClient()
  defer cancel()

  collection := getUsersCollection(mongo_client)
  query_err := collection.FindOne(ctx, bson.M{"username" : username, "password": hashPassword(password)}).Decode(&user)

  return user, query_err
}

func checkUserExists(username string, email string)(bool, error){
  mongo_client, ctx, cancel := createMongoClient()
  defer cancel()
  users_collection := getUsersCollection(mongo_client)

  or_query := []bson.M{}
  or_query = append(or_query, bson.M{"username": username})
  or_query = append(or_query, bson.M{"email": email})

  findQuery := bson.M{"$or": or_query}

   var found_users []User
   cursor, err := users_collection.Find(ctx, findQuery)

  if err != nil {
    return true, errors.New("Internal server error")
  }
  err = cursor.All(ctx, &found_users)

  if err != nil {
    return true, errors.New("Internal server error")
  }
  fmt.Println(found_users)
  if len(found_users) < 1{
    return false, nil
  }
  return true, nil
}

func createUser(username string, password string, email string)(string, error){
  mongo_client, ctx, cancel := createMongoClient()
  defer cancel()

  users_col := getUsersCollection(mongo_client)
  new_user, err := users_col.InsertOne(ctx, bson.M{"username": username, "password": hashPassword(password), "email": email})

  if err != nil {
    return "", errors.New("Database error")
  }

  id := new_user.InsertedID.(primitive.ObjectID)
  str := id.Hex()

  return str, nil
}
