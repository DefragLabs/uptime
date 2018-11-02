package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var mongoClient *mongo.Client

// Datastore is a wrapper around mongo.Client.
type Datastore struct {
	Client       *mongo.Client
	DatabaseName string
}

// New creates a new mongo.Client
func New() *Datastore {
	databaseName, ok := os.LookupEnv("MONGO_DATABASE_NAME")
	if !ok {
		databaseName = "uptime"
	}

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

	return &Datastore{
		Client:       mongoClient,
		DatabaseName: databaseName,
	}
}
