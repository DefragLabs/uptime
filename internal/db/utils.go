package db

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// GenerateObjectID function generates & adds it to passed struct.
func GenerateObjectID() objectid.ObjectID {
	return objectid.New()
}

// CreateUser persists the user to db.
func CreateUser(user User) interface{} {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("users")

	result, _ := collection.InsertOne(
		context.Background(),
		user,
	)
	return result.InsertedID
}

// UpdateUser updates user
func UpdateUser(user User) {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("users")

	collection.FindOneAndUpdate(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", user.ID),
		),
		user,
	)
}

// GetUserByID from db.
func GetUserByID(userID string) User {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("users")

	user := User{}
	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", userID),
		),
	).Decode(&user)

	return user
}

// GetUserByEmail from db.
func GetUserByEmail(email string) User {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("users")

	user := User{}
	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("email", email),
		),
	).Decode(&user)

	return user
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

// GetMonitoringURLS  gets all added url's
func GetMonitoringURLS() []MonitorURL {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("monitorURL")

	count, _ := collection.Count(
		context.Background(),
		bson.NewDocument(),
	)

	monitorURLS := make([]MonitorURL, count)
	cursor, _ := collection.Find(
		context.Background(),
		bson.NewDocument(),
	)

	i := 0
	for cursor.Next(context.Background()) {
		monitorURL := MonitorURL{}
		err := cursor.Decode(&monitorURL)
		if err != nil {
			log.Fatal("error while parsing cursor for monitor urls")
		}

		monitorURLS[i] = monitorURL
		i++
	}

	return monitorURLS
}

// AddResetPassword adds password code with the user id.
func AddResetPassword(resetPassword ResetPassword) interface{} {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("resetPassword")

	result, _ := collection.InsertOne(
		context.Background(),
		resetPassword,
	)
	return result.InsertedID
}

// GetResetPassword gets reset password record from db
func GetResetPassword(uid, code string) ResetPassword {
	dbClient := GetDbClient()
	collection := dbClient.Database("uptime").Collection("resetPassword")

	resetPassword := ResetPassword{}
	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("uid", uid),
			bson.EC.String("code", code),
		),
	).Decode(&resetPassword)

	return resetPassword
}
