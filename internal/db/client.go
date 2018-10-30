package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var mongoClient *mongo.Client

// datastore is a wrapper around mongo.Client.
type datastore struct {
	client *mongo.Client
}

// New creates a new mongo.Client
func New() *datastore {
	mongoPass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	connString := fmt.Sprintf("mongodb://root:%s@db:27017", mongoPass)
	if mongoClient == nil {
		client, err := mongo.NewClient(connString)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		mongoClient = client
	}

	return &datastore{
		client: mongoClient,
	}
}
