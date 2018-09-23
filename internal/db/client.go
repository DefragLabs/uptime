package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// GetDbClient - Get database client.
func GetDbClient() *mongo.Client {
	mongoPass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	connString := fmt.Sprintf("mongodb://root:%s@db:27017", mongoPass)
	client, err := mongo.NewClient(connString)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return client
}
