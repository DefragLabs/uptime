package api

import (
	"encoding/json"
	"net/http"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
)

// AddMonitoringURLHandler api lets an user add an healthcheck url.
func AddMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	decoder := json.NewDecoder(r.Body)
	var monitorURLForm forms.MonitorURLForm
	err := decoder.Decode(&monitorURLForm)
	if err != nil {
		panic(err)
	}

	validationMessage := monitorURLForm.Validate()
	if validationMessage != "" {
		error = true
		errorMsg = validationMessage
	}

	if error {
		errorVal := make(map[string]string)
		errorVal["message"] = errorMsg
		response := Response{
			Success: false,
			Data:    nil,
			Error:   errorVal,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	objectID := db.GenerateObjectID()
	monitorURL.ID = objectID.Hex()

	monitoringURL := db.AddMonitoringURL(monitorURLForm)
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
