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

	"github.com/gorilla/mux"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
)

// Clears integration collection. Also clears users collection.
// We create test user to authenticate the requests. we clear them after
// every test.
func clearIntegrationCollection() {
	datastore := db.New()
	datastore.Client.Database(datastore.DatabaseName).Collection(
		db.UsersCollection).Drop(context.Background())

	datastore.Client.Database(datastore.DatabaseName).Collection(
		db.IntegrationCollection).Drop(context.Background())
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

	return integration.ID
}

func TestAddIntegrationHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	_, jwt := createTestUser()

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

func TestGetIntegrationsHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	addTestIntegration(user.ID)
	defer clearIntegrationCollection()

	req, err := http.NewRequest("GET", "localhost:8080/api/integrations", nil)

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	GetIntegrationsHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}

	response := StructResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if response.Success == false {
		t.Errorf("response success is false")
	}

	integrations := response.Data["integrations"].([]interface{})
	if len(integrations) != 1 {
		t.Errorf("Expected only one integration")
	}
}

func TestDeleteIntegrationHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	integrationID := addTestIntegration(user.ID)
	defer clearIntegrationCollection()

	url := fmt.Sprintf("localhost:8080/api/integrations/%s", integrationID)
	req, err := http.NewRequest("DELETE", url, nil)

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	vars := map[string]string{
		"integrationID": integrationID,
	}
	req = mux.SetURLVars(req, vars)
	DeleteIntegrationHandler(responseWriter, req)
	res := responseWriter.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status No Content, got %v", res.StatusCode)
	}

	datastore := db.New()
	integration := datastore.GetIntegrationByUserID(user.ID, integrationID)

	if integration.ID != "" {
		t.Errorf("Integration is not removed from the database.")
	}
}
