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
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	decoder := json.NewDecoder(r.Body)
	var monitorURLForm forms.MonitorURLForm
	monitorURLForm.UserID = user.ID
	err := decoder.Decode(&monitorURLForm)
	if err != nil {
		writeErrorResponse(w, "Invalid input format")

		log.Info("Invalid input format for forgot password")
		return
	}

	validationMessage := monitorURLForm.Validate()
	if validationMessage != "" {
		writeErrorResponse(w, validationMessage)

		log.Info("Validation failed while adding monitoring URL.")
		return
	}

	objectID := db.GenerateObjectID()
	monitorURLForm.ID = objectID.Hex()

	datastore := db.New()
	monitoringURL := datastore.AddMonitoringURL(monitorURLForm)

	log.Info(fmt.Sprintf("Added monitoring url %s", monitorURLForm.URL))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(monitoringURL)
}

// GetMonitoringURLHandler api returns the monitoring urls configured
func GetMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	user, authErr := db.ValidateJWT(authToken)

	if authErr != nil {
		writeErrorResponse(w, "Authentication failed")

		return
	}

	datastore := db.New()
	monitoringURLS := datastore.GetMonitoringURLSByUserID(user.ID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(monitoringURLS)
}

// UpdateMonitoringURLHandler api lets the user update the details.
func UpdateMonitoringURLHandler(w http.ResponseWriter, r *http.Request) {

}
