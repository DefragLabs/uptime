package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/defraglabs/uptime/internal/utils"
)

func createTestUser() (db.User, string) {
	userRegisterForm := forms.UserRegisterForm{
		FirstName:   "Alice",
		LastName:    "Wonderland",
		Email:       "alice@sample.com",
		Password:    "test@123",
		CompanyName: "skynet",
	}
	newUser := db.RegisterUser(userRegisterForm)
	objectID := db.GenerateObjectID()
	newUser.ID = objectID.Hex()

	datastore := db.New()
	datastore.CreateUser(newUser)
	jwt, _ := db.GetJWT(newUser, userRegisterForm.Password)
	return newUser, jwt
}

// Clears monitor collection. Also clears users collection.
// We create test user to authenticate the requests. we clear them after
// every test.
func clearMonitorCollection() {
	datastore := db.New()

	datastore.Client.Database(datastore.DatabaseName).Collection(
		db.UsersCollection).Drop(context.Background())

	datastore.Client.Database(datastore.DatabaseName).Collection(
		db.MonitorURLCollection).Drop(context.Background())
}

func addTestMonitorURL(userID string) string {
	monitorURLForm := forms.MonitorURLForm{
		UserID:    userID,
		Protocol:  "http",
		URL:       "example.com",
		Frequency: 5,
		Unit:      "minute",
	}
	objectID := db.GenerateObjectID()
	monitorURLForm.ID = objectID.Hex()

	datastore := db.New()
	monitoringURL := datastore.AddMonitoringURL(monitorURLForm)

	return monitoringURL.ID
}

func addTestMonitorURLResult(userID, monitorURLID string) string {
	datastore := db.New()

	monitorURL := datastore.GetMonitoringURLByUserID(userID, monitorURLID)

	status := utils.GetServiceStatus(http.StatusOK)
	responseTime := float64(time.Duration(1*time.Second).Nanoseconds()) / 1000000
	monitorResult := datastore.AddMonitorDetail(
		monitorURL, strconv.Itoa(http.StatusOK), status, "100ms", responseTime)

	return monitorResult.ID
}
