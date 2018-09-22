package db

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// GenerateObjectID function generates & adds it to passed struct.
func GenerateObjectID() objectid.ObjectID {
	return objectid.New()
}

// AddMonitoringURL function persists the value in db.
func AddMonitoringURL(monitorURL MonitorURL) interface{} {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("monitorURL")

	result, _ := collection.InsertOne(
		context.Background(),
		monitorURL,
	)
	return result.InsertedID
}

// GetMonitoringURL function gets monitor url from db.
func GetMonitoringURL() MonitorURL {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("monitorURL")

	monitorURL := MonitorURL{}
	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("protocol", "https"),
		),
	).Decode(&monitorURL)

	return monitorURL
}
