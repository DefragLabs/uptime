package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

const (
	// Database points to the mongo database name.
	Database = "uptime"

	// UsersCollection is the name of the collection which contains users.
	UsersCollection = "users"

	// ResetPasswordCollection is the name of the collection which contains data required
	// for password reset functionality.
	ResetPasswordCollection = "resetPassword"

	// MonitorURLCollection is the name of the collection which contains data of
	// monitor URL's.
	MonitorURLCollection = "monitorURL"
)

// GenerateObjectID generates a new objectid.
func GenerateObjectID() objectid.ObjectID {
	return objectid.New()
}

// CreateUser func persists the user to db.
func CreateUser(user User) interface{} {
	dbClient := GetDbClient()
	collection := dbClient.Database(Database).Collection(UsersCollection)

	result, _ := collection.InsertOne(
		context.Background(),
		user,
	)
	return result.InsertedID
}

// UpdateUser func updates user
func UpdateUser(user User) {
	dbClient := GetDbClient()
	collection := dbClient.Database(Database).Collection(UsersCollection)

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
	collection := dbClient.Database(Database).Collection(UsersCollection)

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
	collection := dbClient.Database(Database).Collection(UsersCollection)

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
	collection := dbClient.Database(Database).Collection(MonitorURLCollection)

	result, _ := collection.InsertOne(
		context.Background(),
		monitorURL,
	)
	return result.InsertedID
}

// GetMonitoringURL function gets monitor url from db.
func GetMonitoringURL() MonitorURL {
	dbClient := GetDbClient()
	collection := dbClient.Database(Database).Collection(MonitorURLCollection)

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
	collection := dbClient.Database(Database).Collection(MonitorURLCollection)

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

// AddMonitorDetail add monitor url detail to the db.
func AddMonitorDetail(monitorURL MonitorURL, status, time string, duration string) {
	dbClient := GetDbClient()
	collection := dbClient.Database(Database).Collection(MonitorURLCollection)

	fmt.Println(duration)

	result := MonitorResult{
		Status:   status,
		Duration: duration,
		Time:     time,
	}

	monitorURL.Results = append(monitorURL.Results, result)

	collection.FindOneAndReplace(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", monitorURL.ID),
		),
		monitorURL,
	)
}

// AddResetPassword adds password code with the user id.
func AddResetPassword(resetPassword ResetPassword) interface{} {
	dbClient := GetDbClient()
	collection := dbClient.Database(Database).Collection(ResetPasswordCollection)

	result, _ := collection.InsertOne(
		context.Background(),
		resetPassword,
	)
	return result.InsertedID
}

// GetResetPassword gets reset password record from db
func GetResetPassword(uid, code string) ResetPassword {
	dbClient := GetDbClient()
	collection := dbClient.Database(Database).Collection(ResetPasswordCollection)

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
