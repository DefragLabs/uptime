package api

import (
	"encoding/json"
	"net/http"

	"github.com/dineshs91/uptime/internal/db"
)

// AddMonitoringURLHandler api lets an user add an healthcheck url.
func AddMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	decoder := json.NewDecoder(r.Body)
	var monitorURL db.MonitorURL
	err := decoder.Decode(&monitorURL)
	if err != nil {
		panic(err)
	}

	objectID := db.GenerateObjectID()
	monitorURL.ID = objectID.Hex()

	monitoringURL := db.AddMonitoringURL(monitorURL)
	json.NewEncoder(w).Encode(monitoringURL)
}

// GetMonitoringURLHandler api returns the monitoring urls configured
func GetMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	monitoringURLS := db.GetMonitoringURLS()
	json.NewEncoder(w).Encode(monitoringURLS)
}

// UpdateMonitoringURLHandler api lets the user update the details.
func UpdateMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {

}
