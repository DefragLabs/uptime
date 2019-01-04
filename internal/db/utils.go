package db

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"

	"github.com/mongodb/mongo-go-driver/options"

	"github.com/defraglabs/uptime/internal/forms"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	log "github.com/sirupsen/logrus"
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

	// MonitorResultCollection stores the result of the pings to the configured monitoring url's
	MonitorResultCollection = "monitorURLResult"

	// IntegrationCollection stores all the integrations configured by an user
	IntegrationCollection = "integration"
)

// AddIndexes adds mongo indexes.
func (datastore *Datastore) AddIndexes() {
	dbClient := datastore.Client

	addIndexesOnMonitorResultCollection(dbClient, datastore)
	addTextIndexesOnMonitorURLCollection(dbClient, datastore)

	log.Info("Added db indexes")
}

func addIndexesOnMonitorResultCollection(dbClient *mongo.Client, datastore *Datastore) {
	monitorResultCollection := dbClient.Database(datastore.DatabaseName).Collection(MonitorResultCollection)

	indexes := monitorResultCollection.Indexes()
	indexes.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{{"time", bsonx.Int32(-1)}},
		},
	)
}

func addTextIndexesOnMonitorURLCollection(dbClient *mongo.Client, datastore *Datastore) {
	monitorURLCollection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	indexes := monitorURLCollection.Indexes()
	indexes.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{"url", bsonx.String("text")},
				{"name", bsonx.String("text")},
			},
		},
	)
}

// GenerateObjectID generates a new objectid.
func GenerateObjectID() objectid.ObjectID {
	return objectid.New()
}

// CreateUser func persists the user to db.
func (datastore *Datastore) CreateUser(user User) interface{} {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

	result, _ := collection.InsertOne(
		context.Background(),
		user,
	)

	return result.InsertedID
}

// UpdateUser func updates user
func (datastore *Datastore) UpdateUser(user User) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(UsersCollection)

	collection.FindOneAndUpdate(
		context.Background(),
		bson.D{
			{"_id", user.ID},
		},
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
		bson.D{
			{"_id", userID},
		},
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
		bson.D{
			{"email", email},
		},
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
		bson.D{
			{"companyName", companyName},
		},
	).Decode(&user)

	return user
}

// AddMonitoringURL function persists the value in db.
func (datastore *Datastore) AddMonitoringURL(monitorURLForm forms.MonitorURLForm) MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	_, err := collection.InsertOne(
		context.Background(),
		monitorURLForm,
	)

	var monitorURL MonitorURL
	if err != nil {
		monitorURL = MonitorURL{}
	} else {
		monitorURL = MonitorURL{
			ID:        monitorURLForm.ID,
			UserID:    monitorURLForm.UserID,
			Protocol:  monitorURLForm.Protocol,
			URL:       monitorURLForm.URL,
			Frequency: monitorURLForm.Frequency,
			Unit:      monitorURLForm.Unit,
		}
	}

	return monitorURL
}

// GetMonitoringURLS  gets all added url's
func (datastore *Datastore) GetMonitoringURLS() []MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{},
	)

	monitorURLS := make([]MonitorURL, count)
	cursor, _ := collection.Find(
		context.Background(),
		bson.D{},
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

// GetMonitoringURLByUserIDAndURL filters with userID, protocol & URL.
func (datastore *Datastore) GetMonitoringURLByUserIDAndURL(userID, protocol, URL string) MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	monitorURL := MonitorURL{}
	collection.FindOne(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"protocol", protocol},
			{"url", URL},
		},
	).Decode(&monitorURL)

	return monitorURL
}

// GetMonitoringURLSByUserIDCount returns the count of monitoring URL's for user.
func (datastore *Datastore) GetMonitoringURLSByUserIDCount(userID string) int64 {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"userID", userID},
		},
	)

	return count
}

// GetMonitoringURLSByUserIDAndStatus returns the monitoring URL's for user based on the status.
func (datastore *Datastore) GetMonitoringURLSByUserIDAndStatus(userID, status string) int64 {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"status", status},
		},
	)

	return count
}

// GetMonitoringURLSByUserID gets all URL's for user.
func (datastore *Datastore) GetMonitoringURLSByUserID(userID string) []MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"userID", userID},
		},
	)

	monitorURLS := make([]MonitorURL, count)
	cursor, _ := collection.Find(
		context.Background(),
		bson.D{
			{"userID", userID},
		},
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

// GetMonitoringURLByUserID gets monitor URL by userID & monitoringURLID
func (datastore *Datastore) GetMonitoringURLByUserID(userID, monitoringURLID string) MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	monitorURL := MonitorURL{}
	collection.FindOne(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"_id", monitoringURLID},
		},
	).Decode(&monitorURL)

	return monitorURL
}

// UpdateMonitoringURLByUserID updates monitor URL
func (datastore *Datastore) UpdateMonitoringURLByUserID(userID, monitoringURLID string, monitorURLForm forms.MonitorURLForm) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	collection.FindOneAndUpdate(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"_id", monitoringURLID},
		},
		bson.D{
			{"$set", bson.D{
				{"protocol", monitorURLForm.Protocol},
				{"frequency", monitorURLForm.Frequency},
				{"unit", monitorURLForm.Unit},
			}},
		},
	)
}

// DeleteMonitoringURL delete's the provided monitorURL
func (datastore *Datastore) DeleteMonitoringURL(userID, monitoringURLID string) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	collection.FindOneAndDelete(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"_id", monitoringURLID},
		},
	)
}

// SearchMonitoringURL searches text fields
func (datastore *Datastore) SearchMonitoringURL(userID, searchText string) []MonitorURL {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"$text", bson.D{
				{"$search", searchText},
			}},
		},
	)

	monitorURLS := make([]MonitorURL, count)
	cursor, _ := collection.Find(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"$text", bson.D{
				{"$search", searchText},
			}},
		},
	)

	i := 0
	for cursor.Next(context.Background()) {
		monitorURL := MonitorURL{}
		err := cursor.Decode(&monitorURL)
		if err != nil {
			log.Info("error while parsing cursor for monitor urls:", err)
		}

		monitorURLS[i] = monitorURL
		i++
	}

	return monitorURLS
}

// AddMonitorDetail add monitor url detail to the db.
// Status UP/DOWN, statusCode is the http response code
func (datastore *Datastore) AddMonitorDetail(monitorURL MonitorURL, statusCode, status, time string, responseTime float64) MonitorResult {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorURLCollection)

	objectID := GenerateObjectID()

	result := MonitorResult{
		ID:                objectID.Hex(),
		MonitorURLID:      monitorURL.ID,
		Status:            status,
		StatusDescription: statusCode,
		ResponseTime:      responseTime,
		Time:              time,
	}

	// Update status in monitorURL
	monitorURL.Status = status
	collection.FindOneAndReplace(
		context.Background(),
		bson.D{
			{"_id", monitorURL.ID},
		},
		monitorURL,
	)

	monitorResultCollection := dbClient.Database(datastore.DatabaseName).Collection(MonitorResultCollection)
	monitorResultCollection.InsertOne(
		context.Background(),
		result,
	)

	return result
}

// GetMonitoringURLStats gets the stats for given monitorURLID
func (datastore *Datastore) GetMonitoringURLStats(monitorURLID string) []MonitorResult {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorResultCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"monitorURLID", monitorURLID},
		},
	)

	cursor, _ := collection.Find(
		context.Background(),
		bson.D{
			{"monitorURLID", monitorURLID},
		},
	)

	monitorResults := make([]MonitorResult, count)

	i := 0
	for cursor.Next(context.Background()) {
		monitorResult := MonitorResult{}
		err := cursor.Decode(&monitorResult)

		if err != nil {
			log.Fatal("error while parsing cursor for monitor urls result")
		}

		monitorResults[i] = monitorResult
		i++
	}

	return monitorResults
}

// GetMonitoringURLStatsInInterval gets the stats for given monitorURLID in the given interval.
// Example:
//     Get url stats for monitorURL in the last 1 hour.
func (datastore *Datastore) GetMonitoringURLStatsInInterval(monitorURLID string, value int32, unit string) []MonitorResult {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorResultCollection)

	now := time.Now()
	currentTime := now.UTC().String()

	var fromTime string
	if unit == "hour" {
		fromTime = now.Add(-(time.Duration(value) * time.Hour)).UTC().String()
	} else if unit == "day" {
		fromTime = now.Add(-(time.Duration(value) * 24 * time.Hour)).UTC().String()
	} else if unit == "week" {
		fromTime = now.Add(-(time.Duration(value) * 24 * 7 * time.Hour)).UTC().String()
	} else if unit == "month" {
		fromTime = now.Add(-(time.Duration(value) * 24 * 31 * time.Hour)).UTC().String()
	}

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"monitorURLID", monitorURLID},
			{"time", bson.D{
				{"$gte", fromTime},
				{"$lte", currentTime},
			}},
		},
	)

	cursor, _ := collection.Find(
		context.Background(),
		bson.D{
			{"monitorURLID", monitorURLID},
			{"time", bson.D{
				{"$gte", fromTime},
				{"$lte", currentTime},
			}},
		},
	)

	monitorResults := make([]MonitorResult, count)

	i := 0
	for cursor.Next(context.Background()) {
		monitorResult := MonitorResult{}
		err := cursor.Decode(&monitorResult)

		if err != nil {
			log.Fatal("error while parsing cursor for monitor urls result")
		}

		monitorResults[i] = monitorResult
		i++
	}

	return monitorResults
}

// GetLastNMonitoringURLStats gets the stats for given monitorURLID
func (datastore *Datastore) GetLastNMonitoringURLStats(monitorURLID string, n int64) []MonitorResult {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(MonitorResultCollection)

	findOptions := options.Find()
	findOptions.Sort = bson.D{
		{"time", -1},
	}
	findOptions.Limit = &n

	cursor, _ := collection.Find(
		context.Background(),
		bson.D{
			{"monitorURLID", monitorURLID},
		},
		findOptions,
	)

	monitorResults := make([]MonitorResult, n)

	i := 0
	for cursor.Next(context.Background()) {
		monitorResult := MonitorResult{}
		err := cursor.Decode(&monitorResult)

		if err != nil {
			log.Fatal("error while parsing cursor for monitor urls result")
		}

		monitorResults[i] = monitorResult
		i++
	}

	return monitorResults
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
		bson.D{
			{"uid", uid},
			{"code", code},
		},
	).Decode(&resetPassword)

	return resetPassword
}

// AddIntegration adds an integration to db
func (datastore *Datastore) AddIntegration(integrationForm forms.IntegrationForm) Integration {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	collection.InsertOne(
		context.Background(),
		integrationForm,
	)

	integration := Integration{
		ID:         integrationForm.ID,
		UserID:     integrationForm.UserID,
		Email:      integrationForm.Email,
		Type:       integrationForm.Type,
		WebhookURL: integrationForm.WebhookURL,
	}

	return integration
}

// GetIntegrations gets all integrations.
func (datastore *Datastore) GetIntegrations() []Integration {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{},
	)

	cursor, _ := collection.Find(
		context.Background(),
		bson.D{},
	)

	integrations := make([]Integration, count)

	i := 0
	for cursor.Next(context.Background()) {
		integration := Integration{}
		err := cursor.Decode(&integration)
		if err != nil {
			log.Fatal("error while parsing cursor for integrations")
		}

		integrations[i] = integration
		i++
	}

	return integrations
}

// GetIntegrationsByUserID gets all integrations added by an user
func (datastore *Datastore) GetIntegrationsByUserID(userID string) []Integration {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	count, _ := collection.Count(
		context.Background(),
		bson.D{
			{"userID", userID},
		},
	)

	cursor, _ := collection.Find(
		context.Background(),
		bson.D{
			{"userID", userID},
		},
	)

	integrations := make([]Integration, count)

	i := 0
	for cursor.Next(context.Background()) {
		integration := Integration{}
		err := cursor.Decode(&integration)
		if err != nil {
			log.Fatal("error while parsing cursor for integrations")
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
		bson.D{
			{"userID", userID},
			{"_id", integrationID},
		},
	).Decode(&integration)

	return integration
}

// DeleteIntegration delete's a given integration
func (datastore *Datastore) DeleteIntegration(userID string, integrationID string) {
	dbClient := datastore.Client
	collection := dbClient.Database(datastore.DatabaseName).Collection(IntegrationCollection)

	collection.FindOneAndDelete(
		context.Background(),
		bson.D{
			{"userID", userID},
			{"_id", integrationID},
		},
	)
}
