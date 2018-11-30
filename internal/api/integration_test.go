package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
)

// Clears integraiton collection. Also clears users collection.
// We create test user to authenticate the requests. we clear them after
// every test.
func clearIntegrationCollection() {
	datastore := db.New()
	datastore.Client.Database(datastore.DatabaseName).Collection(
		db.UsersCollection).Drop(context.Background())

	datastore.Client.Database(datastore.DatabaseName).Collection(
		db.IntegrationCollection).Drop(context.Background())
}

func createTestUser() string {
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

	return jwt
}

func addTestIntegration(userID string) string {
	integrationForm := forms.IntegrationForm{
		UserID: userID,
		Type:   "email",
		Email:  "alice@sample.com",
	}

	objectID := db.GenerateObjectID()
	integrationForm.ID = objectID.Hex()

	datastore := db.New()
	integration := datastore.AddIntegration(integrationForm)

	return integration.ID.String()
}

func TestAddIntegrationHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	jwt := createTestUser()
	defer clearIntegrationCollection()

	integrationForm := forms.IntegrationForm{
		Type:  "email",
		Email: "alice@sample.com",
	}

	byte, _ := json.Marshal(integrationForm)
	req, err := http.NewRequest("POST", "localhost:8080/api/integrations", bytes.NewBuffer(byte))

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()
	AddIntegrationHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status CREATED, got %v", res.StatusCode)
	}

	response := Response{}
	json.NewDecoder(res.Body).Decode(&response)

	if response.Success == false {
		t.Errorf("response success is false")
	}
}
