package db

import (
	"context"
	"fmt"
	"youtube-clone/database/helpers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var FilesDb *mongo.Database
var UrlsCollection *mongo.Collection

func Connect() {
	mongoAddr := helpers.FatalIfEmptyVar("MONGODB_ADDR")
	mongoUser := helpers.FatalIfEmptyVar("MONGODB_USER")
	mongoPass := helpers.FatalIfEmptyVar("MONGODB_PASS")

	clientOptions := options.Client().ApplyURI(mongoAddr)
	clientOptions.SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPass,
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("failed to connect to mongodb")
		panic(err)
	}
	Client = c
	fmt.Println("pinging mongodb ...")
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("failed to connect to mongodb")
		panic(err)
	}
	fmt.Println("connected to mongodb")
	FilesDb = Client.Database("files")
	UrlsCollection = FilesDb.Collection("urls")
}

func Disconnect() {
	err := Client.Disconnect(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
