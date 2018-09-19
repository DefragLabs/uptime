package db

import (
	"context"
)

// AddMonitoringDetail function persists the value in db.
func AddMonitoringDetail(monitorURL MonitorURL) interface{} {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("monitorURL")

	result, _ := collection.InsertOne(
		context.Background(),
		monitorURL,
	)
	return result.InsertedID
}
