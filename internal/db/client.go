package db

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// GetDbClient - Get database client.
func GetDbClient() *mongo.Client {
	client, err := mongo.NewClient("mongodb://root:uix23wr@db:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return client
}
