package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dineshs91/uptime/internal/db"
)

// AddMonitoringDetailHandler api lets an user add an healthcheck url.
func AddMonitoringDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	decoder := json.NewDecoder(r.Body)
	var monitorURL db.MonitorURL
	err := decoder.Decode(&monitorURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(monitorURL)
	monitoringDetail := db.AddMonitoringDetail(monitorURL)
	fmt.Println(monitoringDetail)
	json.NewEncoder(w).Encode(monitoringDetail)
}

// GetMonitoringDetailHandler api returns the monitoring urls configured
func GetMonitoringDetailHandler(w http.ResponseWriter, r *http.Request) {

}

// UpdateMonitoringDetailHandler api lets the user update the details.
func UpdateMonitoringDetailHandler(w http.ResponseWriter, r *http.Request) {

}
