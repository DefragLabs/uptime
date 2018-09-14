package uptime

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/dineshs91/uptime/internal/db"
)

// PingServer function pings any server, and stores the status code.
func PingServer() {
	var frequency = time.Duration(100)
	ticker := time.Tick(frequency)

	dbClient := db.GetDbClient()
	collection := dbClient.Database("uptime").Collection("count")

	go func() {
		for {
			select {
			case t := <-ticker:
				time.Sleep(10 * time.Second)
				result := bson.NewDocument()
				err := collection.FindOne(
					context.Background(),
					bson.NewDocument(
						bson.EC.String("uptime", "test"),
					),
				).Decode(result)

				if err == nil {
					collection.FindOneAndUpdate(
						context.Background(),
						bson.NewDocument(
							bson.EC.String("uptime", "test"),
						),
						bson.NewDocument(
							bson.EC.SubDocumentFromElements(
								"$set",
								bson.EC.String("time", time.Now().Format(time.UnixDate)),
							),
						),
					)
				} else {
					collection.InsertOne(
						context.Background(),
						bson.NewDocument(
							bson.EC.String("uptime", "test"),
							bson.EC.String("time", t.Format(time.UnixDate)),
						),
					)
				}
			}
		}
	}()
}
