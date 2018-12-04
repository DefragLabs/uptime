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
		Unit:      "minutes",
	}

	datastore := db.New()
	monitoringURL := datastore.AddMonitoringURL(monitorURLForm)

	return monitoringURL.ID
}

func TestAddMonitoringURL(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	_, jwt := createTestUser()

	defer clearMonitorCollection()

	monitorURLForm := forms.MonitorURLForm{
		Protocol:  "http",
		URL:       "example.com",
		Frequency: 5,
		Unit:      "minute",
	}

	byte, _ := json.Marshal(monitorURLForm)
	req, err := http.NewRequest("POST", "localhost:8080/api/monitoring-urls", bytes.NewBuffer(byte))

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()
	AddMonitoringURLHandler(responseWriter, req)

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

func TestGetMonitoringURLHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	addTestMonitorURL(user.ID)
	defer clearMonitorCollection()

	req, err := http.NewRequest("GET", "localhost:8080/api/monitoring-urls", nil)

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()
	GetMonitoringURLsHandler(responseWriter, req)

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

	monitoringURLs := response.Data["monitoringURLs"].([]interface{})
	if len(monitoringURLs) != 1 {
		t.Errorf("Expected only one monitoringURL")
	}
}
