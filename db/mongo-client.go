package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log/slog"
	"os"
)

var MongoDb *mongo.Database

func init() {
	uri := os.Getenv("MONGO_CREDENTIALS")
	MongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	databaseName := os.Getenv("MONGO_DATABASE")
	database := MongoClient.Database(databaseName)
	if database.Name() == "" {
		slog.Error("Error connecting to mongo database", slog.String("database", databaseName))
		os.Exit(1)
	}

	MongoDb = database
}
