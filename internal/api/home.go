package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dineshs91/uptime/internal/db"
	"github.com/dineshs91/uptime/internal/tasks"
	"github.com/mongodb/mongo-go-driver/bson"
)

// HomeHandler - handler for root path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	dbClient := db.GetDbClient()
	collection := dbClient.Database("uptime").Collection("count")
	result := bson.NewDocument()
	filter := bson.NewDocument(bson.EC.String("uptime", "test"))
	collection.FindOne(context.Background(), filter).Decode(result)

	timeResult := result.LookupElement("time")
	fmt.Fprintf(w, "Result: %s\n", timeResult.Value().StringValue())
}

// PingHandler - handler for ping path
func PingHandler(w http.ResponseWriter, r *http.Request) {
	uptime.PingServer()
	fmt.Fprintf(w, "Pinging \n")
}
