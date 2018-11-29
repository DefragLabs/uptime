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

func registerTestUser(t *testing.T) string {
	userRegisterForm := forms.UserRegisterForm{
		FirstName:   "Alice",
		LastName:    "Wonderland",
		Email:       "alice@sample.com",
		Password:    "test@123",
		CompanyName: "skynet",
	}
	byte, _ := json.Marshal(userRegisterForm)
	req, err := http.NewRequest("GET", "localhost:8080/api/auth/register", bytes.NewBuffer(byte))

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	RegisterHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	response := Response{}
	json.NewDecoder(res.Body).Decode(&response)

	return response.Data["token"]
}

func TestAddIntegrationHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	jwt := registerTestUser(t)
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
