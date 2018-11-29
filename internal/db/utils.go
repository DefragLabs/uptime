package db

import (
	"context"
	"fmt"
	"log"

	"github.com/defraglabs/uptime/internal/forms"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

const (
	// UsersCollection is the name of the collection which contains users.
	UsersCollection = "users"

	// ResetPasswordCollection is the name of the collection which contains data required
	// for password reset functionality.
	ResetPasswordCollection = "resetPassword"

	// MonitorURLCollection is the name of the collection which contains data of
	// monitor URL's.
	MonitorURLCollection = "monitorURL"

	// IntegrationCollection stores all the integrations configured by an user
	IntegrationCollection = "integration"
)

// GenerateObjectID generates a new objectid.
func GenerateObjectID() objectid.ObjectID {
	return objectid.New()
}

// CreateUser func persists the user to db.
func (datastore *Datastore) CreateUser(user User) interface{} {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

	result, err := collection.InsertOne(
		context.Background(),
		user,
	)
	fmt.Println(result, err)
	return result.InsertedID
}

// UpdateUser func updates user
func (datastore *Datastore) UpdateUser(user User) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

	collection.FindOneAndUpdate(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("_id", user.ID),
		),
		user,
	)
}

// GetUserByID from db.
func (datastore *Datastore) GetUserByID(userID string) User {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

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
func (datastore *Datastore) GetUserByEmail(email string) User {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

	user := User{}
	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("email", email),
		),
	).Decode(&user)

	return user
}

// GetUserByComapnyName from db.
func (datastore *Datastore) GetUserByComapnyName(companyName string) User {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

	user := User{}
	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("companyName", companyName),
		),
	).Decode(&user)

	return user
}

// AddMonitoringURL function persists the value in db.
func (datastore *Datastore) AddMonitoringURL(monitorURLForm forms.MonitorURLForm) interface{} {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	result, _ := collection.InsertOne(
		context.Background(),
		monitorURLForm,
	)
	return result.InsertedID
}

// GetMonitoringURL function gets monitor url from db.
func (datastore *Datastore) GetMonitoringURL() MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

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
func (datastore *Datastore) GetMonitoringURLS() []MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

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

// GetMonitoringURLSByUserID gets all URL's for user.
func (datastore *Datastore) GetMonitoringURLSByUserID(userID string) []MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("userID", userID),
		),
	)

	monitorURLS := make([]MonitorURL, count)
	cursor, _ := collection.Find(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("userID", userID),
		),
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
func (datastore *Datastore) AddMonitorDetail(monitorURL MonitorURL, status, time, duration string) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

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
func (datastore *Datastore) AddResetPassword(resetPassword ResetPassword) interface{} {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(ResetPasswordCollection)

	result, _ := collection.InsertOne(
		context.Background(),
		resetPassword,
	)
	return result.InsertedID
}

// GetResetPassword gets reset password record from db
func (datastore *Datastore) GetResetPassword(uid, code string) ResetPassword {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(ResetPasswordCollection)

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

// AddIntegration adds an integration to db
func (datastore *Datastore) AddIntegration(integrationForm forms.IntegrationForm) interface{} {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	collection.InsertOne(
		context.Background(),
		integrationForm,
	)

	return integrationForm
}

// GetIntegrationsByUserID gets all integrations added by an user
func (datastore *Datastore) GetIntegrationsByUserID(userID string) []Integration {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("userID", userID),
		),
	)

	cursor, _ := collection.Find(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("userID", userID),
		),
	)

	integrations := make([]Integration, count)

	i := 0
	for cursor.Next(context.Background()) {
		integration := Integration{}
		err := cursor.Decode(&integration)
		if err != nil {
			log.Fatal("error while parsing cursor for monitor urls")
		}

		integrations[i] = integration
		i++
	}

	return integrations
}

// GetIntegrationByUserID gets a specific integration added by an user
func (datastore *Datastore) GetIntegrationByUserID(userID string, integrationID string) Integration {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	integration := Integration{}

	collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("userID", userID),
			bson.EC.String("_id", integrationID),
		),
	).Decode(&integration)

	return integration
}

// DeleteIntegration delete's a given integration
func (datastore *Datastore) DeleteIntegration(userID string, integrationID string) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	collection.FindOneAndDelete(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("userID", userID),
			bson.EC.String("_id", integrationID),
		),
	)
}
