package config

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Context context.Context

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(Context, clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	Client = client
	fmt.Println("Connected to MongoDB!")
}
