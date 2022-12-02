package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database
var Users *mongo.Collection
var Movies *mongo.Collection
var Genres *mongo.Collection
var Lists *mongo.Collection
var Reviews *mongo.Collection
var Invitations *mongo.Collection

func Connect() {
	MONGO_URI := "mongodb+srv://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@cluster0.iflqo.mongodb.net/?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(MONGO_URI)

	Client, err := mongo.Connect(context.TODO(), clientOptions)
	handleDbError(err)

	DB = Client.Database("Cinephile")
	Users = DB.Collection("Users")
	Movies = DB.Collection("Movies")
	Genres = DB.Collection("Genres")
	Lists = DB.Collection("Lists")
	Reviews = DB.Collection("Reviews")
	Invitations = DB.Collection("Invitations")
}

func handleDbError(err error) {
	if err != nil {
		log.Fatalf("DB connection failed. Err: %s", err)
	}
}
