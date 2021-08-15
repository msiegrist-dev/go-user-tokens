package main

import (
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "os"
)


func createMongoClient() (*mongo.Client, context.Context, context.CancelFunc) {
  mongo_uri := os.Getenv("MONGO_URI")
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_uri))
  if err != nil {
    log.Fatal(err)
  }
  return client, ctx, cancel
}
