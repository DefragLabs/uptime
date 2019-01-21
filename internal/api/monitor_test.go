package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	"github.com/gorilla/mux"
)

func TestAddMonitoringURL(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	_, jwt := createTestUser()

	defer clearMonitorCollection()

	monitorURLForm := forms.MonitorURLForm{
		Protocol:  "http",
		Name:      "example",
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

func TestAddDuplicateMonitoringURL(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	addTestMonitorURL(user.ID)
	defer clearMonitorCollection()

	monitorURLForm := forms.MonitorURLForm{
		Protocol:  "http",
		Name:      "example",
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

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status BAD REQUEST, got %v", res.StatusCode)
	}

	response := Response{}
	json.NewDecoder(res.Body).Decode(&response)

	if response.Success == true {
		t.Errorf("response success should be false")
	}
}

func TestGetMonitoringURLsHandler(t *testing.T) {
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

	response := SimpleResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if response.Success == false {
		t.Errorf("response success is false")
	}

	monitoringURLs := response.Data.([]interface{})

	if len(monitoringURLs) != 1 {
		t.Errorf("Expected only one monitoringURL")
	}
}

func TestGetMonitoringURLHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	monitoringURLID := addTestMonitorURL(user.ID)
	defer clearMonitorCollection()

	req, err := http.NewRequest("GET", "localhost:8080/api/monitoring-urls", nil)
	vars := map[string]string{
		"monitoringURLID": monitoringURLID,
	}
	req = mux.SetURLVars(req, vars)

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
}

func TestUpdateMonitoringURLHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	monitoringURLID := addTestMonitorURL(user.ID)
	defer clearMonitorCollection()

	monitorURLForm := forms.MonitorURLForm{
		Protocol:  "https",
		Frequency: 30,
		Unit:      "second",
	}

	byte, _ := json.Marshal(monitorURLForm)

	req, err := http.NewRequest("PUT", "localhost:8080/api/monitoring-urls", bytes.NewBuffer(byte))
	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	// Add url path parameter
	vars := map[string]string{
		"monitoringURLID": monitoringURLID,
	}
	req = mux.SetURLVars(req, vars)

	UpdateMonitoringURLHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	response := StructResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}

	unit := response.Data["unit"]
	frequency := response.Data["frequency"]
	protocol := response.Data["protocol"]

	if unit != "second" {
		t.Error("unit should be second after update")
	} else if int(frequency.(float64)) != 30 {
		// https://stackoverflow.com/a/22344130/2134124.
		// The json package will assume float64 when you haven't defined anything else.
		t.Error("frequency should be 30 after update")
	} else if protocol != "https" {
		t.Error("protocol should be http after update")
	}
}

func TestMonitoringURLActionHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	monitoringURLID := addTestMonitorURL(user.ID)
	defer clearMonitorCollection()

	req, err := http.NewRequest("POST", "localhost:8080/api/monitoring-urls/actions?action=pause", nil)
	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()

	// Add url path parameter
	vars := map[string]string{
		"monitoringURLID": monitoringURLID,
	}
	req = mux.SetURLVars(req, vars)

	MonitoringURLActionHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	response := StructResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}

	monitoringStatus := response.Data["monitoringStatus"]
	if monitoringStatus != "paused" {
		fmt.Println("status", monitoringStatus)
		t.Error("monitoring url should be paused")
	}
}

func TestDeleteMonitoringURLHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	monitoringURLID := addTestMonitorURL(user.ID)
	defer clearMonitorCollection()

	req, err := http.NewRequest("DELETE", "localhost:8080/api/monitoring-urls", nil)

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()
	vars := map[string]string{
		"monitoringURLID": monitoringURLID,
	}
	req = mux.SetURLVars(req, vars)

	DeleteMonitoringURLHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status NoContent, got %v", res.StatusCode)
	}

	datastore := db.New()
	monitoringURL := datastore.GetMonitoringURLByUserID(user.ID, monitoringURLID)

	if monitoringURL.ID != "" {
		t.Errorf("Integration is not removed from the database.")
	}
}

func TestGetMonitoringURLStatsHandler(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	monitoringURLID := addTestMonitorURL(user.ID)
	addTestMonitorURLResult(user.ID, monitoringURLID)
	defer clearMonitorCollection()

	url := fmt.Sprintf("localhost:8080/api/monitoring-urls/%s/stats", monitoringURLID)
	req, err := http.NewRequest("DELETE", url, nil)

	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}
	responseWriter := httptest.NewRecorder()
	vars := map[string]string{
		"monitoringURLID": monitoringURLID,
	}
	req = mux.SetURLVars(req, vars)

	GetMonitoringURLStatsHandler(responseWriter, req)

	res := responseWriter.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}
}
