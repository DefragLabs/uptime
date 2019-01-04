package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestGetDashboardStats tests dashboard stats api.
func TestGetDashboardStats(t *testing.T) {
	os.Setenv("MONGO_DATABASE_NAME", "uptime_test")
	user, jwt := createTestUser()
	monitorURLID := addTestMonitorURL(user.ID)
	addTestMonitorURLResult(user.ID, monitorURLID)

	defer clearMonitorCollection()

	req, err := http.NewRequest("GET", "localhost:8080/api/dashboard/stats", nil)
	token := fmt.Sprintf("JWT %s", jwt)
	req.Header.Add("Authorization", token)

	if err != nil {
		t.Errorf("Unable to create a new request")
	}

	responseWriter := httptest.NewRecorder()
	DashboardStatsHandler(responseWriter, req)

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

	monitoringURLCount := response.Data["monitoring_urls_count"].(float64)
	upMonitoringURLCount := response.Data["up_monitoring_urls_count"].(float64)
	downMonitoringURLCount := response.Data["down_monitoring_urls_count"].(float64)

	if monitoringURLCount != 1 {
		t.Errorf("monitoring url count should be 1")
	} else if upMonitoringURLCount != 1 {
		t.Errorf("up monitoring url count should be 1")
	} else if downMonitoringURLCount != 0 {
		t.Errorf("down monitoring url count should be 0")
	}
}
