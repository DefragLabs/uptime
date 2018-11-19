package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/defraglabs/uptime/internal/db"
	"github.com/defraglabs/uptime/internal/forms"
	log "github.com/sirupsen/logrus"
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

	error := false
	errorMsg := ""

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
		log.Info(fmt.Sprintf("Unable to add monitoring url %s", monitorURLForm.URL))
		return
	}

	objectID := db.GenerateObjectID()
	monitorURLForm.ID = objectID.Hex()

	datastore := db.New()
	monitoringURL := datastore.AddMonitoringURL(monitorURLForm)

	log.Info(fmt.Sprintf("Added monitoring url %s", monitorURLForm.URL))
	json.NewEncoder(w).Encode(monitoringURL)
}

// GetMonitoringURLHandler api returns the monitoring urls configured
func GetMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	datastore := db.New()
	monitoringURLS := datastore.GetMonitoringURLS()
	json.NewEncoder(w).Encode(monitoringURLS)
}

// UpdateMonitoringURLHandler api lets the user update the details.
func UpdateMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {

}
